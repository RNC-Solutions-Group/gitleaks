package sources

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func newFilesTestSource(ctx context.Context, dir string) *Files {
	return &Files{
		Config: &config.Config{},
		Path:   dir,
		Sema:   semgroup.NewGroup(ctx, 4),
	}
}

func TestFilesReportsUnreadableFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not enforceable this way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "readable.txt"), []byte("clean\n"), 0o600))
	locked := filepath.Join(dir, "locked.txt")
	require.NoError(t, os.WriteFile(locked, []byte("synthetic secret\n"), 0o600))
	require.NoError(t, os.Chmod(locked, 0o000))
	t.Cleanup(func() { _ = os.Chmod(locked, 0o600) })

	err := newFilesTestSource(context.Background(), dir).Fragments(
		context.Background(),
		func(Fragment, error) error { return nil },
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "locked.txt")
	assert.Contains(t, err.Error(), "could not open file")
}

func TestFilesReportsUnreadableDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not enforceable this way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	blocked := filepath.Join(dir, "blocked")
	require.NoError(t, os.Mkdir(blocked, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(blocked, "secret.txt"), []byte("synthetic secret\n"), 0o600))
	require.NoError(t, os.Chmod(blocked, 0o000))
	t.Cleanup(func() { _ = os.Chmod(blocked, 0o700) })

	err := newFilesTestSource(context.Background(), dir).Fragments(
		context.Background(),
		func(Fragment, error) error { return nil },
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blocked")
}

func TestFilesBoundsManyWorkerErrors(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not enforceable this way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	for i := 0; i < 14; i++ {
		path := filepath.Join(dir, "locked-"+strings.Repeat("x", i)+".txt")
		require.NoError(t, os.WriteFile(path, []byte("content\n"), 0o600))
		require.NoError(t, os.Chmod(path, 0o000))
		pathCopy := path
		t.Cleanup(func() { _ = os.Chmod(pathCopy, 0o600) })
	}

	err := newFilesTestSource(context.Background(), dir).Fragments(
		context.Background(),
		func(Fragment, error) error { return nil },
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "additional scan error(s) omitted")
	assert.Contains(t, err.Error(), "could not open file")
}

func TestFilesPreservesAllowlistedUnreadableFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not enforceable this way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	locked := filepath.Join(dir, "allowed.txt")
	require.NoError(t, os.WriteFile(locked, []byte("content\n"), 0o600))
	require.NoError(t, os.Chmod(locked, 0o000))
	t.Cleanup(func() { _ = os.Chmod(locked, 0o600) })

	allowlist := &config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`allowed\.txt$`)},
	}
	require.NoError(t, allowlist.Validate())
	source := &Files{
		Config: &config.Config{Allowlists: []*config.Allowlist{allowlist}},
		Path:   dir,
		Sema:   semgroup.NewGroup(context.Background(), 4),
	}
	require.NoError(t, source.Fragments(context.Background(), func(Fragment, error) error { return nil }))
}

func TestFilesWaitsForWorkersAfterCancellation(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 16; i++ {
		require.NoError(t, os.WriteFile(
			filepath.Join(dir, "file-"+strings.Repeat("x", i)+".txt"),
			[]byte("content\n"),
			0o600,
		))
	}

	ctx, cancel := context.WithCancel(context.Background())
	var active atomic.Int64
	source := newFilesTestSource(ctx, dir)
	err := source.Fragments(ctx, func(Fragment, error) error {
		active.Add(1)
		defer active.Add(-1)
		cancel()
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, active.Load(), "workers must be joined before Fragments returns")
}

type partialErrorReader struct {
	data []byte
	sent bool
	err  error
}

func (r *partialErrorReader) Read(p []byte) (int, error) {
	if !r.sent {
		r.sent = true
		n := copy(p, r.data)
		return n, r.err
	}
	return 0, r.err
}

func TestFileReportsReadFailureAfterPartialContent(t *testing.T) {
	readErr := errors.New("synthetic read failure")
	file := &File{
		Content: &partialErrorReader{
			data: []byte("synthetic content\n"),
			err:  readErr,
		},
		Path: "input.txt",
	}

	var fragments []Fragment
	err := file.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, readErr)
	require.Len(t, fragments, 1, "content read before the failure must still be scanned")
	assert.Equal(t, "synthetic content\n", fragments[0].Raw)
}

func TestFileReportsTruncatedArchive(t *testing.T) {
	file := &File{
		Content:         bytes.NewReader([]byte("PK\x03\x04truncated")),
		Path:            "broken.zip",
		MaxArchiveDepth: 1,
	}
	err := file.Fragments(context.Background(), func(Fragment, error) error { return nil })
	require.Error(t, err)
}
