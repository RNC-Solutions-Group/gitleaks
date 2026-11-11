//go:build !windows

package detect

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGitleaksIgnoreFollowsRegularSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	require.NoError(t, os.WriteFile(target, []byte("file:rule:1\n"), 0o600))
	link := filepath.Join(dir, ".gitleaksignore")
	require.NoError(t, os.Symlink(target, link))

	d := newIgnoreDetector()
	require.NoError(t, d.AddGitleaksIgnore(link))
	assert.Contains(t, d.gitleaksIgnore, "file:rule:1")
}

func TestAddGitleaksIgnoreRejectsBrokenSymlink(t *testing.T) {
	link := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.Symlink(filepath.Join(t.TempDir(), "missing"), link))

	d := newIgnoreDetector()
	require.Error(t, d.AddGitleaksIgnore(link))
}

func TestAddGitleaksIgnoreRejectsFIFOWithoutBlocking(t *testing.T) {
	fifo := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))

	d := newIgnoreDetector()
	done := make(chan error, 1)
	go func() {
		done <- d.AddGitleaksIgnore(fifo)
	}()

	select {
	case err := <-done:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a regular file")
	case <-time.After(2 * time.Second):
		t.Fatal(".gitleaksignore FIFO blocked open")
	}
}
