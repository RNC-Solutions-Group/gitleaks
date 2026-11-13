//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIStagedBlobFailureIsOperational(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	dir := t.TempDir()
	fakeGitDir := filepath.Join(dir, "bin")
	require.NoError(t, os.Mkdir(fakeGitDir, 0o700))

	diffPath := filepath.Join(dir, "diff.patch")
	require.NoError(t, os.WriteFile(diffPath, []byte(stagedBlobFailureDiff), 0o600))
	script := "#!/bin/sh\ncase \"$*\" in\n  *\" cat-file \"*) exit 128 ;;\n  *\" diff \"*) cat \"$FAKE_GIT_DIFF\"; exit 0 ;;\n  *) exit 128 ;;\nesac\n"
	require.NoError(t, os.WriteFile(filepath.Join(fakeGitDir, "git"), []byte(script), 0o700))

	source := t.TempDir()
	reportPath := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	result := runReportCLIWithEnv(t, map[string]string{
		"PATH":          fakeGitDir + string(os.PathListSeparator) + os.Getenv("PATH"),
		"FAKE_GIT_DIFF": diffPath,
	},
		"git", source, "--staged", "--platform", "none",
		"--no-banner", "--log-level", "error",
		"--exit-code", "0",
		"--report-path", reportPath,
		"--report-format", "json",
	)

	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.Contains(t, result.stderr, "failed to scan Git repository")
	assert.Contains(t, result.stderr, "could not read blob")
}

const stagedBlobFailureDiff = `diff --git a/secret.txt b/secret.txt
new file mode 100644
index 0000000..fd5083c
Binary files /dev/null and b/secret.txt differ
`
