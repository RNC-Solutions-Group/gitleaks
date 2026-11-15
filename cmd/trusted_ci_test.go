package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrustedCIConfigurationContract(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("CPV14_TRUSTED_12345678\n"), 0o600))
	// A hostile repository-local config that suppresses the secret.
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaks.toml"), []byte("title = \"hostile local policy\"\n"), 0o600))
	trusted := filepath.Join(t.TempDir(), "trusted.toml")
	require.NoError(t, os.WriteFile(trusted, []byte(`title = "trusted"
[[rules]]
id = "trusted"
description = "trusted"
regex = '''CPV14_TRUSTED_[0-9]+'''
`), 0o600))

	reportPath := filepath.Join(t.TempDir(), "report.json")
	run := func(args ...string) reportCLIResult {
		base := []string{"dir", source, "--no-banner", "--log-level", "error", "--exit-code", "0", "--report-path", reportPath, "--report-format", "json"}
		return runReportCLI(t, append(base, args...)...)
	}

	// Compatible default: a source-local config is still honored.
	result := run()
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))

	// Trusted CI contract: explicit trusted config + no local discovery.
	result = run("--no-local-config", "--config", trusted)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, string(mustReadFile(t, reportPath)), "CPV14_TRUSTED_12345678")

	// Environment-selected trusted config retains precedence over local.
	result = runReportCLIWithEnv(t, map[string]string{"GITLEAKS_CONFIG": trusted},
		"dir", source, "--no-local-config", "--no-banner", "--log-level", "error", "--exit-code", "0", "--report-path", reportPath, "--report-format", "json")
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, string(mustReadFile(t, reportPath)), "CPV14_TRUSTED_12345678")
}
