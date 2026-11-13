package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/report"
)

func TestFilesCommandScansOnlyExplicitFiles(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(t.TempDir(), "gitleaks.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`title = "explicit files"
[[rules]]
id = "explicit-test"
regex = '''SECRET_[A-Z0-9]{8}'''
`), 0o600))

	first := filepath.Join(dir, "first.txt")
	second := filepath.Join(dir, "second.txt")
	unlisted := filepath.Join(dir, "unlisted.txt")
	require.NoError(t, os.WriteFile(first, []byte("value=SECRET_AAAA1111\n"), 0o600))
	require.NoError(t, os.WriteFile(second, []byte("value=SECRET_BBBB2222\n"), 0o600))
	require.NoError(t, os.WriteFile(unlisted, []byte("value=SECRET_CCCC3333\n"), 0o600))

	readFindings := func(path string) []report.Finding {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)
		var findings []report.Finding
		require.NoError(t, json.Unmarshal(raw, &findings))
		return findings
	}

	firstReport := filepath.Join(t.TempDir(), "first.json")
	result := runReportCLI(t, "files", first, "--config", configPath, "--exit-code", "0",
		"--no-banner", "--log-level", "error", "--report-format", "json", "--report-path", firstReport)
	require.Equal(t, 0, result.exitCode, result.stderr)
	findings := readFindings(firstReport)
	require.Len(t, findings, 1)
	assert.Equal(t, filepath.ToSlash(first), findings[0].File)

	bothReport := filepath.Join(t.TempDir(), "both.json")
	result = runReportCLI(t, "files", first, second, "--config", configPath, "--exit-code", "0",
		"--no-banner", "--log-level", "error", "--report-format", "json", "--report-path", bothReport)
	require.Equal(t, 0, result.exitCode, result.stderr)
	findings = readFindings(bothReport)
	require.Len(t, findings, 2)
	for _, finding := range findings {
		assert.NotEqual(t, filepath.ToSlash(unlisted), finding.File)
	}

	dirResult := runReportCLI(t, "files", dir, "--config", configPath, "--no-banner", "--log-level", "error")
	assert.NotEqual(t, 0, dirResult.exitCode, "directory arguments must not expand into a repository scan")
}
