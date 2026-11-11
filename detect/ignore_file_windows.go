//go:build windows

package detect

import (
	"fmt"
	"os"
)

// openIgnoreFile refuses non-regular files before reading. Windows does not
// expose POSIX FIFOs, but reparse points and other special paths still need the
// same regular-file boundary.
func openIgnoreFile(path string) (*os.File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not inspect .gitleaksignore %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf(".gitleaksignore %q is not a regular file", path)
	}

	file, err := os.Open(path)
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
