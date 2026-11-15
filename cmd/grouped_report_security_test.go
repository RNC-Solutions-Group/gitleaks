package cmd

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const groupedTestKey = "0123456789abcdef0123456789abcdef"

func writeGroupedKeyFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "grouping.key")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func groupedSource(t *testing.T) string {
	t.Helper()
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	return source
}

func runGrouped(t *testing.T, args ...string) reportCLIResult {
	t.Helper()
	return runReportCLI(t, append([]string{"dir", groupedSource(t), "--no-banner", "--log-level", "error", "--exit-code", "0"}, args...)...)
}

func TestGroupedReportRejectsKeyAsDestination(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	before := mustReadFile(t, keyPath)

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", keyPath,
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "collides")
	assert.Equal(t, string(before), string(mustReadFile(t, keyPath)))
}

func TestGroupedReportRejectsKeyInsideScanTarget(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "a.txt"), []byte("ALPHA_12345678\n"), 0o600))
	keyPath := filepath.Join(source, "grouping.key")
	require.NoError(t, os.WriteFile(keyPath, []byte(groupedTestKey), 0o600))
	before := mustReadFile(t, keyPath)

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", filepath.Join(t.TempDir(), "grouped.json"),
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "inside the scan target")
	assert.Equal(t, string(before), string(mustReadFile(t, keyPath)))
}

func TestGroupedReportRejectsSymlinkKey(t *testing.T) {
	source := groupedSource(t)
	realKey := writeGroupedKeyFixture(t, groupedTestKey)
	link := filepath.Join(t.TempDir(), "grouping.link")
	if err := os.Symlink(realKey, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", filepath.Join(t.TempDir(), "grouped.json"),
		"--grouping-key-file", link,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "symbolic link")
}

func TestGroupedReportRejectsHardlinkKeyDestination(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	groupedPath := filepath.Join(t.TempDir(), "grouped.json")
	if err := os.Link(keyPath, groupedPath); err != nil {
		t.Skipf("hardlink unavailable: %v", err)
	}
	before := mustReadFile(t, keyPath)
	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", groupedPath,
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "collides")
	assert.Equal(t, string(before), string(mustReadFile(t, keyPath)))
}

func TestGroupedReportGeneratedKeyWorkflow(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	keyPath := writeGroupedKeyFixture(t, string(key))
	source := groupedSource(t)
	groupedPath := filepath.Join(t.TempDir(), "grouped.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", redactionConfig(t, ""),
		"--grouped-report", groupedPath,
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Contains(t, string(mustReadFile(t, groupedPath)), `"count": 1`)
}

func TestGroupedReportRejectsPermissiveKeyOnPOSIX(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX mode bits are not authoritative on Windows; ACL guidance applies")
	}
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	require.NoError(t, os.Chmod(keyPath, 0o644))
	result := runReportCLI(t,
		"dir", groupedSource(t),
		"--grouped-report", filepath.Join(t.TempDir(), "grouped.json"),
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "group/world accessible")
}

func groupedPolicyConfig(t *testing.T) string {
	t.Helper()
	content := `title = "grouped-policy"
auto-overwrite = false
[[rules]]
id = "alpha"
description = "alpha"
regex = '''ALPHA_[A-Z0-9]{8}'''
`
	path := filepath.Join(t.TempDir(), "policy.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestGroupedReportObeysOutputPolicy(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	cfg := groupedPolicyConfig(t)
	groupedPath := filepath.Join(t.TempDir(), "grouped.json")
	require.NoError(t, os.WriteFile(groupedPath, []byte("SENTINEL"), 0o600))

	base := []string{"dir", source, "--config", cfg, "--grouped-report", groupedPath, "--grouping-key-file", keyPath, "--no-banner", "--log-level", "error", "--exit-code", "0"}

	result := runReportCLI(t, base...)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, groupedPath)))

	result = runReportCLI(t, append(base, "--overwrite")...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.NotEqual(t, "SENTINEL", string(mustReadFile(t, groupedPath)))

	require.NoError(t, os.WriteFile(groupedPath, []byte("SENTINEL"), 0o600))
	result = runReportCLI(t, append(base, "--aggregate")...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, groupedPath)))
	assert.FileExists(t, filepath.Join(filepath.Dir(groupedPath), "grouped_001.json"))
}

func TestGroupedReportAllocatesBeyondV14ACeiling(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	dir := t.TempDir()
	groupedPath := filepath.Join(dir, "grouped.json")
	require.NoError(t, os.WriteFile(groupedPath, []byte("SENTINEL"), 0o600))
	seed := filepath.Join(dir, "occupied.seed")
	require.NoError(t, os.WriteFile(seed, []byte("OCCUPIED"), 0o600))
	for number := 1; number <= 1_001; number++ {
		candidate := numberedReportPath(groupedPath, number)
		if err := os.Link(seed, candidate); err != nil {
			require.NoError(t, os.WriteFile(candidate, []byte("OCCUPIED"), 0o600))
		}
	}

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", groupedPath,
		"--grouping-key-file", keyPath,
		"--config", groupedPolicyConfig(t),
		"--aggregate",
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, groupedPath)))
	allocated := numberedReportPath(groupedPath, 1_002)
	require.FileExists(t, allocated)
	assert.Contains(t, string(mustReadFile(t, allocated)), `"count": 1`)
	assert.NoFileExists(t, numberedReportPath(groupedPath, 1_003))
}

func TestGroupedReportRejectsAggregatePrimaryCollision(t *testing.T) {
	for _, number := range []int{1, 1_000, 10_001} {
		t.Run(fmt.Sprintf("candidate_%d", number), func(t *testing.T) {
			source := groupedSource(t)
			keyPath := writeGroupedKeyFixture(t, groupedTestKey)
			dir := t.TempDir()
			primary := filepath.Join(dir, "report.json")
			groupedPath := numberedReportPath(primary, number)

			result := runReportCLI(t,
				"dir", source,
				"--grouped-report", groupedPath,
				"--grouping-key-file", keyPath,
				"--report-path", primary,
				"--report-format", "json",
				"--aggregate",
				"--exit-code", "0", "--no-banner", "--log-level", "error",
			)
			require.NotEqual(t, 0, result.exitCode)
			assert.Contains(t, result.stderr, "aggregate")
			assert.NoFileExists(t, primary)
			assert.NoFileExists(t, groupedPath)
		})
	}
}

func TestGroupedReportRejectsKeyAsPrimaryDestination(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	before := mustReadFile(t, keyPath)
	result := runReportCLI(t,
		"dir", source,
		"--report-path", keyPath,
		"--report-format", "json",
		"--grouped-report", filepath.Join(t.TempDir(), "grouped.json"),
		"--grouping-key-file", keyPath,
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "collides")
	assert.Equal(t, string(before), string(mustReadFile(t, keyPath)))
}

func TestGroupedReportRejectsKeyAsAggregateCandidate(t *testing.T) {
	for _, number := range []int{1, 1_000, 10_001} {
		t.Run(fmt.Sprintf("candidate_%d", number), func(t *testing.T) {
			source := groupedSource(t)
			dir := t.TempDir()
			keyPath := numberedReportPath(filepath.Join(dir, "grouped.json"), number)
			require.NoError(t, os.WriteFile(keyPath, []byte(groupedTestKey), 0o600))
			groupedPath := filepath.Join(dir, "grouped.json")

			result := runReportCLI(t,
				"dir", source,
				"--grouped-report", groupedPath,
				"--grouping-key-file", keyPath,
				"--aggregate",
				"--exit-code", "0", "--no-banner", "--log-level", "error",
			)
			require.NotEqual(t, 0, result.exitCode)
			assert.Contains(t, result.stderr, "aggregate")
			assert.Equal(t, groupedTestKey, string(mustReadFile(t, keyPath)))
		})
	}
}

func TestGroupedReportRejectsConfigAsHighAggregateCandidate(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	dir := t.TempDir()
	primary := filepath.Join(dir, "report.json")
	configPath := numberedReportPath(primary, 10_001)
	configBody := "title = \"high\"\n[[rules]]\nid = \"alpha\"\ndescription = \"alpha\"\nregex = '''ALPHA_[A-Z0-9]{8}'''\n"
	require.NoError(t, os.WriteFile(configPath, []byte(configBody), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", filepath.Join(dir, "grouped.json"),
		"--grouping-key-file", keyPath,
		"--report-path", primary,
		"--report-format", "json",
		"--config", configPath,
		"--aggregate",
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "protected input")
	assert.Equal(t, configBody, string(mustReadFile(t, configPath)))
	assert.NoFileExists(t, primary)
}

func TestGroupedReportRejectsParentSymlinkCollision(t *testing.T) {
	source := groupedSource(t)
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	realDir := t.TempDir()
	linkDir := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(realDir, linkDir); err != nil {
		t.Skipf("directory symlink unavailable: %v", err)
	}
	primary := filepath.Join(realDir, "report.json")
	groupedPath := filepath.Join(linkDir, "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", groupedPath,
		"--grouping-key-file", keyPath,
		"--report-path", primary,
		"--report-format", "json",
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "collides")
	assert.NoFileExists(t, primary)
}

func groupedInformativeConfig(t *testing.T) string {
	t.Helper()
	content := `title = "grouped-informative"
informative-report = true
[[rules]]
id = "alpha"
description = "alpha"
regex = '''ALPHA_[A-Z0-9]{8}'''
`
	path := filepath.Join(t.TempDir(), "informative.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestGroupedReportRejectsInformativeCollision(t *testing.T) {
	for _, name := range []string{
		"report-informative.json",
		"report_001-informative.json",
		"report_1001-informative.json",
	} {
		t.Run(name, func(t *testing.T) {
			source := groupedSource(t)
			keyPath := writeGroupedKeyFixture(t, groupedTestKey)
			dir := t.TempDir()
			primary := filepath.Join(dir, "report.json")
			groupedPath := filepath.Join(dir, name)

			result := runReportCLI(t,
				"dir", source,
				"--grouped-report", groupedPath,
				"--grouping-key-file", keyPath,
				"--report-path", primary,
				"--report-format", "json",
				"--config", groupedInformativeConfig(t),
				"--aggregate",
				"--exit-code", "0", "--no-banner", "--log-level", "error",
			)
			require.NotEqual(t, 0, result.exitCode)
			assert.Contains(t, result.stderr, "informative")
			assert.NoFileExists(t, primary)
			assert.NoFileExists(t, groupedPath)
		})
	}
}

func TestGroupedReportScansLookalikeSourceFiles(t *testing.T) {
	for _, name := range []string{
		"grouped_001.json",
		"grouped_123.json",
		"grouped_999.json",
		"grouped_1000.json",
		"grouped_123.txt",
	} {
		t.Run(name, func(t *testing.T) {
			source := t.TempDir()
			keyPath := writeGroupedKeyFixture(t, groupedTestKey)
			groupedPath := filepath.Join(source, "grouped.json")
			require.NoError(t, os.WriteFile(filepath.Join(source, name), []byte("ALPHA_12345678\n"), 0o600))

			result := runReportCLI(t,
				"dir", source,
				"--grouped-report", groupedPath,
				"--grouping-key-file", keyPath,
				"--config", groupedPolicyConfig(t),
				"--overwrite",
				"--exit-code", "0", "--no-banner", "--log-level", "error",
			)
			require.Equal(t, 0, result.exitCode, result.stderr)

			var document struct {
				Groups []struct {
					Count int `json:"count"`
				} `json:"groups"`
			}
			raw := mustReadFile(t, groupedPath)
			require.NoError(t, json.Unmarshal(raw, &document))
			require.Len(t, document.Groups, 1, "lookalike source file must not be silently excluded")
			assert.Equal(t, 1, document.Groups[0].Count)
			assert.NotContains(t, string(raw), "ALPHA_12345678")
		})
	}
}

func TestGroupedReportExcludesExactConfiguredDestinationFromScan(t *testing.T) {
	source := t.TempDir()
	keyPath := writeGroupedKeyFixture(t, groupedTestKey)
	groupedPath := filepath.Join(source, "grouped.json")
	require.NoError(t, os.WriteFile(groupedPath, []byte("PREEXISTING ALPHA_12345678\n"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--grouped-report", groupedPath,
		"--grouping-key-file", keyPath,
		"--config", groupedPolicyConfig(t),
		"--overwrite",
		"--exit-code", "0", "--no-banner", "--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.NotContains(t, string(mustReadFile(t, groupedPath)), "ALPHA_12345678")
}
