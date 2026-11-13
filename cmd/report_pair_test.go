package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/report"
)

func TestInformativeCompanionPath(t *testing.T) {
	tests := []struct {
		primary string
		want    string
	}{
		{primary: filepath.Join("reports", "report.json"), want: filepath.Join("reports", "report-informative.json")},
		{primary: filepath.Join("reports", "report.csv"), want: filepath.Join("reports", "report-informative.json")},
		{primary: filepath.Join("reports", "report.sarif"), want: filepath.Join("reports", "report-informative.json")},
		{primary: filepath.Join("reports", "report_001.json"), want: filepath.Join("reports", "report_001-informative.json")},
		{primary: filepath.Join("reports", "report_002.csv"), want: filepath.Join("reports", "report_002-informative.json")},
		{primary: filepath.Join("reports", "report"), want: filepath.Join("reports", "report-informative.json")},
		{primary: filepath.Join("reports", ".report"), want: filepath.Join("reports", ".report-informative.json")},
		{primary: filepath.Join("reports", "archive.tar.gz"), want: filepath.Join("reports", "archive.tar-informative.json")},
		{primary: filepath.Join("unicodé", "报告.json"), want: filepath.Join("unicodé", "报告-informative.json")},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, informativeCompanionPath(test.primary), test.primary)
	}
}

func newInformativeTestOutput(t *testing.T, path string, mode reportMode, render func(string) ([]byte, error)) *reportOutput {
	t.Helper()
	out, err := newReportOutput(path, mode, "", nil, "", ".", "")
	require.NoError(t, err)
	if render == nil {
		render = func(primary string) ([]byte, error) {
			return []byte(fmt.Sprintf("{\"primary\":%q}\n", primary)), nil
		}
	}
	require.NoError(t, out.enableInformative(render))
	return out
}

func TestInformativePairOverwrite(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	primaryPath, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.NoError(t, err)
	assert.Equal(t, dest, primaryPath)
	assert.Equal(t, "NEW-PRIMARY", string(mustReadFile(t, dest)))
	// Overwrite is the explicit replacement policy for both members of the
	// logical pair.
	assertCompanionPrimary(t, companion, dest)
	if runtime.GOOS != "windows" {
		info, statErr := os.Stat(companion)
		require.NoError(t, statErr)
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairAggregateAllocation(t *testing.T) {
	tests := []struct {
		name              string
		existingPrimary   bool
		existingCompanion bool
		wantPrimary       string
		wantCompanion     string
	}{
		{
			name:          "both absent",
			wantPrimary:   "report.json",
			wantCompanion: "report-informative.json",
		},
		{
			name:            "existing primary only",
			existingPrimary: true,
			wantPrimary:     "report_001.json",
			wantCompanion:   "report_001-informative.json",
		},
		{
			name:              "existing companion only",
			existingCompanion: true,
			wantPrimary:       "report_001.json",
			wantCompanion:     "report_001-informative.json",
		},
		{
			name:              "both existing",
			existingPrimary:   true,
			existingCompanion: true,
			wantPrimary:       "report_001.json",
			wantCompanion:     "report_001-informative.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			dest := filepath.Join(dir, "report.json")
			companion := filepath.Join(dir, "report-informative.json")
			if test.existingPrimary {
				require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
			}
			if test.existingCompanion {
				require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))
			}

			out, err := newReportOutput(dest, reportAggregate, "", nil, "", ".", "")
			require.NoError(t, err)
			require.NoError(t, out.enableInformative(func(primary string) ([]byte, error) {
				return []byte(fmt.Sprintf("{\"primary\":%q}\n", primary)), nil
			}))
			primaryPath, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(dir, test.wantPrimary), primaryPath)
			assert.Equal(t, "NEW-PRIMARY", string(mustReadFile(t, primaryPath)))
			companionPath := filepath.Join(dir, test.wantCompanion)
			assertCompanionPrimary(t, companionPath, primaryPath)
			if test.existingPrimary {
				assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
			}
			if test.existingCompanion {
				assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, companion)))
			}
			requireNoTemporaryReports(t, dir)
		})
	}
}

func TestInformativePairAggregateSkipsOccupiedNumberedCandidates(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("BASE"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "report_001.json"), []byte("ONE"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "report_001-informative.json"), []byte("ONE-INFO"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "report_002.json"), []byte("TWO"), 0o600))

	out := newInformativeTestOutput(t, dest, reportAggregate, nil)
	primaryPath, err := out.write(&stringReporter{payload: "NEW"}, nil)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "report_003.json"), primaryPath)
	assert.Equal(t, "NEW", string(mustReadFile(t, primaryPath)))
	assertCompanionPrimary(t, filepath.Join(dir, "report_003-informative.json"), primaryPath)
	assert.Equal(t, "BASE", string(mustReadFile(t, dest)))
	assert.Equal(t, "ONE", string(mustReadFile(t, filepath.Join(dir, "report_001.json"))))
	assert.Equal(t, "ONE-INFO", string(mustReadFile(t, filepath.Join(dir, "report_001-informative.json"))))
	assert.Equal(t, "TWO", string(mustReadFile(t, filepath.Join(dir, "report_002.json"))))
}

func TestInformativePairAggregateConcurrentAllocation(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	out, err := newReportOutput(dest, reportAggregate, "", nil, ".", "", "")
	require.NoError(t, err)
	require.NoError(t, out.enableInformative(func(primary string) ([]byte, error) {
		return []byte(fmt.Sprintf("{\"primary\":%q}\n", primary)), nil
	}))

	const workers = 12
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		paths    []string
		failures []error
	)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			path, err := out.write(&stringReporter{payload: fmt.Sprintf("report-%d", i)}, nil)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failures = append(failures, err)
				return
			}
			paths = append(paths, path)
		}(i)
	}
	wg.Wait()
	require.Empty(t, failures)
	require.Len(t, paths, workers)

	seen := map[string]struct{}{}
	for _, path := range paths {
		_, duplicate := seen[path]
		require.False(t, duplicate, "duplicate primary report path %q", path)
		seen[path] = struct{}{}
		assert.Contains(t, string(mustReadFile(t, path)), "report-")
		companion := informativeCompanionPath(path)
		require.FileExists(t, companion)
		assertCompanionPrimary(t, companion, path)
	}
	assert.Len(t, paths, workers)
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairRejectRequiresBothAbsent(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")

	t.Run("primary exists", func(t *testing.T) {
		require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
		_, err := newReportOutput(dest, reportReject, "", nil, "", ".", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("companion exists", func(t *testing.T) {
		require.NoError(t, os.Remove(dest))
		require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))
		_, err := newReportOutput(dest, reportReject, "", nil, "", ".", "")
		require.NoError(t, err)
		out, err := newReportOutput(dest, reportReject, "", nil, "", ".", "")
		require.NoError(t, err)
		err = out.enableInformative(func(primary string) ([]byte, error) { return []byte("new"), nil })
		require.Error(t, err)
		assert.Contains(t, err.Error(), "informative report path")
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestInformativePairRejectWritesBothWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	out := newInformativeTestOutput(t, dest, reportReject, nil)
	primaryPath, err := out.write(&stringReporter{payload: "PRIMARY"}, nil)
	require.NoError(t, err)
	assert.Equal(t, dest, primaryPath)
	assert.Equal(t, "PRIMARY", string(mustReadFile(t, dest)))
	assertCompanionPrimary(t, informativeCompanionPath(dest), dest)
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairCompanionCollision(t *testing.T) {
	tests := []struct {
		name     string
		source   func(dir, companion string) string
		configs  func(companion string) []string
		baseline func(companion string) string
		ignore   func(companion string) string
		template func(companion string) string
	}{
		{
			name:    "configuration",
			source:  func(string, string) string { return "" },
			configs: func(companion string) []string { return []string{companion} },
		},
		{
			name:     "baseline",
			source:   func(string, string) string { return "" },
			baseline: func(companion string) string { return companion },
		},
		{
			name:   "ignore",
			source: func(string, string) string { return "" },
			ignore: func(companion string) string { return companion },
		},
		{
			name:     "template",
			source:   func(string, string) string { return "" },
			template: func(companion string) string { return companion },
		},
		{
			name:   "scan target",
			source: func(_, companion string) string { return companion },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			dest := filepath.Join(dir, "report.json")
			companion := filepath.Join(dir, "report-informative.json")
			require.NoError(t, os.WriteFile(companion, []byte("PROTECTED"), 0o600))

			var (
				source   string
				configs  []string
				baseline string
				ignore   = "."
				template string
			)
			if test.source != nil {
				source = test.source(dir, companion)
			}
			if test.configs != nil {
				configs = test.configs(companion)
			}
			if test.baseline != nil {
				baseline = test.baseline(companion)
			}
			if test.ignore != nil {
				ignore = test.ignore(companion)
			}
			if test.template != nil {
				template = test.template(companion)
			}

			out, err := newReportOutput(dest, reportOverwrite, source, configs, baseline, ignore, template)
			require.NoError(t, err)
			err = out.enableInformative(func(primary string) ([]byte, error) { return []byte("info"), nil })
			require.Error(t, err)
			assert.Contains(t, err.Error(), "collides with")
			assert.Equal(t, "PROTECTED", string(mustReadFile(t, companion)))
		})
	}
}

func TestInformativePairCompanionSymlinkRejected(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	target := filepath.Join(dir, "target.json")
	require.NoError(t, os.WriteFile(target, []byte("TARGET"), 0o600))
	if err := os.Symlink(target, companion); err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}

	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	err = out.enableInformative(func(string) ([]byte, error) { return []byte("info"), nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "symbolic link")
	assert.Equal(t, "TARGET", string(mustReadFile(t, target)))
}

func TestInformativePairRenderFailurePreservesDestinations(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, func(string) ([]byte, error) {
		return nil, errors.New("render failure")
	})
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "render failure")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, companion)))
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairPrimaryWriterFailurePreservesDestinations(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	_, err := out.write(&stringReporter{payload: "partial", err: errors.New("writer failed")}, nil)
	require.Error(t, err)
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, companion)))
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairCompanionPublicationFailure(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return errors.New("companion publish failed")
		}
		return os.Rename(oldPath, newPath)
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "companion publish failed")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, companion)))
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairOverwriteRollbackPreservesConcurrentCompanion(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		// Simulate an unrelated writer replacing the newly published
		// companion while primary publication is still in progress.
		if err := os.Remove(companion); err != nil {
			return err
		}
		if err := os.WriteFile(companion, []byte("CONCURRENT-INFORMATIVE"), 0o600); err != nil {
			return err
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "primary publish failed")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "CONCURRENT-INFORMATIVE", string(mustReadFile(t, companion)))

	backups, globErr := filepath.Glob(filepath.Join(dir, ".gitleaks-companion-backup-*"))
	require.NoError(t, globErr)
	require.Len(t, backups, 1, "the prior companion must remain recoverable")
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, backups[0])))
}

func TestInformativePairOverwriteRollbackPreservesConcurrentCompanionWithoutBackup(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		if err := os.Remove(companion); err != nil {
			return err
		}
		if err := os.WriteFile(companion, []byte("CONCURRENT-INFORMATIVE"), 0o600); err != nil {
			return err
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "primary publish failed")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "CONCURRENT-INFORMATIVE", string(mustReadFile(t, companion)))
	backups, globErr := filepath.Glob(filepath.Join(dir, ".gitleaks-companion-backup-*"))
	require.NoError(t, globErr)
	assert.Empty(t, backups, "no prior companion means no backup is needed")
}

func TestInformativePairOverwriteFailedRestoreKeepsBackup(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(oldPath, ".gitleaks-companion-backup-") {
			return errors.New("companion restore failed")
		}
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "companion restore failed")
	assert.Contains(t, err.Error(), "preserved backup")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.NoFileExists(t, companion)

	backups, globErr := filepath.Glob(filepath.Join(dir, ".gitleaks-companion-backup-*"))
	require.NoError(t, globErr)
	require.Len(t, backups, 1, "failed restoration must not delete the only recoverable copy")
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, backups[0])))
}

func TestInformativePairOverwritePrimaryFailureRestoresCompanion(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))
	require.NoError(t, os.WriteFile(companion, []byte("OLD-INFORMATIVE"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "primary publish failed")
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.Equal(t, "OLD-INFORMATIVE", string(mustReadFile(t, companion)))
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairOverwritePrimaryFailureRemovesNewCompanion(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	companion := filepath.Join(dir, "report-informative.json")
	require.NoError(t, os.WriteFile(dest, []byte("OLD-PRIMARY"), 0o600))

	out := newInformativeTestOutput(t, dest, reportOverwrite, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Equal(t, "OLD-PRIMARY", string(mustReadFile(t, dest)))
	assert.NoFileExists(t, companion)
	requireNoTemporaryReports(t, dir)
}

func TestInformativePairPrimaryPublicationFailureRollsBackNewCompanion(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	out := newInformativeTestOutput(t, dest, reportAggregate, nil)
	out.rename = func(oldPath, newPath string) error {
		if strings.Contains(newPath, "informative") {
			return os.Rename(oldPath, newPath)
		}
		return errors.New("primary publish failed")
	}
	_, err := out.write(&stringReporter{payload: "NEW-PRIMARY"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "primary publish failed")
	assert.NoFileExists(t, dest)
	assert.NoFileExists(t, informativeCompanionPath(dest))
	requireNoTemporaryReports(t, dir)
}

func assertCompanionPrimary(t *testing.T, companionPath, primaryPath string) {
	t.Helper()
	var payload struct {
		Primary string `json:"primary"`
	}
	require.NoError(t, json.Unmarshal(mustReadFile(t, companionPath), &payload))
	assert.Equal(t, primaryPath, payload.Primary)
}

func TestInformativePairPreScanRejectAndConsoleOnlyCompatibility(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")

	t.Run("console-only caller does not enable a companion", func(t *testing.T) {
		out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
		require.NoError(t, err)
		require.Nil(t, out.companion)
		primary, err := out.write(&stringReporter{payload: "PRIMARY"}, nil)
		require.NoError(t, err)
		assert.Equal(t, dest, primary)
		assert.NoFileExists(t, informativeCompanionPath(dest))
	})

	t.Run("stdout cannot have a companion", func(t *testing.T) {
		out, err := newReportOutput(report.StdoutReportPath, reportOverwrite, "", nil, "", ".", "")
		require.NoError(t, err)
		err = out.enableInformative(func(string) ([]byte, error) { return []byte("info"), nil })
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stdout")
	})
}
