//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIGitStderrDoesNotExposeSecrets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	dir := t.TempDir()
	fakeGitDir := filepath.Join(dir, "bin")
	require.NoError(t, os.Mkdir(fakeGitDir, 0o700))
	script := "#!/bin/sh\necho \"fatal: unable to access 'https://user:ghp_SECRETTOKEN@example.com/repo?token=QUERYVALUE': denied\" >&2\nexit 128\n"
	require.NoError(t, os.WriteFile(filepath.Join(fakeGitDir, "git"), []byte(script), 0o700))

	source := t.TempDir()
	reportPath := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	result := runReportCLIWithEnv(t, map[string]string{
		"PATH": fakeGitDir + string(os.PathListSeparator) + os.Getenv("PATH"),
	},
		"git", source, "--platform", "none",
		"--redact",
		"--no-banner",
		"--log-level", "error",
		"--report-path", reportPath,
		"--report-format", "json",
	)

	require.NotEqual(t, 0, result.exitCode)
	assert.NotContains(t, result.stderr, "ghp_SECRETTOKEN")
	assert.NotContains(t, result.stderr, "QUERYVALUE")
	assert.Contains(t, result.stderr, "https://REDACTED@example.com")
	assert.Contains(t, result.stderr, "failed to scan Git repository")
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
}

func TestCLIStagedBlobReadHonorsTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	dir := t.TempDir()
	fakeGitDir := filepath.Join(dir, "bin")
	require.NoError(t, os.Mkdir(fakeGitDir, 0o700))
	diffPath := filepath.Join(dir, "diff.patch")
	require.NoError(t, os.WriteFile(diffPath, []byte(stagedBlobFailureDiff), 0o600))
	script := "#!/bin/sh\ncase \"$*\" in\n  *\" cat-file \"*) exec sleep 5 ;;\n  *\" diff \"*) cat \"$FAKE_GIT_DIFF\"; exit 0 ;;\n  *) exit 128 ;;\nesac\n"
	require.NoError(t, os.WriteFile(filepath.Join(fakeGitDir, "git"), []byte(script), 0o700))

	source := t.TempDir()
	reportPath := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	start := time.Now()
	result := runReportCLIWithEnv(t, map[string]string{
		"PATH":          fakeGitDir + string(os.PathListSeparator) + os.Getenv("PATH"),
		"FAKE_GIT_DIFF": diffPath,
	},
		"git", source, "--staged", "--platform", "none",
		"--no-banner", "--log-level", "error",
		"--timeout", "1",
		"--report-path", reportPath,
		"--report-format", "json",
	)

	require.NotEqual(t, 0, result.exitCode)
	assert.Less(t, time.Since(start), 4*time.Second)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "failed to scan Git repository")
}
