package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const scanStatusSecret = "aws_access_key_id = \"AKIAABCDEFGHIJKLMNOP\"\n"

func initScanStatusRepo(t *testing.T, withSecret bool) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available")
	}

	repo := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		commandArgs := append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)
		cmd := exec.Command(gitPath, commandArgs...)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output
		require.NoError(t, cmd.Run(), "git %v: %s", args, output.String())
	}
	runGit("init", "-q")
	runGit("config", "core.fsmonitor", "false")
	runGit("config", "user.email", "cpv3@example.com")
	runGit("config", "user.name", "CPV3")
	runGit("config", "core.autocrlf", "false")

	name := "readme.txt"
	content := "ordinary repository content\n"
	if withSecret {
		name = "creds.txt"
		content = scanStatusSecret
	}
	require.NoError(t, os.WriteFile(filepath.Join(repo, name), []byte(content), 0o600))
	runGit("add", name)
	runGit("commit", "-q", "-m", "fixture")
	return repo
}

func TestCLIGitOperationalFailureStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("non-repository fails without a report", func(t *testing.T) {
		source := t.TempDir()
		reportPath := filepath.Join(source, "report.json")
		result := runReportCLI(t,
			"git", source,
			"--no-banner",
			"--log-level", "error",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.NotEqual(t, 0, result.exitCode)
		assert.NoFileExists(t, reportPath)
		assert.Contains(t, result.stderr, "failed to scan Git repository")
	})

	t.Run("invalid log range fails and preserves the report", func(t *testing.T) {
		repo := initScanStatusRepo(t, true)
		reportPath := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

		result := runReportCLI(t,
			"git", repo,
			"--no-banner",
			"--log-level", "error",
			"--log-opts", "does-not-exist..HEAD",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.NotEqual(t, 0, result.exitCode)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
		assert.Contains(t, result.stderr, "failed to scan Git repository")
	})

	t.Run("exit-code zero cannot mask Git failure", func(t *testing.T) {
		repo := initScanStatusRepo(t, true)
		result := runReportCLI(t,
			"git", repo,
			"--no-banner",
			"--log-level", "error",
			"--exit-code", "0",
			"--log-opts", "does-not-exist..HEAD",
		)
		require.NotEqual(t, 0, result.exitCode)
		assert.Contains(t, result.stderr, "failed to scan Git repository")
	})

	t.Run("successful finding scan honors exit-code zero", func(t *testing.T) {
		repo := initScanStatusRepo(t, true)
		reportPath := filepath.Join(t.TempDir(), "report.json")
		result := runReportCLI(t,
			"git", repo,
			"--no-banner",
			"--log-level", "error",
			"--exit-code", "0",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, string(mustReadFile(t, reportPath)), `"RuleID"`)
	})

	t.Run("successful clean scan stays successful", func(t *testing.T) {
		repo := initScanStatusRepo(t, false)
		reportPath := filepath.Join(t.TempDir(), "report.json")
		result := runReportCLI(t,
			"git", repo,
			"--no-banner",
			"--log-level", "error",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
	})
}

func TestCLIDirectoryFailureStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("clean and finding controls", func(t *testing.T) {
		clean := cleanScanDir(t)
		result := runReportCLI(t, reportCommandArgs(clean, filepath.Join(t.TempDir(), "clean.json"), "")...)
		require.Equal(t, 0, result.exitCode, result.stderr)

		finding := cleanScanDir(t)
		require.NoError(t, os.WriteFile(filepath.Join(finding, "creds.txt"), []byte(scanStatusSecret), 0o600))
		reportPath := filepath.Join(t.TempDir(), "finding.json")
		args := reportCommandArgs(finding, reportPath, "")
		args = append(args, "--exit-code", "0")
		result = runReportCLI(t, args...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, string(mustReadFile(t, reportPath)), `"RuleID"`)
	})

	t.Run("unreadable file fails despite exit-code zero and preserves report", func(t *testing.T) {
		source := cleanScanDir(t)
		require.NoError(t, os.WriteFile(filepath.Join(source, "readable.txt"), []byte(scanStatusSecret), 0o600))
		locked := filepath.Join(source, "locked.txt")
		require.NoError(t, os.WriteFile(locked, []byte("ordinary content\n"), 0o600))
		lockScanFile(t, locked)

		reportPath := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))
		args := reportCommandArgs(source, reportPath, "")
		args = append(args, "--exit-code", "0")
		result := runReportCLI(t, args...)

		require.NotEqual(t, 0, result.exitCode)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
		assert.Contains(t, result.stderr, "partial scan completed")
		assert.Contains(t, result.stderr, "locked.txt")
	})

	t.Run("unreadable file without findings is an operational failure", func(t *testing.T) {
		source := cleanScanDir(t)
		locked := filepath.Join(source, "locked.txt")
		require.NoError(t, os.WriteFile(locked, []byte("ordinary content\n"), 0o600))
		lockScanFile(t, locked)

		reportPath := filepath.Join(t.TempDir(), "report.json")
		args := reportCommandArgs(source, reportPath, "")
		args = append(args, "--exit-code", "0")
		result := runReportCLI(t, args...)

		require.NotEqual(t, 0, result.exitCode)
		assert.NoFileExists(t, reportPath)
		assert.Contains(t, result.stderr, "partial scan completed")
	})
}
