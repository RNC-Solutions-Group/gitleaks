package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cpv6IgnoreConfig = `title = "cpv6 ignore"
[[rules]]
id = "cpv6-secret"
description = "synthetic"
regex = '''CPV6_[A-Z0-9]{8}'''
`

func cpv6IgnoreSource(t *testing.T) (string, string) {
	t.Helper()
	source := t.TempDir()
	configPath := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(cpv6IgnoreConfig), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("CPV6_12345678\n"), 0o600))
	return source, configPath
}

func TestCLIIgnoreFileSuppressesFinding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaksignore"), []byte(filepath.ToSlash(filepath.Join(source, "secret.txt"))+":cpv6-secret:1\n"), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLIExplicitIgnoreFileSuppressesFinding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	ignorePath := filepath.Join(t.TempDir(), "ignore")
	require.NoError(t, os.WriteFile(ignorePath, []byte(filepath.ToSlash(filepath.Join(source, "secret.txt"))+":cpv6-secret:1\n"), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--gitleaks-ignore-path", ignorePath,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLIIgnoreDirectorySuppressesFinding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	ignoreDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(ignoreDir, ".gitleaksignore"), []byte(filepath.ToSlash(filepath.Join(source, "secret.txt"))+":cpv6-secret:1\n"), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--gitleaks-ignore-path", ignoreDir,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLIEmptyIgnoreFileIsAccepted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaksignore"), nil, 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode, result.stderr)
	assert.Contains(t, string(mustReadFile(t, reportPath)), "CPV6_12345678")
}

func TestCLIExplicitMissingIgnorePathFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--gitleaks-ignore-path", filepath.Join(t.TempDir(), "missing"),
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "could not inspect .gitleaksignore path")
}

func TestCLIOversizedIgnoreFileFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source, configPath := cpv6IgnoreSource(t)
	ignorePath := filepath.Join(t.TempDir(), ".gitleaksignore")
	require.NoError(t, os.WriteFile(ignorePath, nil, 0o600))
	require.NoError(t, os.Truncate(ignorePath, (16<<20)+1))
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--config", configPath,
		"--gitleaks-ignore-path", ignorePath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "exceeds limit")
}
