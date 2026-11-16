package detect

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

// TestScanMergesOptIn builds the #1028 shape: a secret introduced only by an
// evil merge. The default must remain unchanged (Git omits merge patches), and
// the explicit opt-in must recover the merge-only finding. --no-diff-merges
// remains an explicit override.
func TestScanMergesOptIn(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
		return strings.TrimSpace(string(out))
	}
	run("init", "-q", "-b", "main")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@example.com")

	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.txt"), []byte("base\n"), 0o600))
	run("add", "app.txt")
	run("commit", "-q", "-m", "base")

	run("checkout", "-q", "-b", "feature")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "feature.txt"), []byte("feature\n"), 0o600))
	run("add", "feature.txt")
	run("commit", "-q", "-m", "feature")

	run("checkout", "-q", "main")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.txt"), []byte("main\n"), 0o600))
	run("add", "main.txt")
	run("commit", "-q", "-m", "main")

	// Evil merge: the branch itself is clean; the secret exists only in the
	// merge commit's tree.
	run("merge", "--no-ff", "--no-commit", "feature")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "merge-only.txt"), []byte("aws_access_key_id = AKIAZQ2M4X6P7R2T4V6Z\n"), 0o600))
	run("add", "merge-only.txt")
	run("commit", "-q", "-m", "evil merge")
	mergeSHA := run("rev-parse", "HEAD")

	scan := func(options sources.GitLogOptions, logOpts string) []report.Finding {
		t.Helper()
		d, err := NewDetectorDefaultConfig()
		require.NoError(t, err)
		d.MaxDecodeDepth = 5
		cmd, err := sources.NewGitLogCmdContextWithOptions(context.Background(), dir, logOpts, options)
		require.NoError(t, err)
		findings, err := d.DetectSource(context.Background(), &sources.Git{
			Cmd:    cmd,
			Config: &d.Config,
		})
		require.NoError(t, err)
		return findings
	}

	assert.Empty(t, scan(sources.GitLogOptions{}, ""), "default history selection must not change")

	findings := scan(sources.GitLogOptions{ScanMerges: true}, "")
	require.Len(t, findings, 1)
	assert.Equal(t, "aws-access-token", findings[0].RuleID)
	assert.Equal(t, mergeSHA, findings[0].Commit)
	assert.Equal(t, "merge-only.txt", findings[0].File)

	assert.Empty(t, scan(sources.GitLogOptions{ScanMerges: true}, "--no-diff-merges"), "user override must remain authoritative")
}
