//go:build windows

package sources

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
)

func TestFilesReportsWindowsShareLockedFile(t *testing.T) {
	dir := t.TempDir()
	locked := filepath.Join(dir, "locked.txt")
	require.NoError(t, os.WriteFile(locked, []byte("synthetic secret\n"), 0o600))

	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(locked),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = syscall.CloseHandle(handle) })

	source := &Files{
		Config: &config.Config{},
		Path:   dir,
		Sema:   semgroup.NewGroup(context.Background(), 4),
	}
	err = source.Fragments(context.Background(), func(Fragment, error) error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "locked.txt")
	assert.Contains(t, err.Error(), "could not open file")
}
