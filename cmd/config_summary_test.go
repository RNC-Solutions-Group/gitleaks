package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

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

func TestSanitizeConfigSummary(t *testing.T) {
	input := "ok\x1b[31mred\nline\tTab\u2028sep\u202ebidi"
	got := sanitizeConfigSummary(input, 80)
	assert.NotContains(t, got, "\x1b")
	assert.NotContains(t, got, "\n")
	assert.NotContains(t, got, "\t")
	assert.NotContains(t, got, "\u2028")
	assert.NotContains(t, got, "\u202e")
	assert.Contains(t, got, `\x1b`)
	assert.Contains(t, got, `\n`)
	assert.Contains(t, got, `\u2028`)
	assert.Contains(t, got, `\u202e`)

	long := strings.Repeat("\u00e9", 300)
	truncated := sanitizeConfigSummary(long, 10)
	assert.Equal(t, 10, utf8.RuneCountInString(strings.TrimSuffix(truncated, "...")))
	assert.True(t, strings.HasSuffix(truncated, "..."))
}

func TestConfigSummaryEscapesHostileTitle(t *testing.T) {
	source := t.TempDir()
	title := "evil\\u001b[31mred\\ninjected\\u2028separator"
	config := "title = \"" + title + "\"\n[[rules]]\nid = \"alpha\"\ndescription = \"alpha\"\nregex = '''ALPHA_[A-Z0-9]{8}'''\n"
	configPath := filepath.Join(t.TempDir(), "hostile.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(config), 0o600))

	result := runReportCLI(t, "config", source, "--config", configPath, "--no-banner", "--log-level", "error")
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.NotContains(t, result.stdout, "\x1b")
	assert.NotContains(t, result.stdout, "red\ninjected")
	assert.Contains(t, result.stdout, `\x1b`)
	assert.Contains(t, result.stdout, `\n`)
	assert.NotContains(t, result.stdout, "ALPHA_[A-Z0-9]{8}")
}

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
