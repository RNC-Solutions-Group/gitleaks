//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLITimeoutIsOperationalFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess timeout test in short mode")
	}

	dir := t.TempDir()
	fakeGitDir := filepath.Join(dir, "bin")
	require.NoError(t, os.Mkdir(fakeGitDir, 0o700))

	diffPath := filepath.Join(dir, "diff.patch")
	require.NoError(t, os.WriteFile(diffPath, []byte(gitTestDiffForTimeout), 0o600))
	script := "#!/bin/sh\ncat \"$FAKE_GIT_STDOUT\"\nexec sleep 30\n"
	require.NoError(t, os.WriteFile(filepath.Join(fakeGitDir, "git"), []byte(script), 0o700))

	source := t.TempDir()
	reportPath := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(reportPath, []byte("SENTINEL"), 0o600))

	env := map[string]string{
		"PATH":            fakeGitDir + string(os.PathListSeparator) + os.Getenv("PATH"),
		"FAKE_GIT_STDOUT": diffPath,
	}
	result := runReportCLIWithEnv(t, env,
		"git", source,
		"--no-banner",
		"--log-level", "error",
		"--platform", "none",
		"--timeout", "1",
		"--exit-code", "0",
		"--report-path", reportPath,
		"--report-format", "json",
	)

	require.NotEqual(t, 0, result.exitCode)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, reportPath)))
	assert.True(t, strings.Contains(result.stderr, "context deadline exceeded") ||
		strings.Contains(result.stderr, "failed to scan Git repository"), result.stderr)
}

const gitTestDiffForTimeout = `commit 6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b
Author: Test <test@example.com>
Date:   Sun Feb 11 21:44:04 2024 +0000

    a commit

diff --git a/first.txt b/first.txt
new file mode 100644
index 0000000..1111111
--- /dev/null
+++ b/first.txt
@@ -0,0 +1 @@
+first
`
