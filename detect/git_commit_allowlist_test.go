package detect

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func scanGitHistory(t *testing.T, repo string, cfg config.Config) []report.Finding {
	t.Helper()
	return scanGitHistoryWithOpts(t, repo, cfg, "")
}

func scanGitHistoryWithOpts(t *testing.T, repo string, cfg config.Config, logOpts string) []report.Finding {
	t.Helper()

	detector := NewDetector(cfg)
	gitCmd, err := sources.NewGitLogCmdContext(context.Background(), repo, logOpts)
	require.NoError(t, err)

	findings, err := detector.DetectSource(context.Background(), &sources.Git{
		Cmd:             gitCmd,
		Config:          &detector.Config,
		Remote:          sources.NewRemoteInfo(scm.NoPlatform, repo),
		Sema:            detector.Sema,
		MaxArchiveDepth: detector.MaxArchiveDepth,
	})
	require.NoError(t, err)
	return findings
}

func runGitExpectConflict(t *testing.T, repo string, args ...string) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	require.NoError(t, err)
	cmd := exec.Command(gitPath, append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)...)
	output, err := cmd.CombinedOutput()
	require.Error(t, err, "expected git %v to fail with a merge conflict:\n%s", args, output)
	return string(output)
}

func commitTestFile(t *testing.T, repo, name, content, message string) string {
	t.Helper()
	path := filepath.Join(repo, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	runGitInRepo(t, repo, "add", "--", name)
	runGitInRepo(t, repo, "commit", "-q", "-m", message)
	return strings.TrimSpace(runGitInRepo(t, repo, "rev-parse", "HEAD"))
}

func commitAllowlistConfig(allowlistedCommit string) config.Config {
	rule := config.Rule{
		RuleID:      "cpv11-commit-test",
		Description: "synthetic CP-V11 commit allowlist credential",
		Regex:       regexp.MustCompile(`CPV11_(?:ALLOWED|OTHER)_[A-Z0-9]{4}`),
		Keywords:    []string{"cpv11_"},
	}
	cfg := config.Config{
		Rules:        map[string]config.Rule{rule.RuleID: rule},
		Keywords:     map[string]struct{}{"cpv11_": {}},
		OrderedRules: []string{rule.RuleID},
	}
	if allowlistedCommit != "" {
		cfg.Allowlists = []*config.Allowlist{{Commits: []string{allowlistedCommit}}}
	}
	return cfg
}

// TestGitCommitAllowlistSkipsEntireCommit exercises the source-level callback
// path, not just the pure allowlist helper. The pre-fix code logged the match
// but continued scanning because the continue applied to the inner allowlist
// loop. Both commits contain synthetic credentials; only the allowed commit
// must be suppressed.
func TestGitLogNeutralizesRepositoryFormatting(t *testing.T) {
	repo := initStagedTestRepo(t)
	runGitInRepo(t, repo, "config", "format.pretty", "%H")
	runGitInRepo(t, repo, "config", "log.date", "relative")
	sha := commitTestFile(t, repo, "secret.txt", "cpv11_token = \"CPV11_OTHER_0003\"\n", "add formatted-history credential")

	findings := scanGitHistory(t, repo, commitAllowlistConfig(""))
	require.Len(t, findings, 1)
	assert.Equal(t, sha, findings[0].Commit)

	authorDate := strings.TrimSpace(runGitInRepo(t, repo, "show", "-s", "--format=%aI", sha))
	parsed, err := time.Parse(time.RFC3339, authorDate)
	require.NoError(t, err)
	assert.Equal(t, parsed.UTC().Format(time.RFC3339), findings[0].Date)
}

func TestGitMergeOnlyFindingRequiresExplicitMergeDiff(t *testing.T) {
	repo := initStagedTestRepo(t)
	baseBranch := strings.TrimSpace(runGitInRepo(t, repo, "rev-parse", "--abbrev-ref", "HEAD"))

	runGitInRepo(t, repo, "checkout", "-q", "-b", "feature")
	commitTestFile(t, repo, "conflict.txt", "value=feature\n", "feature side")
	runGitInRepo(t, repo, "checkout", "-q", baseBranch)
	commitTestFile(t, repo, "conflict.txt", "value=main\n", "main side")

	// The merge has a conflict; resolve it cleanly and introduce a credential
	// only in the merge commit itself.
	runGitExpectConflict(t, repo, "merge", "--no-ff", "feature")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "conflict.txt"), []byte("resolved\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, "merge-only.txt"), []byte("cpv11_token = \"CPV11_OTHER_0005\"\n"), 0o600))
	runGitInRepo(t, repo, "add", "--", "conflict.txt", "merge-only.txt")
	runGitInRepo(t, repo, "commit", "-q", "-m", "merge with merge-only credential")
	mergeSHA := strings.TrimSpace(runGitInRepo(t, repo, "rev-parse", "HEAD"))

	cfg := commitAllowlistConfig("")
	defaultFindings := scanGitHistory(t, repo, cfg)
	for _, finding := range defaultFindings {
		assert.False(t, finding.Commit == mergeSHA && finding.Secret == "CPV11_OTHER_0005",
			"default git log output omits merge-only changes")
	}

	optInFindings := scanGitHistoryWithOpts(t, repo, cfg, "--diff-merges=first-parent")
	found := false
	for _, finding := range optInFindings {
		if finding.Commit == mergeSHA && finding.Secret == "CPV11_OTHER_0005" {
			found = true
			break
		}
	}
	assert.True(t, found, "explicit --diff-merges=first-parent must expose merge-only changes")
}

func TestGitLogDoesNotInvokeRepositoryTextconv(t *testing.T) {
	repo := initStagedTestRepo(t)
	runGitInRepo(t, repo, "config", "diff.evil.textconv", "gitleaks-v12-missing-textconv")
	require.NoError(t, os.WriteFile(filepath.Join(repo, ".gitattributes"), []byte("*.bin diff=evil\n"), 0o600))
	runGitInRepo(t, repo, "add", "--", ".gitattributes")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.bin"), []byte("\x00\x01CPV11_OTHER_0004\n"), 0o600))
	runGitInRepo(t, repo, "add", "--", "secret.bin")
	runGitInRepo(t, repo, "commit", "-q", "-m", "binary with repository textconv")

	findings := scanGitHistory(t, repo, commitAllowlistConfig(""))
	assert.Empty(t, findings)
}

func TestGitCommitAllowlistSkipsEntireCommit(t *testing.T) {
	repo := initStagedTestRepo(t)
	allowedSHA := commitTestFile(t, repo, "allowed.txt", "cpv11_token = \"CPV11_ALLOWED_0001\"\n", "add allowed credential")
	otherSHA := commitTestFile(t, repo, "other.txt", "cpv11_token = \"CPV11_OTHER_0002\"\n", "add other credential")

	control := scanGitHistory(t, repo, commitAllowlistConfig(""))
	require.Len(t, control, 2, "control scan must see both synthetic credentials")

	findings := scanGitHistory(t, repo, commitAllowlistConfig(allowedSHA))
	require.Len(t, findings, 1, "the allowlisted commit must not produce findings")
	assert.Equal(t, otherSHA, findings[0].Commit)
	assert.Equal(t, "other.txt", findings[0].File)
	assert.Equal(t, "CPV11_OTHER_0002", findings[0].Secret)
}
