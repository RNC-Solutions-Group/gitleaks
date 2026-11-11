package sources

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gitleaks/go-gitdiff/gitdiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
)

func newIndexTestRepo(t *testing.T) string {
	t.Helper()
	gitPath := requireGit(t)
	repo := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(gitPath, append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)...)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v: %s", args, output)
	}
	run("init", "-q")
	run("config", "core.fsmonitor", "false")
	run("config", "user.email", "cpv4@example.com")
	run("config", "user.name", "CPV4")
	return repo
}

func TestStagedIndexBlobFailureIsReturned(t *testing.T) {
	repo := newIndexTestRepo(t)
	source := &Git{
		Cmd:    &GitCmd{repoPath: repo, staged: true},
		Config: &config.Config{},
	}
	err := source.processDiffFile(
		context.Background(),
		&gitdiff.File{IsBinary: true, NewName: "missing.txt"},
		true,
		"",
		nil,
		func(Fragment, error) error { return nil },
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not read blob")
}

func TestStagedIndexBlobCancellationIsReturned(t *testing.T) {
	repo := newIndexTestRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	source := &Git{
		Cmd:    &GitCmd{repoPath: repo, staged: true},
		Config: &config.Config{},
	}
	err := source.processDiffFile(
		ctx,
		&gitdiff.File{IsBinary: true, NewName: "missing.txt"},
		true,
		"",
		nil,
		func(Fragment, error) error { return nil },
	)
	require.Error(t, err)
}

func TestIndexBlobReaderHandlesSpecialPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not permit ':' in a file name")
	}
	repo := newIndexTestRepo(t)
	name := "colon:name.txt"
	require.NoError(t, os.WriteFile(filepath.Join(repo, name), []byte("synthetic staged content\n"), 0o600))

	gitPath := requireGit(t)
	cmd := exec.Command(gitPath, "-c", "core.fsmonitor=false", "-C", repo, "add", "--", name)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	gitCmd := &GitCmd{repoPath: repo, staged: true}
	reader, err := gitCmd.NewIndexBlobReaderContext(context.Background(), name)
	require.NoError(t, err)
	data, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	require.NoError(t, readErr)
	require.NoError(t, closeErr)
	assert.Equal(t, "synthetic staged content\n", string(data))
}

type failingReadCloser struct{ err error }

func (r *failingReadCloser) Read([]byte) (int, error) { return 0, r.err }
func (r *failingReadCloser) Close() error             { return nil }

// TestBlobReaderHelperProcess is a no-op child used to give blobReader.Close a
// real command to wait for while the read side fails independently.
func TestBlobReaderHelperProcess(t *testing.T) {}

func TestBlobReaderCloseReturnsReadError(t *testing.T) {
	child := exec.Command(os.Args[0], "-test.run=^TestBlobReaderHelperProcess$")
	child.Stdout = io.Discard
	child.Stderr = io.Discard
	require.NoError(t, child.Start())

	readErr := errors.New("synthetic blob read failure")
	reader := &blobReader{
		ReadCloser: &failingReadCloser{err: readErr},
		cmd:        child,
	}
	err := reader.Close()
	require.Error(t, err)
	assert.ErrorIs(t, err, readErr)
}

func TestBlobReaderCloseAfterSkipStopsProducer(t *testing.T) {
	child := exec.Command(os.Args[0], "-test.run=^TestGitOutputHelper$")
	child.Env = append(os.Environ(),
		"GITLEAKS_GIT_OUTPUT_HELPER=1",
		"GITLEAKS_GIT_SLEEP_MS=10000",
		"GITLEAKS_GIT_EXIT_CODE=0",
	)
	stdout, err := child.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, child.Start())

	reader := &blobReader{ReadCloser: stdout, cmd: child}
	start := time.Now()
	require.NoError(t, reader.CloseAfterSkip())
	assert.Less(t, time.Since(start), 5*time.Second)
}
