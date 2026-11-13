package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test documents the established upstream configuration contract: an
// automatically discovered source/.gitleaks.toml is a deliberate custom
// configuration and can replace the default rule set. CP-V6 does not change
// that policy; untrusted callers must supply an explicit or environment
// configuration. The explicit and environment sources still take precedence.
func TestCLINoLocalConfigUsesDefaultRules(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("AKIAQ7Z3R5T7Y2U2I4O6\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaks.toml"), []byte("title = \"untrusted local policy without rules\"\n"), 0o600))

	run := func(extra ...string) reportCLIResult {
		reportPath := filepath.Join(t.TempDir(), "report.json")
		args := []string{"dir", source, "--no-banner", "--log-level", "error", "--report-format", "json", "--report-path", reportPath}
		args = append(args, extra...)
		result := runReportCLI(t, args...)
		if _, err := os.Stat(reportPath); err == nil {
			result.stdout = string(mustReadFile(t, reportPath))
		}
		return result
	}

	without := run()
	require.Equal(t, 0, without.exitCode, without.stderr)
	assert.JSONEq(t, "[]", without.stdout)

	with := run("--no-local-config")
	require.Equal(t, 1, with.exitCode, with.stderr)
	assert.Contains(t, with.stdout, "AKIAQ7Z3R5T7Y2U2I4O6")
}

func TestCLIConfigurationPrecedenceIsPreserved(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("CPV6_12345678\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaks.toml"), []byte("title = \"local policy without rules\"\n"), 0o600))

	customConfig := filepath.Join(t.TempDir(), "custom.toml")
	require.NoError(t, os.WriteFile(customConfig, []byte(cpv6IgnoreConfig), 0o600))
	content := cpv6IgnoreConfig

	tests := []struct {
		name  string
		env   map[string]string
		args  []string
		found bool
	}{
		{
			name:  "automatic local config",
			found: false,
		},
		{
			name:  "explicit config",
			args:  []string{"--config", customConfig},
			found: true,
		},
		{
			name:  "environment config path",
			env:   map[string]string{"GITLEAKS_CONFIG": customConfig},
			found: true,
		},
		{
			name:  "environment config content",
			env:   map[string]string{"GITLEAKS_CONFIG_TOML": content},
			found: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportPath := filepath.Join(t.TempDir(), "report.json")
			args := []string{"dir", source, "--no-banner", "--log-level", "error", "--report-path", reportPath, "--report-format", "json"}
			args = append(args, tt.args...)
			result := runReportCLIWithEnv(t, tt.env, args...)
			if tt.found {
				require.Equal(t, 1, result.exitCode, result.stderr)
				assert.Contains(t, string(mustReadFile(t, reportPath)), "CPV6_12345678")
			} else {
				require.Equal(t, 0, result.exitCode, result.stderr)
				assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
			}
		})
	}
}
