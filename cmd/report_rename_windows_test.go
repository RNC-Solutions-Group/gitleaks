//go:build windows

package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
)

func TestIsTransientWindowsRenameError(t *testing.T) {
	assert.True(t, isTransientWindowsRenameError(&os.LinkError{Op: "rename", Err: windows.ERROR_ACCESS_DENIED}))
	assert.True(t, isTransientWindowsRenameError(&os.LinkError{Op: "rename", Err: windows.ERROR_SHARING_VIOLATION}))
	assert.False(t, isTransientWindowsRenameError(errors.New("not a sharing error")))
	assert.False(t, isTransientWindowsRenameError(os.ErrNotExist))
}
