package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIStrictConfigRejectsUnknownFields(t *testing.T) {
	source := t.TempDir()
	valid := filepath.Join(t.TempDir(), "valid.toml")
	require.NoError(t, os.WriteFile(valid, []byte(configSummaryConfig), 0o600))

	result := runReportCLI(t, "config", source, "--config", valid, "--strict-config", "--no-banner", "--log-level", "error")
	require.Equal(t, 0, result.exitCode, result.stderr)

	unknown := filepath.Join(t.TempDir(), "unknown.toml")
	require.NoError(t, os.WriteFile(unknown, []byte("title = \"bad\"\nunknown-field = \"sensitive-value\"\n"), 0o600))
	result = runReportCLI(t, "config", source, "--config", unknown, "--strict-config", "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "unsupported field")
	assert.NotContains(t, result.stderr, "sensitive-value")

	// Default non-strict behavior remains compatible.
	result = runReportCLI(t, "config", source, "--config", unknown, "--no-banner", "--log-level", "error")
	require.Equal(t, 0, result.exitCode, result.stderr)
}

func TestCLIStrictConfigRejectsUnknownRuleField(t *testing.T) {
	source := t.TempDir()
	path := filepath.Join(t.TempDir(), "unknown-rule.toml")
	content := "title = \"bad\"\n[[rules]]\nid = \"alpha\"\ndescription = \"alpha\"\nregex = '''ALPHA_[A-Z0-9]{8}'''\nunknownRuleField = \"sensitive-value\"\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	result := runReportCLI(t, "config", source, "--config", path, "--strict-config", "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "rules[0]")
	assert.Contains(t, result.stderr, "unknownRuleField")
	assert.NotContains(t, result.stderr, "sensitive-value")
}

func TestCLIStrictConfigValidatesExtensionChain(t *testing.T) {
	source := t.TempDir()
	child := filepath.Join(t.TempDir(), "child.toml")
	require.NoError(t, os.WriteFile(child, []byte("title = \"child\"\nunknownChildField = true\n"), 0o600))
	base := filepath.Join(t.TempDir(), "base.toml")
	content := fmt.Sprintf("title = \"base\"\n[extend]\npath = %q\nuseDefault = false\n", filepath.ToSlash(child))
	require.NoError(t, os.WriteFile(base, []byte(content), 0o600))

	result := runReportCLI(t, "config", source, "--config", base, "--strict-config", "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "strict configuration")
	assert.Contains(t, result.stderr, "unsupported field")
}
