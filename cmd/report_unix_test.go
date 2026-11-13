//go:build unix

package cmd

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportOutputRejectsSpecialFile(t *testing.T) {
	dir := t.TempDir()
	fifo := filepath.Join(dir, "report.fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Skipf("could not create a FIFO: %v", err)
	}

	_, err := newReportOutput(fifo, reportOverwrite, "", nil, "", ".", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a regular file")
}
