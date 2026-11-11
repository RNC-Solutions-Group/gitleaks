package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	redactionAlpha = "ALPHA_12345678"
	redactionBeta  = "BETA_87654321"
	redactionAux   = "AUX_87654321"
)

func redactionConfig(t *testing.T, extra string) string {
	t.Helper()
	content := `title = "cpv5-redaction"
[[rules]]
id = "alpha"
description = "alpha"
regex = '''ALPHA_[A-Z0-9]{8}'''
[[rules]]
id = "beta"
description = "beta"
regex = '''BETA_[A-Z0-9]{8}'''
` + extra
	path := filepath.Join(t.TempDir(), "gitleaks.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func redactionSource(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	}
	return dir
}

func TestCLIRedactionAdjacentSecrets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, "")
	source := redactionSource(t, map[string]string{
		"adjacent.txt": "left=" + redactionAlpha + " right=" + redactionBeta + "\n",
	})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--verbose",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.NotContains(t, result.stdout, redactionAlpha)
	assert.NotContains(t, result.stdout, redactionBeta)
	assert.NotContains(t, string(mustReadFile(t, reportPath)), redactionAlpha)
	assert.NotContains(t, string(mustReadFile(t, reportPath)), redactionBeta)
	assert.Contains(t, result.stdout, "REDACTED")
}

func TestCLIRedactionDefaultReportRemainsRaw(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, "")
	source := redactionSource(t, map[string]string{"secret.txt": redactionAlpha + "\n"})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.Contains(t, string(mustReadFile(t, reportPath)), redactionAlpha)
}

func TestCLIRedactionEmptySecret(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, `
[[rules]]
id = "empty"
description = "empty capture"
regex = '''EMPTYTOKEN=()([A-Za-z0-9]+)'''
secretGroup = 1
`)
	source := redactionSource(t, map[string]string{"empty.txt": "EMPTYTOKEN=SUPERSECRET\n"})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--verbose",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.NotContains(t, result.stdout, "SUPERSECRET")
	assert.NotContains(t, result.stdout, "REDACTEDEREDACTED")
	assert.NotContains(t, string(mustReadFile(t, reportPath)), "SUPERSECRET")
	assert.Contains(t, result.stdout, "REDACTED")
}

func TestCLIRedactionPathOnlyDoesNotInsertMaskBetweenCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, `
[[rules]]
id = "path-only"
description = "path only"
path = '''path-only\.txt$'''
`)
	source := redactionSource(t, map[string]string{"path-only.txt": "ordinary content\n"})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--verbose",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.Contains(t, result.stdout, "file detected:")
	assert.NotContains(t, result.stdout, "REDACTEDfREDACTED")
	assert.NotContains(t, result.stdout, "REDACTEDpREDACTED")
	assert.Contains(t, string(mustReadFile(t, reportPath)), "file detected:")
}

func TestCLIRedactionRequiredFindings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, `
[[rules]]
id = "primary"
description = "primary"
regex = '''PRIMARY_[A-Z0-9]{8}'''
  [[rules.required]]
  id = "aux"
  withinLines = 2
[[rules]]
id = "aux"
description = "aux"
regex = '''AUX_[A-Z0-9]{8}'''
`)
	source := redactionSource(t, map[string]string{"required.txt": "PRIMARY_12345678\n" + redactionAux + "\n"})

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--verbose",
		"--no-banner",
		"--log-level", "error",
	)
	require.Equal(t, 1, result.exitCode)
	assert.NotContains(t, result.stdout, redactionAux)
	assert.Contains(t, result.stdout, "Required:")
	assert.Contains(t, result.stdout, "REDACTED")
}

func TestCLIRedactionUnicodePartial(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	secret := "pässwörd🔑"
	cfg := redactionConfig(t, `
[[rules]]
id = "unicode"
description = "unicode"
regex = '''pässwörd🔑'''
`)
	source := redactionSource(t, map[string]string{"unicode.txt": "token=" + secret + "\n"})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact=45",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	reportBytes := mustReadFile(t, reportPath)
	assert.True(t, utf8.Valid(reportBytes))
	assert.NotContains(t, string(reportBytes), secret)
	assert.NotContains(t, string(reportBytes), "\ufffd")

	var findings []map[string]any
	require.NoError(t, json.Unmarshal(reportBytes, &findings))
	require.NotEmpty(t, findings)
	for _, field := range []string{"Secret", "Match"} {
		value, ok := findings[0][field].(string)
		require.True(t, ok)
		assert.True(t, utf8.ValidString(value), field)
		assert.NotContains(t, value, secret)
	}

	result = runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact=45",
		"--verbose",
		"--no-banner",
		"--log-level", "error",
	)
	require.Equal(t, 1, result.exitCode)
	assert.True(t, utf8.ValidString(result.stdout))
	assert.NotContains(t, result.stdout, secret)
	assert.NotContains(t, result.stdout, "\ufffd")
}

func TestCLIRedactionAllReportFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, "")
	source := redactionSource(t, map[string]string{"secret.txt": redactionAlpha + "\n"})
	templatePath := filepath.Join(t.TempDir(), "report.tmpl")
	require.NoError(t, os.WriteFile(templatePath, []byte(`{{range .}}{{.Secret}}|{{.Match}}|{{.Line}}{{end}}`), 0o600))

	for _, format := range []string{"json", "csv", "junit", "sarif", "template"} {
		t.Run(format, func(t *testing.T) {
			reportPath := filepath.Join(t.TempDir(), "report."+format)
			args := []string{
				"dir", source,
				"--config", cfg,
				"--redact",
				"--no-banner",
				"--log-level", "error",
				"--report-path", reportPath,
				"--report-format", format,
			}
			if format == "template" {
				args = append(args, "--report-template", templatePath)
			}
			result := runReportCLI(t, args...)
			require.Equal(t, 1, result.exitCode, result.stderr)
			reportBytes := mustReadFile(t, reportPath)
			assert.NotContains(t, string(reportBytes), redactionAlpha)
			assert.True(t, utf8.Valid(reportBytes))
		})
	}
}

func TestCLIRedactionLogsDoNotExposeSecrets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, `
[[rules]]
id = "trace"
description = "trace"
regex = '''TRACE_[A-Z0-9]{8}'''
[[rules]]
id = "low-entropy"
description = "low entropy"
regex = '''LOWENTROPY_[A-Z0-9]{8}'''
entropy = 10
[[rules]]
id = "allowed"
description = "allowed"
regex = '''ALLOWED_[A-Z0-9]{8}'''
[allowlist]
description = "synthetic allowlist"
regexes = ['''ALLOWED_[A-Z0-9]{8}''']
`)
	source := redactionSource(t, map[string]string{
		"trace.txt":   "TRACE_12345678 # gitleaks:allow\n",
		"encoded.txt": "QUxQSEFfMTIzNDU2Nzg=\n",
		"allowed.txt": "ALLOWED_12345678\n",
		"entropy.txt": "LOWENTROPY_12345678\n",
	})
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--no-banner",
		"--log-level", "trace",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.NotContains(t, result.stderr, "TRACE_12345678")
	assert.NotContains(t, result.stderr, redactionAlpha)
	assert.NotContains(t, result.stderr, "QUxQSEFfMTIzNDU2Nzg=")
	assert.NotContains(t, result.stderr, "ALLOWED_12345678")
	assert.NotContains(t, result.stderr, "LOWENTROPY_12345678")
	assert.Contains(t, result.stderr, "encoded segment found")
	assert.Contains(t, result.stderr, "gitleaks:allow")
	assert.Contains(t, result.stderr, "allowlist")
	assert.Contains(t, result.stderr, "low entropy")
}

func TestCLIConfigTOMLIsNotLogged(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := redactionSource(t, map[string]string{"readme.txt": "ordinary content\n"})
	valid := `title = "SUPERSECRETCONFIG"
[[rules]]
id = "alpha"
description = "SUPERSECRETCONFIG"
regex = '''ALPHA_[A-Z0-9]{8}'''
`
	result := runReportCLIWithEnv(t, map[string]string{"GITLEAKS_CONFIG_TOML": valid},
		"dir", source,
		"--redact",
		"--no-banner",
		"--log-level", "debug",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.NotContains(t, result.stderr, "SUPERSECRETCONFIG")
	assert.NotContains(t, result.stderr, "content=")
	assert.Contains(t, result.stderr, "using gitleaks config from GITLEAKS_CONFIG_TOML")

	invalid := "title = \"SUPERSECRETCONFIG\"\n[[rules]\n"
	result = runReportCLIWithEnv(t, map[string]string{"GITLEAKS_CONFIG_TOML": invalid},
		"dir", source,
		"--redact",
		"--no-banner",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.NotContains(t, result.stderr, "SUPERSECRETCONFIG")
	assert.NotContains(t, result.stderr, "content=")
	assert.Contains(t, result.stderr, "unable to load gitleaks config")

	badConfig := filepath.Join(t.TempDir(), "bad.toml")
	require.NoError(t, os.WriteFile(badConfig, []byte("title = \"SUPERSECRETCONFIG\"\n[[rules]\n"), 0o600))
	result = runReportCLI(t,
		"dir", source,
		"--config", badConfig,
		"--redact",
		"--no-banner",
		"--log-level", "error",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.NotContains(t, result.stderr, "SUPERSECRETCONFIG")
	assert.NotContains(t, result.stderr, "content=")
}

func TestCLIRedactionStdoutReport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	cfg := redactionConfig(t, "")
	source := redactionSource(t, map[string]string{"secret.txt": redactionAlpha + "\n"})
	result := runReportCLI(t,
		"dir", source,
		"--config", cfg,
		"--redact",
		"--no-banner",
		"--log-level", "error",
		"--report-path", "-",
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.NotContains(t, result.stdout, redactionAlpha)
	assert.Contains(t, result.stdout, "REDACTED")
	assert.True(t, strings.Contains(result.stdout, "["), "stdout report should remain JSON")
}
