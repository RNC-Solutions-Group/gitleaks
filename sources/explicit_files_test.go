package sources

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collectExplicitFragments(t *testing.T, source *ExplicitFiles) []Fragment {
	t.Helper()
	var fragments []Fragment
	err := source.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	return fragments
}

func TestExplicitFilesOnlyScansNamedFiles(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.txt")
	second := filepath.Join(dir, "second.txt")
	require.NoError(t, os.WriteFile(first, []byte("first\n"), 0o600))
	require.NoError(t, os.WriteFile(second, []byte("second\n"), 0o600))

	fragments := collectExplicitFragments(t, &ExplicitFiles{Paths: []string{first}})
	require.Len(t, fragments, 1)
	assert.Equal(t, filepath.ToSlash(first), fragments[0].FilePath)
	assert.Equal(t, "first\n", fragments[0].Raw)
}

func TestExplicitFilesRejectsDirectories(t *testing.T) {
	dir := t.TempDir()
	err := (&ExplicitFiles{Paths: []string{dir}}).Fragments(context.Background(), func(Fragment, error) error {
		return nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is a directory")
}

func TestExplicitFilesRejectsMissingAndSkipsEmpty(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty.txt")
	require.NoError(t, os.WriteFile(empty, nil, 0o600))
	var fragments []Fragment
	err := (&ExplicitFiles{
		Paths: []string{filepath.Join(dir, "missing.txt"), empty},
	}).Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	assert.Empty(t, fragments)
}

func TestExplicitFilesDeduplicatesPaths(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "one.txt")
	require.NoError(t, os.WriteFile(file, []byte("one\n"), 0o600))
	fragments := collectExplicitFragments(t, &ExplicitFiles{Paths: []string{file, file}})
	require.Len(t, fragments, 1)
}
