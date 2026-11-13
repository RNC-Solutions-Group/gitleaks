package sources

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
)

const gitTestDiff = `commit 6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b
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

const gitTestMalformedDiff = gitTestDiff + `diff --git a/bad.txt b/bad.txt
index 2222222..3333333 100644
--- a/bad.txt
+++ b/bad.txt
@@ -231,18446744073709551615 +231,1 @@
+bad
diff --git a/second.txt b/second.txt
new file mode 100644
index 0000000..4444444
--- /dev/null
+++ b/second.txt
@@ -0,0 +1 @@
+second
`

// TestGitOutputHelper is a subprocess used to produce controlled Git stdout,
// stderr, exit status, and lifetime without invoking a real Git executable.
func TestGitOutputHelper(t *testing.T) {
	if os.Getenv("GITLEAKS_GIT_OUTPUT_HELPER") != "1" {
		return
	}
	if path := os.Getenv("GITLEAKS_GIT_STDOUT"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			os.Exit(91)
		}
		if _, err := os.Stdout.Write(data); err != nil {
			os.Exit(92)
		}
	}
	if path := os.Getenv("GITLEAKS_GIT_STDERR"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			os.Exit(93)
		}
		if _, err := os.Stderr.Write(data); err != nil {
			os.Exit(94)
		}
	}
	if value := os.Getenv("GITLEAKS_GIT_SLEEP_MS"); value != "" {
		milliseconds, err := strconv.Atoi(value)
		if err != nil {
			os.Exit(95)
		}
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	}
	exitCode, err := strconv.Atoi(os.Getenv("GITLEAKS_GIT_EXIT_CODE"))
	if err != nil {
		os.Exit(96)
	}
	os.Exit(exitCode)
}

func newControlledGitCmd(
	t *testing.T,
	ctx context.Context,
	stdout, stderr string,
	exitCode int,
	sleep time.Duration,
) *GitCmd {
	t.Helper()
	dir := t.TempDir()
	stdoutPath := filepath.Join(dir, "stdout")
	stderrPath := filepath.Join(dir, "stderr")
	require.NoError(t, os.WriteFile(stdoutPath, []byte(stdout), 0o600))
	require.NoError(t, os.WriteFile(stderrPath, []byte(stderr), 0o600))

	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGitOutputHelper$")
	cmd.Env = append(os.Environ(),
		"GITLEAKS_GIT_OUTPUT_HELPER=1",
		"GITLEAKS_GIT_STDOUT="+stdoutPath,
		"GITLEAKS_GIT_STDERR="+stderrPath,
		"GITLEAKS_GIT_EXIT_CODE="+strconv.Itoa(exitCode),
		"GITLEAKS_GIT_SLEEP_MS="+strconv.Itoa(int(sleep/time.Millisecond)),
	)
	gitCmd, err := newGitCmd(cmd, dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = gitCmd.Wait()
	})
	return gitCmd
}

func collectGitFragments(t *testing.T, ctx context.Context, gitCmd *GitCmd) ([]Fragment, error) {
	t.Helper()
	source := &Git{
		Cmd:    gitCmd,
		Config: &config.Config{},
		Sema:   semgroup.NewGroup(ctx, 4),
	}
	var fragments []Fragment
	err := source.Fragments(ctx, func(fragment Fragment, err error) error {
		if err != nil {
			return err
		}
		fragments = append(fragments, fragment)
		return nil
	})
	return fragments, err
}

func TestGitFragmentsReportsNonzeroExitWithoutStderr(t *testing.T) {
	first := newControlledGitCmd(t, context.Background(), "", "", 128, 0)
	_, err := collectGitFragments(t, context.Background(), first)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "git command failed")
	assert.Contains(t, err.Error(), "exit status 128")
}

func TestGitFragmentsReportsNonzeroExitWithStderr(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), "", "fatal: synthetic git failure\n", 128, 0)
	_, err := collectGitFragments(t, context.Background(), gitCmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fatal: synthetic git failure")
}

func TestGitFragmentsAllowsNonfatalStderr(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), gitTestDiff, "warning: synthetic nonfatal warning\n", 0, 0)
	fragments, err := collectGitFragments(t, context.Background(), gitCmd)
	require.NoError(t, err)
	require.Len(t, fragments, 1)
	assert.Equal(t, "first.txt", fragments[0].FilePath)
}

func TestGitFragmentsReportsFatalStderrWithZeroExit(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), "", "fatal: synthetic zero-exit failure\n", 0, 0)
	_, err := collectGitFragments(t, context.Background(), gitCmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "git reported an error on stderr")
	assert.Contains(t, err.Error(), "synthetic zero-exit failure")
}

func TestGitFragmentsReportsIncompleteDiff(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), gitTestMalformedDiff, "", 0, 0)
	finished := make(chan error, 1)
	var fragments []Fragment
	go func() {
		var err error
		fragments, err = collectGitFragments(t, context.Background(), gitCmd)
		finished <- err
	}()

	select {
	case err := <-finished:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "diff parsing stopped early")
		assert.Contains(t, err.Error(), "results are incomplete")
		require.Len(t, fragments, 1)
		assert.Equal(t, "first.txt", fragments[0].FilePath)
	case <-time.After(10 * time.Second):
		t.Fatal("Git.Fragments did not return after the parser stopped early")
	}
}

func TestGitFragmentsCancellationWaitsForProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	gitCmd := newControlledGitCmd(t, ctx, gitTestDiff, "", 0, 10*time.Second)
	start := time.Now()
	_, err := collectGitFragments(t, ctx, gitCmd)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, elapsed, 5*time.Second, "cancellation must not wait for the helper's full sleep")
	require.NotNil(t, gitCmd.cmd.ProcessState, "the Git process must be reaped before Fragments returns")
}

func TestGitCmdWaitIsIdempotent(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), gitTestDiff, "", 0, 0)
	for range gitCmd.DiffFilesCh() {
	}
	for range gitCmd.ErrCh() {
	}
	require.NoError(t, gitCmd.Wait())
	require.NoError(t, gitCmd.Wait())
}

const gitTestTruncatedDiff = `commit 6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b
Author: Test <test@example.com>
Date:   Sun Feb 11 21:44:04 2024 +0000

    a commit

diff --git a/truncated.txt b/truncated.txt
new file mode 100644
index 0000000..1111111
--- /dev/null
+++ b/truncated.txt
@@ -0,0 +1,2 @@
+only one line of a two-line hunk
`

func TestGitFragmentsReportsTruncatedDiff(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), gitTestTruncatedDiff, "", 0, 0)
	_, err := collectGitFragments(t, context.Background(), gitCmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "diff parsing stopped early")
	assert.Contains(t, err.Error(), "results are incomplete")
}

func TestGitCmdStartFailure(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := NewGitLogCmdContext(context.Background(), ".", "")
	require.Error(t, err)
}

func TestGitOutputTrackerCountsDiffSections(t *testing.T) {
	tracker := &gitOutputTracker{}
	_, _ = tracker.Write([]byte("    diff --git a/message b/message\n"))
	_, _ = tracker.Write([]byte("diff --git a/one b/one\n--- a/one\n+++ b/one\n"))
	_, _ = tracker.Write([]byte("diff --git a/two b/two"))
	tracker.finish()
	require.Equal(t, int64(2), tracker.sectionCount())
}

func TestSanitizeGitDiagnostic(t *testing.T) {
	tests := map[string]string{
		"fatal: synthetic git failure": "fatal: synthetic git failure",
		"fatal: unable to access 'https://user:ghp_SECRETTOKEN@example.com/repo?token=QUERYVALUE': denied": "fatal: unable to access 'https://REDACTED@example.com/repo?token=REDACTED': denied",
		"fatal: bad token AKIAABCDEFGHIJKLMNOP": "fatal: bad token REDACTED",
	}
	for input, expected := range tests {
		assert.Equal(t, expected, sanitizeGitDiagnostic(input))
	}
}

func TestGitFragmentsSanitizesStderrSecrets(t *testing.T) {
	gitCmd := newControlledGitCmd(t, context.Background(), "",
		"fatal: unable to access 'https://user:ghp_SECRETTOKEN@example.com/repo?token=QUERYVALUE': denied\n",
		128, 0)
	_, err := collectGitFragments(t, context.Background(), gitCmd)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "ghp_SECRETTOKEN")
	assert.NotContains(t, err.Error(), "QUERYVALUE")
	assert.Contains(t, err.Error(), "https://REDACTED@example.com")
	assert.Contains(t, err.Error(), "token=REDACTED")
}

func TestIsCommitAllowed(t *testing.T) {
	tests := map[string]struct {
		allowlists  []*config.Allowlist
		sha         string
		wantAllowed bool
		wantCommit  string
	}{
		"no allowlists": {
			allowlists:  nil,
			sha:         "commitA",
			wantAllowed: false,
		},
		"single allowlist matches": {
			allowlists:  []*config.Allowlist{{Commits: []string{"commitA"}}},
			sha:         "commitA",
			wantAllowed: true,
			wantCommit:  "commitA",
		},
		"single allowlist does not match": {
			allowlists:  []*config.Allowlist{{Commits: []string{"commitB"}}},
			sha:         "commitA",
			wantAllowed: false,
		},
		"empty sha is never allowed": {
			allowlists:  []*config.Allowlist{{Commits: []string{"commitA"}}},
			sha:         "",
			wantAllowed: false,
		},
		// The original bug: the allowed commit lived in a later allowlist and the
		// misscoped continue let it fall through to be scanned anyway. Match on a
		// non-first entry guards against that regression.
		"match in a later allowlist": {
			allowlists: []*config.Allowlist{
				{Commits: []string{"commitB"}},
				{Commits: []string{"commitA"}},
			},
			sha:         "commitA",
			wantAllowed: true,
			wantCommit:  "commitA",
		},
		// Overlapping allowlists: both cover the commit. We stop at the first
		// match and never scan, regardless of how many entries match.
		"overlapping allowlists both match": {
			allowlists: []*config.Allowlist{
				{Commits: []string{"commitA"}},
				{Commits: []string{"commitA"}},
			},
			sha:         "commitA",
			wantAllowed: true,
			wantCommit:  "commitA",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotAllowed, gotCommit := isCommitAllowed(tt.allowlists, tt.sha)
			assert.Equal(t, tt.wantAllowed, gotAllowed)
			if tt.wantAllowed {
				assert.Equal(t, tt.wantCommit, gotCommit)
			}
		})
	}
}
