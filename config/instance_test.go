package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadWithLocalViper(t *testing.T, path string) Config {
	t.Helper()
	v := viper.New()
	v.SetConfigFile(path)
	require.NoError(t, v.ReadInConfig())
	var vc ViperConfig
	require.NoError(t, v.Unmarshal(&vc))
	vc.ConfigPath = path
	cfg, err := vc.Translate()
	require.NoError(t, err)
	return cfg
}

func TestTranslateUsesExplicitConfigPath(t *testing.T) {
	dir := t.TempDir()
	wanted := filepath.Join(dir, "wanted.toml")
	decoy := filepath.Join(dir, "decoy.toml")
	require.NoError(t, os.WriteFile(wanted, []byte("title = \"wanted\"\n"), 0o600))
	require.NoError(t, os.WriteFile(decoy, []byte("title = \"decoy\"\n"), 0o600))

	// A package-global Viper value must not influence translation. This is the
	// regression that made two configurations unsafe to load in one process.
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.SetConfigFile(decoy)
	require.NoError(t, viper.ReadInConfig())

	cfg := loadWithLocalViper(t, wanted)
	assert.Equal(t, wanted, cfg.Path)
	assert.Equal(t, "wanted", cfg.Title)
	require.Len(t, cfg.ConfigFiles, 1)
	assert.Equal(t, wanted, cfg.ConfigFiles[0])
}

func TestTranslateConcurrentConfigsStayIndependent(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
		return path
	}

	baseA := write("base-a.toml", `
title = "base-a"
[[rules]]
id = "rule-a"
regex = 'alpha'
`)
	baseB := write("base-b.toml", `
title = "base-b"
[[rules]]
id = "rule-b"
regex = 'beta'
`)
	mainA := write("main-a.toml", `
title = "main-a"
[extend]
path = "`+filepath.ToSlash(baseA)+`"
`)
	mainB := write("main-b.toml", `
title = "main-b"
[extend]
path = "`+filepath.ToSlash(baseB)+`"
`)

	type result struct {
		cfg Config
		err error
	}
	load := func(path string) result {
		v := viper.New()
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return result{err: err}
		}
		var vc ViperConfig
		if err := v.Unmarshal(&vc); err != nil {
			return result{err: err}
		}
		vc.ConfigPath = path
		cfg, err := vc.Translate()
		return result{cfg: cfg, err: err}
	}

	var wg sync.WaitGroup
	results := make([]result, 2)
	for i, path := range []string{mainA, mainB} {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			results[i] = load(path)
		}(i, path)
	}
	wg.Wait()

	require.NoError(t, results[0].err)
	require.NoError(t, results[1].err)
	assert.Equal(t, "main-a", results[0].cfg.Title)
	assert.Equal(t, "main-b", results[1].cfg.Title)
	require.Contains(t, results[0].cfg.Rules, "rule-a")
	assert.NotContains(t, results[0].cfg.Rules, "rule-b")
	require.Contains(t, results[1].cfg.Rules, "rule-b")
	assert.NotContains(t, results[1].cfg.Rules, "rule-a")
	assert.Equal(t, mainA, results[0].cfg.Path)
	assert.Equal(t, mainB, results[1].cfg.Path)
}
