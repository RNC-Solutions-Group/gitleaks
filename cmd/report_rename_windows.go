//go:build windows

package cmd

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// renameReportFile handles the short Windows sharing/access-denied window that
// can occur while replacing a newly reserved report file. The allocation itself
// remains exclusive; only the final rename is retried.
func renameReportFile(oldPath, newPath string) error {
	return renameReportFileWith(os.Rename, time.Sleep, isTransientWindowsRenameError, oldPath, newPath)
}

func isTransientWindowsRenameError(err error) bool {
	return errors.Is(err, windows.ERROR_ACCESS_DENIED) ||
		errors.Is(err, windows.ERROR_SHARING_VIOLATION)
}
