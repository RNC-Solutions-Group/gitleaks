package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configSummaryConfig = `title = "summary-test"
[[rules]]
id = "alpha"
description = "alpha"
regex = '''ALPHA_[A-Z0-9]{8}'''
[[rules]]
id = "beta"
description = "beta"
regex = '''BETA_[A-Z0-9]{8}'''
[[allowlists]]
regexes = ['''BETA_[A-Z0-9]{8}''']
[[rules.allowlists]]
regexes = ['''ALPHA_[A-Z0-9]{8}''']
`

func TestConfigSummaryIsSafeAndEffective(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "summary.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(configSummaryConfig), 0o600))
	source := t.TempDir()

	result := runReportCLI(t, "config", source, "--config", configPath, "--no-banner", "--log-level", "error")
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, result.stdout, "Rules active: 2")
	assert.Contains(t, result.stdout, "Global allowlists: 1")
	assert.Contains(t, result.stdout, "Rule allowlists: 1")
	assert.Contains(t, result.stdout, "Report policy: auto-overwrite=true aggregation=false")
	assert.Contains(t, result.stdout, "Sensitive rule/allowlist patterns are intentionally not printed.")
	assert.NotContains(t, result.stdout, "ALPHA_[A-Z0-9]{8}")
	assert.NotContains(t, result.stdout, "BETA_[A-Z0-9]{8}")
	assert.NotContains(t, result.stdout, "ALPHA_12345678")
}

func TestConfigSummaryUsesTargetLocalConfig(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaks.toml"), []byte(configSummaryConfig), 0o600))

	result := runReportCLI(t, "config", source, "--no-banner", "--log-level", "error")
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, result.stdout, "Configuration source: target-local .gitleaks.toml")
	assert.Contains(t, result.stdout, "Rules active: 2")
}

func TestConfigSummaryInvalidConfigFailsClosed(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "invalid.toml")
	require.NoError(t, os.WriteFile(configPath, []byte("[[rules]]\nid=\"bad\"\nregex='''[unterminated'''\n"), 0o600))

	result := runReportCLI(t, "config", t.TempDir(), "--config", configPath, "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.NotEmpty(t, result.stderr)
}
