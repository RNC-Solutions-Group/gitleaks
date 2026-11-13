package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const informativeTestRule = `
[[rules]]
id = "cpv7a-test-rule"
description = "CP-V7A test credential"
regex = '''CPV7A_[A-Z0-9]{8}'''
keywords = ["cpv7a_"]
`

func informativeConfig(informative, consoleOnly bool, extra string) string {
	return fmt.Sprintf(`title = "CP-V7A informative test"
informative-report = %t
console-only = %t
%s
%s`, informative, consoleOnly, extra, informativeTestRule)
}

func writeInformativeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gitleaks.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func decodeInformativeReport(t *testing.T, path string) map[string]any {
	t.Helper()
	var document map[string]any
	require.NoError(t, json.Unmarshal(mustReadFile(t, path), &document))
	return document
}

func informativeSummary(document map[string]any) map[string]any {
	return document["summary"].(map[string]any)
}

func informativeExclusions(document map[string]any) []any {
	return document["exclusions"].([]any)
}

func TestCLIInformativeOutputMatrix(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("disabled no report path keeps original console behavior", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(false, false, ""))
		result := runReportCLI(t, "dir", source, "--no-banner", "--config", configPath, "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.NotContains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
		assert.Empty(t, filesWithSuffix(t, source, "-informative.json"))
	})

	t.Run("disabled file report has no companion", func(t *testing.T) {
		source := cleanScanDir(t)
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		configPath := writeInformativeConfig(t, informativeConfig(false, false, ""))
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
		assert.NoFileExists(t, filepath.Join(dir, "report-informative.json"))
		assert.NotContains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("enabled no report path writes console only", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
		result := runReportCLI(t, "dir", source, "--no-banner", "--config", configPath, "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
		assert.Contains(t, result.stderr, "Excluded targets:")
		assert.Contains(t, result.stderr, "Generated at (UTC):")
		assert.Empty(t, filesWithSuffix(t, source, "-informative.json"))
	})

	t.Run("enabled file report writes a paired JSON document", func(t *testing.T) {
		source := cleanScanDir(t)
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
		start := time.Now().UTC().Add(-time.Second)
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		end := time.Now().UTC().Add(time.Second)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
		companion := filepath.Join(dir, "report-informative.json")
		document := decodeInformativeReport(t, companion)
		assert.EqualValues(t, 1, document["schema_version"])
		scan := document["scan"].(map[string]any)
		assert.Equal(t, "directory", scan["source_type"])
		assert.Equal(t, true, scan["completed"])
		assert.Equal(t, dest, scan["primary_report"])
		generatedAt, ok := scan["generated_at"].(string)
		require.True(t, ok, "informative scan must have a generated_at string")
		parsed, err := time.Parse(time.RFC3339, generatedAt)
		require.NoError(t, err)
		assert.False(t, parsed.Before(start))
		assert.False(t, parsed.After(end))
		assert.NotContains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
		assert.Contains(t, result.stderr, "informative report written")
	})

	t.Run("single file source type", func(t *testing.T) {
		sourceDir := cleanScanDir(t)
		source := filepath.Join(sourceDir, "readme.txt")
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		document := decodeInformativeReport(t, informativeCompanionPath(dest))
		assert.Equal(t, "file", document["scan"].(map[string]any)["source_type"])
	})

	t.Run("console only with file writes no companion", func(t *testing.T) {
		source := cleanScanDir(t)
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		configPath := writeInformativeConfig(t, informativeConfig(true, true, ""))
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
		assert.NoFileExists(t, filepath.Join(dir, "report-informative.json"))
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("console only without a report path is unchanged", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(true, true, ""))
		result := runReportCLI(t, "dir", source, "--no-banner", "--config", configPath, "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
		assert.Contains(t, result.stderr, "Generated at (UTC):")
		assert.Empty(t, filesWithSuffix(t, source, "-informative.json"))
	})

	t.Run("stdout report is not contaminated", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
		result := runReportCLI(t,
			"dir", source,
			"--no-banner",
			"--config", configPath,
			"--report-path", "-",
			"--report-format", "json",
			"--log-level", "error",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", result.stdout)
		assert.NotContains(t, result.stdout, "INFORMATIVE")
		assert.NotContains(t, result.stdout, "generated_at")
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
		assert.Empty(t, filesWithSuffix(t, source, "-informative.json"))
	})

	t.Run("console only with stdout keeps the primary on stdout", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(true, true, ""))
		result := runReportCLI(t,
			"dir", source,
			"--no-banner",
			"--config", configPath,
			"--report-path", "-",
			"--report-format", "json",
			"--log-level", "error",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", result.stdout)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("console only is inert when informative reporting is disabled", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := writeInformativeConfig(t, informativeConfig(false, true, ""))
		result := runReportCLI(t, "dir", source, "--no-banner", "--config", configPath, "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.NotContains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})
}

func TestCLIInformativePairedFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	templatePath := filepath.Join(t.TempDir(), "report.tmpl")
	require.NoError(t, os.WriteFile(templatePath, []byte(`findings={{ len . }}`), 0o600))

	tests := []struct {
		name       string
		filename   string
		format     string
		extra      []string
		expectText string
	}{
		{name: "json", filename: "report.json", format: "json", expectText: "[]"},
		{name: "csv", filename: "report.csv", format: "csv", expectText: ""},
		{name: "sarif", filename: "report.sarif", format: "sarif", expectText: `"version"`},
		{name: "junit", filename: "report.xml", format: "junit", expectText: "<testsuites"},
		{name: "template", filename: "report.txt", format: "template", extra: []string{"--report-template", templatePath}, expectText: "findings=0"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := cleanScanDir(t)
			dir := t.TempDir()
			dest := filepath.Join(dir, test.filename)
			configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
			args := []string{
				"dir", source,
				"--no-banner",
				"--config", configPath,
				"--report-path", dest,
				"--report-format", test.format,
			}
			args = append(args, test.extra...)
			result := runReportCLI(t, args...)
			require.Equal(t, 0, result.exitCode, result.stderr)
			if test.expectText != "" {
				assert.Contains(t, string(mustReadFile(t, dest)), test.expectText)
			}
			companion := informativeCompanionPath(dest)
			document := decodeInformativeReport(t, companion)
			scan := document["scan"].(map[string]any)
			assert.Equal(t, dest, scan["primary_report"])
			assert.Equal(t, "directory", scan["source_type"])
			_, err := time.Parse(time.RFC3339, scan["generated_at"].(string))
			require.NoError(t, err)
		})
	}
}

func TestCLIInformativeAggregationPair(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := cleanScanDir(t)
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("SENTINEL-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("SENTINEL-INFORMATIVE"), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, "aggregation = true\nauto-overwrite = true\n"))

	for index, expected := range []struct {
		primary   string
		companion string
	}{
		{filepath.Join(dir, "report_001.json"), filepath.Join(dir, "report_001-informative.json")},
		{filepath.Join(dir, "report_002.json"), filepath.Join(dir, "report_002-informative.json")},
	} {
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, expected.primary)))
		document := decodeInformativeReport(t, expected.companion)
		scan := document["scan"].(map[string]any)
		assert.Equal(t, expected.primary, scan["primary_report"], "scan %d", index)
	}
	assert.Equal(t, "SENTINEL-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "SENTINEL-INFORMATIVE", string(mustReadFile(t, companion)))
}

func TestCLIInformativeExclusionCategories(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "readme.txt"), []byte("ordinary text\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "empty.txt"), nil, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "skip-me.txt"), []byte("CPV7A_SKIP0001\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "binary.dat"), []byte("%PDF-1.7\nCPV7A_BINARY01\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "large.txt"), append([]byte("CPV7A_LARGE001\n"), make([]byte, 1_000_010)...), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "archive.zip"), zipFixture(t, "inside.txt", "CPV7A_ARCHIVE1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "report-informative.json"), []byte("CPV7A_OLDFILE1\n"), 0o600))

	configPath := writeInformativeConfig(t, informativeConfig(true, false, `
[[allowlists]]
paths = ['''skip-me\.txt''']
`))
	dest := filepath.Join(source, "report.json")
	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--max-archive-depth", "0",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	companion := informativeCompanionPath(dest)
	document := decodeInformativeReport(t, companion)
	summary := informativeSummary(document)
	assert.EqualValues(t, 7, summary["observed_targets"])
	assert.EqualValues(t, 1, summary["inspected_targets"])
	assert.EqualValues(t, 6, summary["excluded_targets"])
	reasons := summary["exclusions_by_reason"].(map[string]any)
	assert.EqualValues(t, 1, reasons["empty_file"])
	assert.EqualValues(t, 1, reasons["path_allowlist"])
	assert.EqualValues(t, 1, reasons["size_limit"])
	assert.EqualValues(t, 1, reasons["binary_content"])
	assert.EqualValues(t, 1, reasons["archive_depth"])
	assert.EqualValues(t, 1, reasons["excluded_input"])

	allExclusions, err := json.Marshal(informativeExclusions(document))
	require.NoError(t, err)
	encoded := string(allExclusions)
	assert.NotContains(t, encoded, "CPV7A_SKIP0001")
	assert.NotContains(t, encoded, "CPV7A_BINARY01")
	assert.NotContains(t, encoded, "CPV7A_LARGE001")
	assert.NotContains(t, encoded, "CPV7A_ARCHIVE1")
	assert.NotContains(t, encoded, "CPV7A_OLDFILE1")
	assert.NotContains(t, result.stderr, "CPV7A_SKIP0001")
}

func TestCLIInformativeArchiveMemberSizeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "big-member.zip"),
		zipFixture(t, "big.txt", "CPV7A_MEMBER01\n"+strings.Repeat("x", 1_000_000)), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--max-archive-depth", "1",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	exclusions := informativeExclusions(document)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "archive_member", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
	assert.Contains(t, exclusion["path"], "!big.txt")
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_MEMBER01")
}

func TestCLIInformativeStdinSizeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	input := []byte(strings.Repeat("A", 1_000_001) + "CPV7A_STDIN001\n")
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runStdinCLI(t, input,
		"--no-banner",
		"--config", configPath,
		"--max-target-megabytes", "1",
		"--report-path", dest,
		"--report-format", "json",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	scan := document["scan"].(map[string]any)
	assert.Equal(t, "stdin", scan["source_type"])
	summary := informativeSummary(document)
	assert.EqualValues(t, 1, summary["observed_targets"])
	assert.EqualValues(t, 0, summary["inspected_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
	exclusions := informativeExclusions(document)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "stdin", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
	_, hasSize := exclusion["size_bytes"]
	assert.False(t, hasSize, "a bounded stream limit is not an exact size")
	assert.EqualValues(t, 1_000_000, exclusion["limit_bytes"])
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_STDIN001")
}

func TestCLIInformativeNonseekableArchiveLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	input := zipStoredFixture(t, "inside.txt", strings.Repeat("x", 1_000_100))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runStdinCLI(t, input,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--max-archive-depth", "1",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	summary := informativeSummary(document)
	assert.EqualValues(t, 1, summary["observed_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
	assert.EqualValues(t, 1, summary["exclusions_by_reason"].(map[string]any)["size_limit"])
	exclusions := informativeExclusions(document)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "stdin", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
	_, hasSize := exclusion["size_bytes"]
	assert.False(t, hasSize, "a bounded nonseekable copy is not an exact size")
}

func TestCLIInformativeZeroFindingsCanStillHaveExclusions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "readme.txt"), []byte("ordinary\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, "empty.txt"), nil, 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	summary := document["summary"].(map[string]any)
	assert.EqualValues(t, 1, summary["inspected_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
}

func TestCLIInformativeCompanionDoesNotDiscloseFindings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	const secret = "CPV7A_SECRET01"
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "finding.txt"),
		[]byte("token = "+secret+"\n"), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
	require.Equal(t, 1, result.exitCode, result.stderr)
	assert.Contains(t, string(mustReadFile(t, dest)), secret)
	companion := string(mustReadFile(t, informativeCompanionPath(dest)))
	assert.NotContains(t, companion, secret)
	assert.NotContains(t, result.stderr, secret)
}

func TestCLIInformativeSuppressedFindingIsNotAnExcludedTarget(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "allowed.txt"),
		[]byte("token = CPV7A_ALLOW001 # gitleaks:allow\n"), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	summary := informativeSummary(document)
	assert.EqualValues(t, 1, summary["observed_targets"])
	assert.EqualValues(t, 1, summary["inspected_targets"])
	assert.EqualValues(t, 0, summary["excluded_targets"])
	assert.Empty(t, informativeExclusions(document))
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_ALLOW001")
}

func TestCLIInformativeConfigurationSources(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	content := informativeConfig(true, true, "")

	t.Run("explicit config", func(t *testing.T) {
		configPath := writeInformativeConfig(t, content)
		source := cleanScanDir(t)
		result := runReportCLI(t, "dir", source, "--no-banner", "--config", configPath, "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("environment config path", func(t *testing.T) {
		configPath := writeInformativeConfig(t, content)
		source := cleanScanDir(t)
		result := runReportCLIWithEnv(t, map[string]string{"GITLEAKS_CONFIG": configPath},
			"dir", source, "--no-banner", "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("environment config content", func(t *testing.T) {
		source := cleanScanDir(t)
		result := runReportCLIWithEnv(t, map[string]string{"GITLEAKS_CONFIG_TOML": content},
			"dir", source, "--no-banner", "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})

	t.Run("source config", func(t *testing.T) {
		source := cleanScanDir(t)
		require.NoError(t, os.WriteFile(filepath.Join(source, ".gitleaks.toml"), []byte(content), 0o600))
		result := runReportCLI(t, "dir", source, "--no-banner", "--log-level", "error")
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "INFORMATIVE SCAN SUMMARY")
	})
}

func TestCLIInformativeSymlinkExclusion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation generally requires elevated privileges on Windows")
	}

	source := t.TempDir()
	target := filepath.Join(source, "target.txt")
	link := filepath.Join(source, "link.txt")
	require.NoError(t, os.WriteFile(target, []byte("ordinary\n"), 0o600))
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	reasons := informativeSummary(document)["exclusions_by_reason"].(map[string]any)
	assert.EqualValues(t, 1, reasons["symlink_not_followed"])
	assert.NotContains(t, result.stderr, "skipping symlink")
}

func TestCLIInformativeIncompleteScanDoesNotPublish(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	// A malformed ZIP is a genuine archive identification/opening failure, not
	// an intentional exclusion.
	require.NoError(t, os.WriteFile(filepath.Join(source, "broken.zip"),
		[]byte("PK\x03\x04this is not a valid zip payload"), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--exit-code", "0",
		"--max-archive-depth", "1",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.NoFileExists(t, dest)
	assert.NoFileExists(t, informativeCompanionPath(dest))
	requireNoTemporaryReports(t, dir)
}

func TestCLIInformativeGitExclusions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	repo := initCLIStagedRepo(t)
	require.NoError(t, os.WriteFile(filepath.Join(repo, ".gitattributes"), []byte("binary.txt -diff\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "binary.txt"), []byte("CPV7A_GIT00001\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "archive.zip"), zipFixture(t, "inside.txt", "CPV7A_GIT00002\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "allowlisted.txt"), []byte("CPV7A_GIT00003\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "ordinary.txt"), []byte("ordinary text\n"), 0o600))
	commitAllCLIFiles(t, repo, "informative git fixture")

	configPath := writeInformativeConfig(t, informativeConfig(true, false, `
[[allowlists]]
paths = ['''allowlisted\.txt''']
`))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"git", repo,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-archive-depth", "0",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	scan := document["scan"].(map[string]any)
	assert.Equal(t, "git", scan["source_type"])
	reasons := informativeSummary(document)["exclusions_by_reason"].(map[string]any)
	assert.GreaterOrEqual(t, reasons["binary_content"], float64(1))
	assert.GreaterOrEqual(t, reasons["archive_depth"], float64(1))
	assert.GreaterOrEqual(t, reasons["path_allowlist"], float64(1))
	companion := string(mustReadFile(t, informativeCompanionPath(dest)))
	assert.NotContains(t, companion, "CPV7A_GIT00001")
	assert.NotContains(t, companion, "CPV7A_GIT00002")
	assert.NotContains(t, companion, "CPV7A_GIT00003")
}

func TestCLIInformativeGitStagedBlobSizeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	repo := initCLIStagedRepo(t)
	stageCLIFile(t, repo, "large.txt", "CPV7A_STAGE001\n\x00"+strings.Repeat("x", 1_000_000))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"git", repo,
		"--pre-commit", "--staged",
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	reasons := informativeSummary(document)["exclusions_by_reason"].(map[string]any)
	assert.GreaterOrEqual(t, reasons["size_limit"], float64(1))
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_STAGE001")
}

func commitAllCLIFiles(t *testing.T, repo, message string) {
	t.Helper()
	cmd := exec.Command("git", "-c", "core.fsmonitor=false", "-C", repo, "add", "-A")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
	cmd = exec.Command("git", "-c", "core.fsmonitor=false", "-C", repo, "commit", "-q", "-m", message)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}

func TestCLIInformativeCompressedStreamSizeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "large.txt.gz"),
		gzipFixture(t, "CPV7A_GZIP0001\n"+strings.Repeat("x", 1_000_000)), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--max-archive-depth", "1",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	summary := informativeSummary(document)
	assert.EqualValues(t, 1, summary["observed_targets"])
	assert.EqualValues(t, 0, summary["inspected_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
	exclusions := informativeExclusions(document)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "file", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
	_, hasSize := exclusion["size_bytes"]
	assert.False(t, hasSize, "bounded decompression does not prove an exact size")
	assert.Contains(t, exclusion["detail"], "stopped after")
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_GZIP0001")
}

func TestCLIInformativeNestedArchiveMemberPath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	inner := zipFixture(t, "big.txt", "CPV7A_NESTED01\n"+strings.Repeat("x", 1_000_000))
	require.NoError(t, os.WriteFile(filepath.Join(source, "outer.tar.gz"),
		tarGzFixture(t, "inner.zip", inner), 0o600))
	configPath := writeInformativeConfig(t, informativeConfig(true, false, ""))
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--config", configPath,
		"--report-path", dest,
		"--report-format", "json",
		"--max-target-megabytes", "1",
		"--max-archive-depth", "2",
		"--log-level", "error",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	document := decodeInformativeReport(t, informativeCompanionPath(dest))
	summary := informativeSummary(document)
	assert.EqualValues(t, 3, summary["observed_targets"])
	assert.EqualValues(t, 2, summary["inspected_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
	exclusions := informativeExclusions(document)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "archive_member", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
	assert.Contains(t, exclusion["path"], "outer.tar.gz!inner.zip!big.txt")
	assert.NotContains(t, string(mustReadFile(t, informativeCompanionPath(dest))), "CPV7A_NESTED01")
}

func gzipFixture(t *testing.T, content string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return buffer.Bytes()
}

func tarGzFixture(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWriter := tar.NewWriter(gzipWriter)
	require.NoError(t, tarWriter.WriteHeader(&tar.Header{
		Name: name,
		Mode: 0o600,
		Size: int64(len(content)),
	}))
	_, err := tarWriter.Write(content)
	require.NoError(t, err)
	require.NoError(t, tarWriter.Close())
	require.NoError(t, gzipWriter.Close())
	return buffer.Bytes()
}

func zipStoredFixture(t *testing.T, name, content string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	header := &zip.FileHeader{Name: name, Method: zip.Store}
	entry, err := writer.CreateHeader(header)
	require.NoError(t, err)
	_, err = entry.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return buffer.Bytes()
}

func zipFixture(t *testing.T, name, content string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	entry, err := writer.Create(name)
	require.NoError(t, err)
	_, err = entry.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return buffer.Bytes()
}

func filesWithSuffix(t *testing.T, dir, suffix string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	var matches []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), suffix) {
			matches = append(matches, filepath.Join(dir, entry.Name()))
		}
	}
	return matches
}
