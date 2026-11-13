//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
)

func TestCLIIgnoreDiscoveryRejectsFIFOWithoutBlocking(t *testing.T) {
	fifo := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))
	d := detect.NewDetector(config.Config{})

	start := time.Now()
	require.Error(t, loadGitleaksIgnoreFile(d, fifo))
	assert.Less(t, time.Since(start), 2*time.Second)

	start = time.Now()
	require.Error(t, loadGitleaksIgnorePath(d, fifo, true))
	assert.Less(t, time.Since(start), 2*time.Second)
}

func TestCLIIgnoreDiscoveryFollowsRegularSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	require.NoError(t, os.WriteFile(target, []byte("file:rule:1\n"), 0o600))
	link := filepath.Join(dir, ".gitleaksignore")
	require.NoError(t, os.Symlink(target, link))
	d := detect.NewDetector(config.Config{})

	require.NoError(t, loadGitleaksIgnoreFile(d, link))
}

func TestCLIIgnoreDiscoveryRejectsBrokenSymlink(t *testing.T) {
	link := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.Symlink(filepath.Join(t.TempDir(), "missing"), link))
	d := detect.NewDetector(config.Config{})
	require.Error(t, loadGitleaksIgnoreFile(d, link))
	require.Error(t, loadGitleaksIgnorePath(d, link, true))
}
