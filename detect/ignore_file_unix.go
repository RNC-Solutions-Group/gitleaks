//go:build !windows

package detect

import (
	"fmt"
	"os"
	"syscall"
)

// openIgnoreFile refuses non-regular files before reading. O_NONBLOCK makes a
// FIFO replacement between the metadata check and open non-blocking; the
// post-open fstat check then rejects it before any content is consumed.
func openIgnoreFile(path string) (*os.File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not inspect .gitleaksignore %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf(".gitleaksignore %q is not a regular file", path)
	}

	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("could not open .gitleaksignore %q: %w", path, err)
	}
	openedInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("could not stat .gitleaksignore %q: %w", path, err)
	}
	if !openedInfo.Mode().IsRegular() {
		_ = file.Close()
		return nil, fmt.Errorf(".gitleaksignore %q is not a regular file", path)
	}
	return file, nil
}
