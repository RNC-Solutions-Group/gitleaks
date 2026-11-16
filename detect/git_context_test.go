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

func gitTestRepo(t *testing.T) (string, func(args ...string) string) {
	t.Helper()
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
	return dir, run
}

func scanGitWithOptions(t *testing.T, dir string, options sources.GitLogOptions) []report.Finding {
	t.Helper()
	d, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	d.MaxDecodeDepth = 5
	cmd, err := sources.NewGitLogCmdContextWithOptions(context.Background(), dir, "", options)
	require.NoError(t, err)
	findings, err := d.DetectSource(context.Background(), &sources.Git{Cmd: cmd, Config: &d.Config})
	require.NoError(t, err)
	return findings
}

func TestGitContextLinesRecoverMultiLineRule(t *testing.T) {
	dir, run := gitTestRepo(t)

	initial := `apiVersion: v1
kind: Secret
metadata:
  name: some-secret
  namespace: default
data: {}
`
	final := `apiVersion: v1
kind: Secret
metadata:
  name: some-secret
  namespace: default
data:
  something: c29tZXRoaW5nZ2VoZWltZXMK
  else: c29tZXRoaW5nZ2VoZWltZXMK
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "yaml.yaml"), []byte(initial), 0o600))
	run("add", "yaml.yaml")
	run("commit", "-q", "-m", "initial")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "yaml.yaml"), []byte(final), 0o600))
	run("add", "yaml.yaml")
	run("commit", "-q", "-m", "update")

	assert.Empty(t, scanGitWithOptions(t, dir, sources.GitLogOptions{}), "default addition-only behavior must remain")

	findings := scanGitWithOptions(t, dir, sources.GitLogOptions{ContextLines: 10})
	require.NotEmpty(t, findings)
	for _, finding := range findings {
		assert.Equal(t, "kubernetes-secret-yaml", finding.RuleID)
		assert.Equal(t, "yaml.yaml", finding.File)
	}
}

func TestGitContextLinesDoNotCreditUnchangedContext(t *testing.T) {
	dir, run := gitTestRepo(t)

	initial := `apiVersion: v1
kind: Secret
metadata:
  name: some-secret
  namespace: default
data:
  something: c29tZXRoaW5nZ2VoZWltZXMK
`
	// The secret is already present. The second commit adds only a line after
	// the secret block, so -U3 includes the secret as context but it must not
	// be reported a second time at the later commit.
	updated := `apiVersion: v1
kind: Secret
metadata:
  name: some-secret
  namespace: default
data:
  something: c29tZXRoaW5nZ2VoZWltZXMK
extra: unrelated-change
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "yaml.yaml"), []byte(initial), 0o600))
	run("add", "yaml.yaml")
	run("commit", "-q", "-m", "initial")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "yaml.yaml"), []byte(updated), 0o600))
	run("add", "yaml.yaml")
	run("commit", "-q", "-m", "rename namespace")

	defaultFindings := scanGitWithOptions(t, dir, sources.GitLogOptions{})
	contextFindings := scanGitWithOptions(t, dir, sources.GitLogOptions{ContextLines: 10})
	require.NotEmpty(t, defaultFindings)
	assert.Equal(t, len(defaultFindings), len(contextFindings), "context-only matches must be rejected")
}
