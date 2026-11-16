package sources

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitLogOptionsContextBounds(t *testing.T) {
	assert.NoError(t, GitLogOptions{}.validate())
	assert.NoError(t, GitLogOptions{ContextLines: maxGitContextLines}.validate())
	assert.Error(t, GitLogOptions{ContextLines: -1}.validate())
	assert.Error(t, GitLogOptions{ContextLines: maxGitContextLines + 1}.validate())
}

func TestQuotedNoNotesReachesGitCommand(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
	}
	run("init", "-q")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@example.com")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content\n"), 0o600))
	run("add", "file.txt")
	run("commit", "-q", "-m", "initial")

	cmd, err := NewGitLogCmdContext(context.Background(), dir, "'--no-notes'")
	require.NoError(t, err)
	assert.Contains(t, cmd.String(), "--no-notes")
	assert.NotContains(t, cmd.String(), "'--no-notes'")
	for range cmd.DiffFilesCh() {
	}
	require.NoError(t, cmd.Wait())
}
