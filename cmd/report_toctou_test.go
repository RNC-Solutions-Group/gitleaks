package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests exercise filesystem changes between report validation and the
// final rename. They are adversarial fixtures, not promises that two
// filesystem operations are atomic.

func TestReportOverwriteReplacesHardlinkWithoutTouchingTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(target, []byte("PROTECTED-TARGET"), 0o600))
	require.NoError(t, os.Link(target, dest))

	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	_, err = out.write(&stringReporter{payload: "NEW-REPORT"}, nil)
	require.NoError(t, err)

	assert.Equal(t, "PROTECTED-TARGET", string(mustReadFile(t, target)),
		"replacing a hardlinked report destination must not modify the other link")
	assert.Equal(t, "NEW-REPORT", string(mustReadFile(t, dest)))
	requireNoTemporaryReports(t, dir)
}

func TestReportInputHardlinkCollisionRejected(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(configPath, []byte("title = \"protected\"\n"), 0o600))
	require.NoError(t, os.Link(configPath, dest))

	_, err := newReportOutput(dest, reportOverwrite, "", []string{configPath}, "", "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "collides with scan input")
	assert.Equal(t, "title = \"protected\"\n", string(mustReadFile(t, configPath)))
}

func TestReportOverwriteDirectorySubstitutionDoesNotPublish(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-REPORT"), 0o600))

	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	out.rename = func(oldPath, newPath string) error {
		require.NoError(t, os.Remove(newPath))
		require.NoError(t, os.Mkdir(newPath, 0o700))
		return os.Rename(oldPath, newPath)
	}

	_, err = out.write(&stringReporter{payload: "NEW-REPORT"}, nil)
	require.Error(t, err)
	info, statErr := os.Stat(dest)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir(), "the substituted directory must remain a directory")
	requireNoTemporaryReports(t, dir)
}

func TestReportOverwriteSymlinkSubstitutionDoesNotFollowTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(target, []byte("PROTECTED-TARGET"), 0o600))
	require.NoError(t, os.WriteFile(dest, []byte("OLD-REPORT"), 0o600))

	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	out.rename = func(oldPath, newPath string) error {
		if err := os.Remove(newPath); err != nil {
			return err
		}
		if err := os.Symlink(target, newPath); err != nil {
			return err
		}
		return os.Rename(oldPath, newPath)
	}

	_, err = out.write(&stringReporter{payload: "NEW-REPORT"}, nil)
	if runtime.GOOS == "windows" && err != nil {
		t.Skipf("Windows could not replace a symlink destination: %v", err)
	}
	require.NoError(t, err)
	assert.Equal(t, "PROTECTED-TARGET", string(mustReadFile(t, target)),
		"publication must replace the link, not the link target")
	assert.Equal(t, "NEW-REPORT", string(mustReadFile(t, dest)))
	info, statErr := os.Lstat(dest)
	require.NoError(t, statErr)
	assert.Zero(t, info.Mode()&os.ModeSymlink, "destination should be the published regular file")
	requireNoTemporaryReports(t, dir)
}

func TestReportOverwriteParentSubstitutionDoesNotPublishToNewParent(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "reports")
	parked := filepath.Join(base, "parked-reports")
	require.NoError(t, os.Mkdir(dir, 0o700))
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-REPORT"), 0o600))

	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	out.rename = func(oldPath, newPath string) error {
		if err := os.Rename(dir, parked); err != nil {
			return err
		}
		if err := os.Mkdir(dir, 0o700); err != nil {
			return err
		}
		return os.Rename(oldPath, newPath)
	}

	_, err = out.write(&stringReporter{payload: "NEW-REPORT"}, nil)
	require.Error(t, err)
	assert.NoFileExists(t, filepath.Join(dir, "report.json"),
		"a substituted parent must not receive the report")
	// The temporary file created in the original parent can outlive a failed
	// publication when that parent is renamed. This is a documented cleanup
	// limitation; t.TempDir removes the parked directory after the test.
	leftovers, globErr := filepath.Glob(filepath.Join(parked, ".*.tmp-*"))
	require.NoError(t, globErr)
	t.Logf("temporary files left in renamed parent: %v", leftovers)
}

func TestReportOverwriteConcurrentWritersPublishCompleteFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)

	const writers = 12
	payloads := make(chan string, writers)
	errs := make(chan error, writers)
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		payload := fmt.Sprintf("REPORT-%02d", i)
		payloads <- payload
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, writeErr := out.write(&stringReporter{payload: payload}, nil)
			errs <- writeErr
		}()
	}
	wg.Wait()
	close(errs)
	close(payloads)
	for writeErr := range errs {
		require.NoError(t, writeErr)
	}

	final := string(mustReadFile(t, dest))
	expected := map[string]struct{}{}
	for payload := range payloads {
		expected[payload] = struct{}{}
	}
	_, complete := expected[final]
	assert.True(t, complete, "final report must be one complete writer's payload, got %q", final)
	requireNoTemporaryReports(t, dir)
}
