package sources

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gitleaks/go-gitdiff/gitdiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const colorTestSecret = "aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\""

func requireGit(t *testing.T) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available")
	}
	return gitPath
}

func gitRun(t *testing.T, gitPath, repo string, args ...string) {
	t.Helper()
	cmd := exec.Command(gitPath, append([]string{"-C", repo}, args...)...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, output)
}

func initColorRepo(t *testing.T, colorKeys ...string) (string, string) {
	t.Helper()
	gitPath := requireGit(t)
	repo := t.TempDir()
	gitRun(t, gitPath, repo, "init", "-q")
	gitRun(t, gitPath, repo, "config", "user.email", "cpv1@example.com")
	gitRun(t, gitPath, repo, "config", "user.name", "CPV1")
	gitRun(t, gitPath, repo, "config", "core.autocrlf", "false")
	gitRun(t, gitPath, repo, "config", "core.fsmonitor", "false")
	for _, key := range colorKeys {
		gitRun(t, gitPath, repo, "config", key, "always")
	}
	return repo, gitPath
}

func collectGitFiles(t *testing.T, cmd *GitCmd) []*gitdiff.File {
	t.Helper()

	var (
		wg   sync.WaitGroup
		errs []error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range cmd.ErrCh() {
			errs = append(errs, err)
		}
	}()

	var files []*gitdiff.File
	for file := range cmd.DiffFilesCh() {
		files = append(files, file)
	}
	wg.Wait()
	require.NoError(t, cmd.Wait())
	require.Empty(t, errs)
	return files
}

func filesContain(files []*gitdiff.File, name, content string) bool {
	for _, file := range files {
		if file.NewName != name {
			continue
		}
		for _, fragment := range file.TextFragments {
			if strings.Contains(fragment.Raw(gitdiff.OpAdd), content) {
				return true
			}
		}
	}
	return false
}

func TestGitDiffIgnoresForcedColor(t *testing.T) {
	tests := []struct {
		name      string
		colorKeys []string
	}{
		{name: "default color behavior"},
		{name: "color.ui always", colorKeys: []string{"color.ui"}},
		{name: "color.diff always", colorKeys: []string{"color.diff"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, gitPath := initColorRepo(t, test.colorKeys...)
			require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte(colorTestSecret+"\n"), 0o600))
			gitRun(t, gitPath, repo, "add", "secret.txt")

			cmd, err := NewGitDiffCmdContext(context.Background(), repo, true)
			require.NoError(t, err)
			files := collectGitFiles(t, cmd)
			require.True(t, filesContain(files, "secret.txt", "AKIAQ7Z3R5T7Y9U2I4O6"),
				"git color settings must not hide the staged diff: %#v", files)
		})
	}
}

func TestGitUnstagedDiffIgnoresForcedColor(t *testing.T) {
	repo, gitPath := initColorRepo(t, "color.ui", "color.diff")
	path := filepath.Join(repo, "secret.txt")
	require.NoError(t, os.WriteFile(path, []byte("clean\n"), 0o600))
	gitRun(t, gitPath, repo, "add", "secret.txt")
	gitRun(t, gitPath, repo, "commit", "-q", "-m", "base")
	require.NoError(t, os.WriteFile(path, []byte(colorTestSecret+"\n"), 0o600))

	cmd, err := NewGitDiffCmdContext(context.Background(), repo, false)
	require.NoError(t, err)
	files := collectGitFiles(t, cmd)
	require.True(t, filesContain(files, "secret.txt", "AKIAQ7Z3R5T7Y9U2I4O6"),
		"forced color must not hide the unstaged diff: %#v", files)
}

func TestGitLogIgnoresForcedColor(t *testing.T) {
	repo, gitPath := initColorRepo(t, "color.ui", "color.diff")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte(colorTestSecret+"\n"), 0o600))
	gitRun(t, gitPath, repo, "add", "secret.txt")
	gitRun(t, gitPath, repo, "commit", "-q", "-m", "add secret")

	tests := []struct {
		name    string
		logOpts string
	}{
		{name: "default log"},
		{name: "custom options", logOpts: "--all"},
		{name: "explicit color", logOpts: "--all --color"},
		{name: "explicit color always", logOpts: "--all --color=always"},
		{name: "revision plus color", logOpts: "--all HEAD --color=always"},
		{name: "pathspec separator", logOpts: "--all -- ."},
		{name: "color before pathspec separator", logOpts: "--all --color=always -- ."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, err := NewGitLogCmdContext(context.Background(), repo, test.logOpts)
			require.NoError(t, err)
			files := collectGitFiles(t, cmd)
			require.True(t, filesContain(files, "secret.txt", "AKIAQ7Z3R5T7Y9U2I4O6"),
				"forced color or explicit color options must not hide history: %#v", files)
		})
	}
}

func TestForceNoColor(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "appends after options",
			in:   []string{"log", "--color=always", "--all"},
			want: []string{"log", "--color=always", "--all", "--no-color"},
		},
		{
			name: "inserts before pathspec separator",
			in:   []string{"log", "--color", "--all", "--", "src"},
			want: []string{"log", "--color", "--all", "--no-color", "--", "src"},
		},
		{
			name: "no user color still appends",
			in:   []string{"log", "--all"},
			want: []string{"log", "--all", "--no-color"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, forceNoColor(append([]string(nil), test.in...)))
		})
	}
}
