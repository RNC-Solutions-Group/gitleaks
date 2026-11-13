//go:build !windows

package cmd

import "os"

func renameReportFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
