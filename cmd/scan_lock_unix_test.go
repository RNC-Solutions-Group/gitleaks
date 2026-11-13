//go:build !windows

package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func lockScanFile(t *testing.T, path string) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}
	require.NoError(t, os.Chmod(path, 0o000))
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
}
