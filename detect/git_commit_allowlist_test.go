package detect

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

	detector := NewDetector(cfg)
	gitCmd, err := sources.NewGitLogCmdContext(context.Background(), repo, "")
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
