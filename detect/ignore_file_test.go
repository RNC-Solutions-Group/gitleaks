package detect

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
)

type ignoreErrorReader struct {
	data []byte
	err  error
}

func (r *ignoreErrorReader) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		return n, nil
	}
	return 0, r.err
}

func newIgnoreDetector() *Detector {
	return NewDetector(config.Config{})
}

func TestAddGitleaksIgnoreNormalizesAndDeduplicates(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.WriteFile(path, []byte("# comment\n\nfoo\\bar:rule:1\nfoo/bar:rule:1\ncommit:foo\\bar:rule:2\ninvalid\n"), 0o600))

	d := newIgnoreDetector()
	require.NoError(t, d.AddGitleaksIgnore(path))
	assert.Contains(t, d.gitleaksIgnore, "foo/bar:rule:1")
	assert.Contains(t, d.gitleaksIgnore, "commit:foo/bar:rule:2")
	assert.Contains(t, d.gitleaksIgnore, "invalid")
	assert.Len(t, d.gitleaksIgnore, 3)
}

func TestAddGitleaksIgnoreEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.WriteFile(path, nil, 0o600))

	d := newIgnoreDetector()
	require.NoError(t, d.AddGitleaksIgnore(path))
	assert.Empty(t, d.gitleaksIgnore)
}

func TestAddGitleaksIgnoreRejectsDirectory(t *testing.T) {
	d := newIgnoreDetector()
	err := d.AddGitleaksIgnore(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a regular file")
}

func TestAddGitleaksIgnoreRejectsOversizedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.WriteFile(path, nil, 0o600))
	require.NoError(t, os.Truncate(path, maxGitleaksIgnoreBytes+1))

	d := newIgnoreDetector()
	err := d.AddGitleaksIgnore(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds limit")
}

func TestAddGitleaksIgnoreRejectsOversizedLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.WriteFile(path, []byte(strings.Repeat("a", maxGitleaksIgnoreLine+1)), 0o600))

	d := newIgnoreDetector()
	err := d.AddGitleaksIgnore(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token too long")
}

func TestAddGitleaksIgnoreRejectsTooManyEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitleaksignore")
	var content strings.Builder
	for i := 0; i <= maxGitleaksIgnoreEntries; i++ {
		content.WriteString("file:rule:")
		content.WriteString("1\n")
	}
	require.NoError(t, os.WriteFile(path, []byte(content.String()), 0o600))

	d := newIgnoreDetector()
	err := d.AddGitleaksIgnore(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "more than")
}

func TestAddGitleaksIgnoreReadErrorDoesNotMergePartialEntries(t *testing.T) {
	d := newIgnoreDetector()
	d.gitleaksIgnore["existing"] = struct{}{}

	entries, err := parseGitleaksIgnore(&ignoreErrorReader{
		data: []byte("file:rule:1\n"),
		err:  errors.New("synthetic read failure"),
	})
	require.Error(t, err)
	assert.Nil(t, entries)
	assert.Contains(t, err.Error(), "synthetic read failure")
	assert.NotContains(t, d.gitleaksIgnore, "file:rule:1")
	assert.Contains(t, d.gitleaksIgnore, "existing")
}

func TestAddGitleaksIgnoreMissingFileIsError(t *testing.T) {
	d := newIgnoreDetector()
	require.Error(t, d.AddGitleaksIgnore(filepath.Join(t.TempDir(), "missing")))
}
