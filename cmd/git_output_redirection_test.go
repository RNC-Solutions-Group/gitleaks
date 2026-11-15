package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIGitOutputRedirectionFailsBeforeGit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}
	repo := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
	}
	runGit("init", "-q")
	runGit("config", "user.name", "Test")
	runGit("config", "user.email", "test@example.com")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("ALPHA_12345678\n"), 0o600))
	runGit("add", ".")
	runGit("commit", "-q", "-m", "initial")

	cfg := redactionConfig(t, "")
	for _, logOpts := range []string{
		"--all --output=" + filepath.Join(t.TempDir(), "leak.txt"),
		"--all --output " + filepath.Join(t.TempDir(), "leak-separated.txt"),
	} {
		outputPath := logOpts[strings.LastIndex(logOpts, " ")+1:]
		if idx := strings.Index(outputPath, "="); idx >= 0 {
			outputPath = outputPath[idx+1:]
		}
		result := runReportCLI(t,
			"git", repo,
			"--config", cfg,
			"--log-opts", logOpts,
			"--exit-code", "0",
			"--no-banner",
			"--log-level", "error",
		)
		require.NotEqual(t, 0, result.exitCode, result.stderr)
		assert.Contains(t, result.stderr, "output redirection")
		assert.NoFileExists(t, outputPath)
	}
}
