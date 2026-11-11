package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cliStagedPAT = "ghp_Zq8xR2mV7kLp4Wn9Tb3Yc6Hd1Fs5Ga0Ju2Ke"

func initCLIStagedRepo(t *testing.T) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available")
	}

	repo := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(gitPath, append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)...)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v: %s", args, output)
	}
	run("init", "-q")
	run("config", "core.fsmonitor", "false")
	run("config", "user.email", "cpv4@example.com")
	run("config", "user.name", "CPV4")
	run("commit", "-q", "--allow-empty", "-m", "root")
	return repo
}

func stageCLIFile(t *testing.T, repo, name, content string) {
	t.Helper()
	path := filepath.Join(repo, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	cmd := exec.Command("git", "-c", "core.fsmonitor=false", "-C", repo, "add", "--", name)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}

func TestCLIStagedMarkedBinaryIsReported(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("gitattributes", func(t *testing.T) {
		repo := initCLIStagedRepo(t)
		stageCLIFile(t, repo, ".gitattributes", "secret.txt -diff\n")
		stageCLIFile(t, repo, "secret.txt", "token = \""+cliStagedPAT+"\"\n")
		reportPath := filepath.Join(t.TempDir(), "report.json")

		result := runReportCLI(t,
			"git", repo, "--pre-commit", "--staged",
			"--no-banner", "--log-level", "error",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.Equal(t, 1, result.exitCode)
		assert.Contains(t, string(mustReadFile(t, reportPath)), cliStagedPAT)

		result = runReportCLI(t,
			"git", repo, "--pre-commit", "--staged",
			"--no-banner", "--log-level", "error",
			"--exit-code", "0",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.Equal(t, 0, result.exitCode)
		assert.Contains(t, string(mustReadFile(t, reportPath)), cliStagedPAT)
	})

	t.Run("worktree divergence", func(t *testing.T) {
		repo := initCLIStagedRepo(t)
		stageCLIFile(t, repo, "secret.txt", "token = \""+cliStagedPAT+"\"\n")
		require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("clean worktree\n"), 0o600))
		reportPath := filepath.Join(t.TempDir(), "report.json")

		result := runReportCLI(t,
			"git", repo, "--pre-commit", "--staged",
			"--no-banner", "--log-level", "error",
			"--report-path", reportPath,
			"--report-format", "json",
		)
		require.Equal(t, 1, result.exitCode)
		assert.Contains(t, string(mustReadFile(t, reportPath)), cliStagedPAT)
	})
}

func TestCLIStagedActualBinaryRemainsExcluded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	repo := initCLIStagedRepo(t)
	binary := append([]byte("\x89PNG\r\n\x1a\n\x00"), []byte("token = \""+cliStagedPAT+"\"\n")...)
	stageCLIFile(t, repo, "image.png", string(binary))
	reportPath := filepath.Join(t.TempDir(), "report.json")

	result := runReportCLI(t,
		"git", repo, "--pre-commit", "--staged",
		"--no-banner", "--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, reportPath)))
}
