package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIFileLinksConsoleOnly(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--file-links",
		"--no-color",
		"--verbose",
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, result.stdout, "File-Link:")
	assert.Contains(t, result.stdout, "file:///")
	assert.Contains(t, result.stdout, "#L1")
	assert.NotContains(t, result.stdout, "\x1b")
	assert.NotContains(t, string(mustReadFile(t, reportPath)), "File-Link")
}
