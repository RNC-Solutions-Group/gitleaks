package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func translateTOML(t *testing.T, content string) Config {
	t.Helper()
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.SetConfigType("toml")
	require.NoError(t, viper.ReadConfig(strings.NewReader(content)))

	var vc ViperConfig
	require.NoError(t, viper.Unmarshal(&vc))
	cfg, err := vc.Translate()
	require.NoError(t, err)
	return cfg
}

func TestInformativeSettingsDefaultFalse(t *testing.T) {
	cfg := translateTOML(t, `title = "default informative settings"`)
	assert.False(t, cfg.InformativeReport)
	assert.False(t, cfg.ConsoleOnly)
}

func TestInformativeSettingsExplicitValues(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		informative bool
		consoleOnly bool
	}{
		{
			name:        "explicit false",
			content:     "informative-report = false\nconsole-only = false\n",
			informative: false,
			consoleOnly: false,
		},
		{
			name:        "explicit true",
			content:     "informative-report = true\nconsole-only = true\n",
			informative: true,
			consoleOnly: true,
		},
		{
			name:        "console only without informative report",
			content:     "console-only = true\n",
			informative: false,
			consoleOnly: true,
		},
		{
			name:        "informative report only",
			content:     "informative-report = true\n",
			informative: true,
			consoleOnly: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := translateTOML(t, test.content)
			assert.Equal(t, test.informative, cfg.InformativeReport)
			assert.Equal(t, test.consoleOnly, cfg.ConsoleOnly)
		})
	}
}

func TestInformativeSettingsSurviveDefaultExtension(t *testing.T) {
	cfg := translateTOML(t, `
title = "custom rules with informative settings"
informative-report = true
console-only = true

[extend]
useDefault = true
`)
	assert.True(t, cfg.InformativeReport)
	assert.True(t, cfg.ConsoleOnly)
	assert.NotEmpty(t, cfg.Rules)
	assert.Contains(t, cfg.Rules, "generic-api-key")
}

func TestInformativeSettingsDoNotComeFromExtendedConfig(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "base.toml")
	require.NoError(t, os.WriteFile(basePath, []byte(
		"title = \"base\"\ninformative-report = true\nconsole-only = true\n",
	), 0o600))
	cfg := translateTOML(t, fmt.Sprintf(`
title = "top level"
informative-report = false
console-only = false

[extend]
path = %q
`, filepath.ToSlash(basePath)))
	assert.False(t, cfg.InformativeReport)
	assert.False(t, cfg.ConsoleOnly)
	assert.Contains(t, cfg.ConfigFiles, filepath.ToSlash(basePath))
}
