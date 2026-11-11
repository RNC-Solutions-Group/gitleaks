//go:build windows

package cmd

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func lockScanFile(t *testing.T, path string) {
	t.Helper()
	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(path),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = syscall.CloseHandle(handle) })
}
