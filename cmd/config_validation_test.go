package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIInvalidConfigFailsBeforeScanning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	configPath := filepath.Join(t.TempDir(), "invalid.toml")
	require.NoError(t, os.WriteFile(configPath, []byte("[[rules]]\nid = \"bad\"\nregex = \"[a-z\"\n"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "invalid rule regex")
}

func TestCLIInvalidAllowlistConfigFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	configPath := filepath.Join(t.TempDir(), "invalid.toml")
	require.NoError(t, os.WriteFile(configPath, []byte("[[allowlists]]\npaths = [\"*.test.js\"]\n"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "invalid allowlist path regex")
}

func TestCLIInvalidExtendedConfigFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	basePath := filepath.Join(t.TempDir(), "base.toml")
	require.NoError(t, os.WriteFile(basePath, []byte("[[rules]]\nid = \"bad\"\nregex = \"[a-z\"\n"), 0o600))
	configPath := filepath.Join(t.TempDir(), "main.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(fmt.Sprintf("[extend]\npath = %q\n", filepath.ToSlash(basePath))), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "failed to load extended config")
}
