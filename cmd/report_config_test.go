package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeConfigFixture(t *testing.T, path, content string) []byte {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return []byte(content)
}

func assertConfigCollision(t *testing.T, result reportCLIResult, path string, original []byte) {
	t.Helper()
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, filepath.Base(path))
	assert.Contains(t, result.stderr, "collides with scan input")
	assert.NotContains(t, result.stderr, "scanned ~")
	assert.Equal(t, original, mustReadFile(t, path))
}

func assertNoConfigSidecars(t *testing.T, dir, base string) {
	t.Helper()
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	matches, err := filepath.Glob(filepath.Join(dir, stem+"_*"+ext))
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestCLIExplicitConfigCollision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)
	configContent := "title = \"explicit-config\"\nauto-overwrite = true\naggregation = false\n"

	t.Run("absolute paths", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(t.TempDir(), "active.toml")
		original := writeConfigFixture(t, configPath, configContent)

		result := runReportCLI(t, reportCommandArgs(source, configPath, configPath)...)
		assertConfigCollision(t, result, configPath, original)
	})

	t.Run("relative paths", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(t.TempDir(), "active.toml")
		original := writeConfigFixture(t, configPath, configContent)
		relative, err := filepath.Rel(cwd, configPath)
		if err != nil {
			t.Skipf("could not form a relative path: %v", err)
		}

		result := runReportCLI(t, reportCommandArgs(source, relative, relative)...)
		assertConfigCollision(t, result, relative, original)
	})
}

func TestCLIAutoDiscoveredConfigCollision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)
	configContent := "title = \"auto-discovered-config\"\n"

	t.Run("absolute paths", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(source, ".gitleaks.toml")
		original := writeConfigFixture(t, configPath, configContent)

		result := runReportCLI(t, reportCommandArgs(source, configPath, "")...)
		assertConfigCollision(t, result, configPath, original)
		assertNoConfigSidecars(t, source, ".gitleaks.toml")
	})

	t.Run("relative report path", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(source, ".gitleaks.toml")
		original := writeConfigFixture(t, configPath, configContent)
		relative, err := filepath.Rel(cwd, configPath)
		if err != nil {
			t.Skipf("could not form a relative path: %v", err)
		}

		result := runReportCLI(t, reportCommandArgs(source, relative, "")...)
		assertConfigCollision(t, result, relative, original)
		assertNoConfigSidecars(t, source, ".gitleaks.toml")
	})
}

func TestCLIGitleaksConfigEnvCollision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)
	configContent := "title = \"environment-config\"\n"

	t.Run("absolute environment path", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(t.TempDir(), "active.toml")
		original := writeConfigFixture(t, configPath, configContent)

		result := runReportCLIWithEnv(t,
			map[string]string{"GITLEAKS_CONFIG": configPath},
			reportCommandArgs(source, configPath, "")...,
		)
		assertConfigCollision(t, result, configPath, original)
	})

	t.Run("relative environment path", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := filepath.Join(t.TempDir(), "active.toml")
		original := writeConfigFixture(t, configPath, configContent)
		relative, err := filepath.Rel(cwd, configPath)
		if err != nil {
			t.Skipf("could not form a relative path: %v", err)
		}

		result := runReportCLIWithEnv(t,
			map[string]string{"GITLEAKS_CONFIG": relative},
			reportCommandArgs(source, configPath, "")...,
		)
		assertConfigCollision(t, result, configPath, original)
	})
}

func TestCLIExtendedConfigCollision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("extended base", func(t *testing.T) {
		source := cleanScanDir(t)
		configDir := t.TempDir()
		basePath := filepath.Join(configDir, "base.toml")
		mainPath := filepath.Join(configDir, "main.toml")
		baseOriginal := writeConfigFixture(t, basePath, "title = \"base-config\"\n")
		mainOriginal := writeConfigFixture(t, mainPath,
			"title = \"main-config\"\n[extend]\npath = \""+filepath.ToSlash(basePath)+"\"\n")

		result := runReportCLI(t, reportCommandArgs(source, basePath, mainPath)...)
		assertConfigCollision(t, result, basePath, baseOriginal)
		assert.Equal(t, mainOriginal, mustReadFile(t, mainPath))
	})

	t.Run("nested extended base", func(t *testing.T) {
		source := cleanScanDir(t)
		configDir := t.TempDir()
		grandPath := filepath.Join(configDir, "grand.toml")
		basePath := filepath.Join(configDir, "base.toml")
		mainPath := filepath.Join(configDir, "main.toml")
		grandOriginal := writeConfigFixture(t, grandPath, "title = \"grand-config\"\n")
		baseOriginal := writeConfigFixture(t, basePath,
			"title = \"base-config\"\n[extend]\npath = \""+filepath.ToSlash(grandPath)+"\"\n")
		mainOriginal := writeConfigFixture(t, mainPath,
			"title = \"main-config\"\n[extend]\npath = \""+filepath.ToSlash(basePath)+"\"\n")

		result := runReportCLI(t, reportCommandArgs(source, grandPath, mainPath)...)
		assertConfigCollision(t, result, grandPath, grandOriginal)
		assert.Equal(t, baseOriginal, mustReadFile(t, basePath))
		assert.Equal(t, mainOriginal, mustReadFile(t, mainPath))
	})
}

func TestCLIConfigCollisionCannotBeOverridden(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := cleanScanDir(t)
	configPath := filepath.Join(source, ".gitleaks.toml")
	original := writeConfigFixture(t, configPath, "title = \"auto-discovered-config\"\n")

	for _, test := range []struct {
		name  string
		extra []string
	}{
		{name: "overwrite", extra: []string{"--overwrite"}},
		{name: "aggregate", extra: []string{"--aggregate"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := reportCommandArgs(source, configPath, "")
			args = append(args, test.extra...)
			result := runReportCLI(t, args...)
			assertConfigCollision(t, result, configPath, original)
			assertNoConfigSidecars(t, source, ".gitleaks.toml")
		})
	}
}

func TestCLIConfigTrackingAllowsOrdinaryReport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := cleanScanDir(t)
	configPath := filepath.Join(source, ".gitleaks.toml")
	original := writeConfigFixture(t, configPath, "title = \"auto-discovered-config\"\n")
	dest := filepath.Join(source, "report.json")

	result := runReportCLI(t, reportCommandArgs(source, dest, "")...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	assert.Equal(t, original, mustReadFile(t, configPath))
	assert.True(t, containsPath(result.stderr, dest), result.stderr)
}
