package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeGroupingKey(t *testing.T, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "grouping.key")
	require.NoError(t, os.WriteFile(path, content, 0o600))
	return path
}

func TestCLIKeyedGroupedReportIsReportOnly(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "b.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	keyPath := writeGroupingKey(t, []byte("0123456789abcdef0123456789abcdef"))
	primary := filepath.Join(t.TempDir(), "primary.json")
	grouped := filepath.Join(t.TempDir(), "grouped.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--grouped-report", grouped,
		"--grouping-key-file", keyPath,
		"--report-path", primary,
		"--report-format", "json",
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)

	var primaryFindings []map[string]interface{}
	require.NoError(t, json.Unmarshal(mustReadFile(t, primary), &primaryFindings))
	assert.Len(t, primaryFindings, 2)

	var document struct {
		SchemaVersion int `json:"schema_version"`
		Groups        []struct {
			ID          string `json:"id"`
			Count       int    `json:"count"`
			Occurrences []struct {
				File string `json:"file"`
			} `json:"occurrences"`
		} `json:"groups"`
	}
	raw := mustReadFile(t, grouped)
	require.NoError(t, json.Unmarshal(raw, &document))
	require.Equal(t, 1, document.SchemaVersion)
	require.Len(t, document.Groups, 1)
	assert.Equal(t, 2, document.Groups[0].Count)
	assert.Contains(t, document.Groups[0].ID, "hmac-sha256:")
	assert.Len(t, document.Groups[0].Occurrences, 2)
	assert.NotContains(t, string(raw), "ALPHA_12345678")
	assert.NotContains(t, string(raw), "REDACTED")
}

func TestCLIKeyedGroupedReportSurvivesRedaction(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "b.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	keyPath := writeGroupingKey(t, []byte("0123456789abcdef0123456789abcdef"))
	grouped := filepath.Join(t.TempDir(), "grouped.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--grouped-report", grouped,
		"--grouping-key-file", keyPath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	var document struct {
		Groups []struct {
			Count int `json:"count"`
		} `json:"groups"`
	}
	raw := mustReadFile(t, grouped)
	require.NoError(t, json.Unmarshal(raw, &document))
	require.Len(t, document.Groups, 1)
	assert.Equal(t, 2, document.Groups[0].Count)
	assert.NotContains(t, string(raw), "ALPHA_12345678")
}

func TestCLIKeyedGroupedReportRejectsPrimaryCollision(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	keyPath := writeGroupingKey(t, []byte("0123456789abcdef0123456789abcdef"))
	dest := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--grouped-report", dest,
		"--grouping-key-file", keyPath,
		"--report-path", dest,
		"--report-format", "json",
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "must not use the primary report path")
	assert.NoFileExists(t, dest)
}

func TestCLIKeyedGroupedReportRequiresStrongKey(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	grouped := filepath.Join(t.TempDir(), "grouped.json")
	keyPath := writeGroupingKey(t, []byte("short"))

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--grouped-report", grouped,
		"--grouping-key-file", keyPath,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "at least 32 bytes")
	assert.NoFileExists(t, grouped)
}

func TestCLIKeyedGroupedReportRequiresBothFlags(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")
	grouped := filepath.Join(t.TempDir(), "grouped.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--grouped-report", grouped,
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "must be used together")
	assert.NoFileExists(t, grouped)
}
