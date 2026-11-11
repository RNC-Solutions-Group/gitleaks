package cmd

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIArchiveLikePlaintextIsScanned(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "notes.7z.txt"), []byte("aws_access_key_id = \"AKIAABCDEFGHIJKLMNOP\"\n"), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	assert.Contains(t, string(mustReadFile(t, reportPath)), "AKIAABCDEFGHIJKLMNOP")
}

func TestCLIDirectoryMaxTargetSkipsOversizedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	content := strings.Repeat("A", 1_000_001) + `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"` + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(source, "big.txt"), []byte(content), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--max-target-megabytes", "1",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLIArchiveMemberMaxTargetSkipsMember(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	w, err := zw.Create("inside.txt")
	require.NoError(t, err)
	_, err = w.Write([]byte(strings.Repeat("A", 1_000_001) + `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	require.NoError(t, os.WriteFile(filepath.Join(source, "archive.zip"), archive.Bytes(), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--max-archive-depth", "1",
		"--max-target-megabytes", "1",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLICompressedExpansionMaxTargetSkipsStream(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	var compressed bytes.Buffer
	gw := gzip.NewWriter(&compressed)
	_, err := gw.Write([]byte(strings.Repeat("A", 1_000_001) + `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`))
	require.NoError(t, err)
	require.NoError(t, gw.Close())
	require.NoError(t, os.WriteFile(filepath.Join(source, "archive.gz"), compressed.Bytes(), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--max-archive-depth", "1",
		"--max-target-megabytes", "1",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func TestCLIStdinMaxTargetSkipsStream(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	reportPath := filepath.Join(t.TempDir(), "report.json")
	input := strings.Repeat("A", 1_000_001) + `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"` + "\n"
	result := runStdinCLI(t, []byte(input),
		"--max-target-megabytes", "1",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}

func runStdinCLI(t *testing.T, input []byte, args ...string) reportCLIResult {
	t.Helper()
	args = append([]string{"stdin"}, args...)
	encoded, err := json.Marshal(args)
	require.NoError(t, err)

	cmd := exec.Command(os.Args[0], "-test.run=^TestReportCLIChild$")
	env := make([]string, 0, len(os.Environ())+2)
	for _, value := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(value), "GITLEAKS_") {
			continue
		}
		env = append(env, value)
	}
	env = append(env, "GITLEAKS_REPORT_CHILD=1", "GITLEAKS_TEST_ARGS="+string(encoded))
	cmd.Env = env
	cmd.Stdin = bytes.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	result := reportCLIResult{stdout: stdout.String(), stderr: stderr.String()}
	if err == nil {
		return result
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result
	}
	t.Fatalf("could not run stdin CLI child: %v", err)
	return result
}

func TestCLICorruptArchiveIsOperationalFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "broken.zip"), []byte("PK\x03\x04truncated"), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	result := runReportCLI(t,
		"dir", source,
		"--max-archive-depth", "1",
		"--exit-code", "0",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "failed scan directory")
}

func TestCLIRealArchiveIsStillScanned(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	w, err := zw.Create("inside.txt")
	require.NoError(t, err)
	_, err = w.Write([]byte("aws_access_key_id = \"AKIAABCDEFGHIJKLMNOP\"\n"))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	require.NoError(t, os.WriteFile(filepath.Join(source, "archive.zip"), archive.Bytes(), 0o600))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"dir", source,
		"--max-archive-depth", "1",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 1, result.exitCode)
	report := string(mustReadFile(t, reportPath))
	assert.Contains(t, report, "AKIAABCDEFGHIJKLMNOP")
	assert.Contains(t, report, "archive.zip!inside.txt")
}
