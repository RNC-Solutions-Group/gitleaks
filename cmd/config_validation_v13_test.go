package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigValidationEmptyRuleBlockFailsClosed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.toml")
	require.NoError(t, os.WriteFile(path, []byte("title = \"empty\"\n[[rules]]\n"), 0o600))

	result := runReportCLI(t, "config", t.TempDir(), "--config", path, "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "rule |id| is missing or empty")
}

func TestConfigValidationMalformedTOMLFailsClosed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.toml")
	require.NoError(t, os.WriteFile(path, []byte("title = \"bad\"\n[[rules]\n"), 0o600))

	result := runReportCLI(t, "config", t.TempDir(), "--config", path, "--no-banner", "--log-level", "error")
	require.NotEqual(t, 0, result.exitCode)
	assert.NotEmpty(t, result.stderr)
}
