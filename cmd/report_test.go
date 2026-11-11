package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/report"
)

type stringReporter struct {
	payload string
	err     error
}

func (r *stringReporter) Write(out io.WriteCloser, _ []report.Finding) error {
	if r.payload != "" {
		if _, err := io.WriteString(out, r.payload); err != nil {
			return err
		}
	}
	return r.err
}

func policyTestCommand(t *testing.T, args ...string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("overwrite", false, "")
	cmd.Flags().Bool("aggregate", false, "")
	require.NoError(t, cmd.ParseFlags(args))
	return cmd
}

func TestResolveReportMode(t *testing.T) {
	configs := []struct {
		name          string
		autoOverwrite bool
		aggregation   bool
		wantDefault   reportMode
	}{
		{name: "true false", autoOverwrite: true, aggregation: false, wantDefault: reportOverwrite},
		{name: "true true", autoOverwrite: true, aggregation: true, wantDefault: reportAggregate},
		{name: "false true", autoOverwrite: false, aggregation: true, wantDefault: reportAggregate},
		{name: "false false", autoOverwrite: false, aggregation: false, wantDefault: reportReject},
	}

	for _, testConfig := range configs {
		cfg := config.Config{
			AutoOverwrite: testConfig.autoOverwrite,
			Aggregation:   testConfig.aggregation,
		}
		t.Run(testConfig.name, func(t *testing.T) {
			t.Run("no flags uses configuration", func(t *testing.T) {
				mode, err := resolveReportMode(cfg, policyTestCommand(t))
				require.NoError(t, err)
				assert.Equal(t, testConfig.wantDefault, mode)
			})
			t.Run("--overwrite overrides configuration", func(t *testing.T) {
				mode, err := resolveReportMode(cfg, policyTestCommand(t, "--overwrite"))
				require.NoError(t, err)
				assert.Equal(t, reportOverwrite, mode)
			})
			t.Run("--aggregate overrides configuration", func(t *testing.T) {
				mode, err := resolveReportMode(cfg, policyTestCommand(t, "--aggregate"))
				require.NoError(t, err)
				assert.Equal(t, reportAggregate, mode)
			})
			t.Run("both flags conflict", func(t *testing.T) {
				_, err := resolveReportMode(cfg, policyTestCommand(t, "--overwrite", "--aggregate"))
				require.Error(t, err)
				assert.Contains(t, err.Error(), "--overwrite")
				assert.Contains(t, err.Error(), "--aggregate")
			})
		})
	}

	t.Run("explicit false flag does not force behavior", func(t *testing.T) {
		mode, err := resolveReportMode(
			config.Config{AutoOverwrite: false, Aggregation: false},
			policyTestCommand(t, "--overwrite=false"),
		)
		require.NoError(t, err)
		assert.Equal(t, reportReject, mode)

		mode, err = resolveReportMode(
			config.Config{AutoOverwrite: true, Aggregation: true},
			policyTestCommand(t, "--aggregate=false"),
		)
		require.NoError(t, err)
		assert.Equal(t, reportAggregate, mode)
	})
}

func TestNewReportOutput(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	require.NoError(t, os.Mkdir(source, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
	input := filepath.Join(source, "input.txt")

	t.Run("absent destination is accepted by every mode", func(t *testing.T) {
		for _, mode := range []reportMode{reportOverwrite, reportAggregate, reportReject} {
			_, err := newReportOutput(filepath.Join(root, "new-"+mode.String()+".json"), mode, source, nil, "", ".", "")
			require.NoError(t, err, mode)
		}
	})

	t.Run("existing destination is accepted by overwrite and aggregate", func(t *testing.T) {
		dest := filepath.Join(root, "existing.json")
		require.NoError(t, os.WriteFile(dest, []byte("UNRELATED"), 0o600))
		for _, mode := range []reportMode{reportOverwrite, reportAggregate} {
			_, err := newReportOutput(dest, mode, source, nil, "", ".", "")
			require.NoError(t, err, mode)
		}
	})

	t.Run("existing destination is rejected without overwrite or aggregation", func(t *testing.T) {
		dest := filepath.Join(root, "reject.json")
		require.NoError(t, os.WriteFile(dest, []byte("UNRELATED"), 0o600))
		_, err := newReportOutput(dest, reportReject, source, nil, "", ".", "")
		require.Error(t, err)
		assert.True(t, containsPath(err.Error(), dest), err.Error())
		assert.Contains(t, err.Error(), "auto-overwrite")
		assert.Contains(t, err.Error(), "aggregation")
		assert.Contains(t, err.Error(), "--overwrite")
		assert.Contains(t, err.Error(), "--aggregate")
	})

	t.Run("existing in-tree destination is no longer a special case", func(t *testing.T) {
		dest := filepath.Join(source, "in-tree.json")
		require.NoError(t, os.WriteFile(dest, []byte("UNRELATED"), 0o600))
		for _, mode := range []reportMode{reportOverwrite, reportAggregate} {
			_, err := newReportOutput(dest, mode, source, nil, "", ".", "")
			require.NoError(t, err, mode)
		}
	})

	t.Run("source file collision is rejected", func(t *testing.T) {
		for _, mode := range []reportMode{reportOverwrite, reportAggregate, reportReject} {
			_, err := newReportOutput(input, mode, input, nil, "", ".", "")
			require.Error(t, err, mode)
			assert.Contains(t, err.Error(), "collides with the scan target")
		}
	})

	t.Run("scan input collisions are rejected", func(t *testing.T) {
		configPath := filepath.Join(root, "gitleaks.toml")
		baselinePath := filepath.Join(root, "baseline.json")
		templatePath := filepath.Join(root, "report.tmpl")
		ignoreDir := filepath.Join(root, "ignore")
		require.NoError(t, os.Mkdir(ignoreDir, 0o700))
		ignorePath := filepath.Join(ignoreDir, ".gitleaksignore")
		for _, path := range []string{configPath, baselinePath, templatePath, ignorePath} {
			require.NoError(t, os.WriteFile(path, []byte("input"), 0o600))
		}

		for _, test := range []struct {
			name       string
			reportPath string
			configPath string
			baseline   string
			ignore     string
			template   string
		}{
			{name: "config", reportPath: configPath, configPath: configPath, ignore: "."},
			{name: "baseline", reportPath: baselinePath, baseline: baselinePath, ignore: "."},
			{name: "template", reportPath: templatePath, template: templatePath, ignore: "."},
			{name: "ignore", reportPath: ignorePath, ignore: ignoreDir},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := newReportOutput(
					test.reportPath,
					reportOverwrite,
					source,
					[]string{test.configPath},
					test.baseline,
					test.ignore,
					test.template,
				)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "collides with scan input")
			})
		}
	})

	t.Run("symlink destination is rejected", func(t *testing.T) {
		target := filepath.Join(root, "target.json")
		link := filepath.Join(root, "link.json")
		require.NoError(t, os.WriteFile(target, []byte("target"), 0o600))
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlinks are not available: %v", err)
		}
		for _, mode := range []reportMode{reportOverwrite, reportAggregate, reportReject} {
			_, err := newReportOutput(link, mode, source, nil, "", ".", "")
			require.Error(t, err, mode)
			assert.Contains(t, err.Error(), "symbolic link")
		}
	})

	t.Run("directory destination is rejected", func(t *testing.T) {
		dir := filepath.Join(root, "destination-dir")
		require.NoError(t, os.Mkdir(dir, 0o700))
		_, err := newReportOutput(dir, reportOverwrite, source, nil, "", ".", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is a directory")
	})

	t.Run("missing parent is rejected before scanning", func(t *testing.T) {
		_, err := newReportOutput(filepath.Join(root, "missing", "report.json"), reportOverwrite, source, nil, "", ".", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not accessible")
	})
}

func TestNumberedReportPath(t *testing.T) {
	tests := []struct {
		path   string
		number int
		want   string
	}{
		{path: filepath.Join("reports", "report.json"), number: 1, want: filepath.Join("reports", "report_001.json")},
		{path: filepath.Join("reports", "report.csv"), number: 12, want: filepath.Join("reports", "report_012.csv")},
		{path: filepath.Join("reports", "archive.tar.gz"), number: 2, want: filepath.Join("reports", "archive.tar_002.gz")},
		{path: filepath.Join("reports", "report"), number: 3, want: filepath.Join("reports", "report_003")},
		{path: filepath.Join("reports", ".report"), number: 4, want: filepath.Join("reports", ".report_004")},
		{path: filepath.Join("unicodé", "报告.json"), number: 5, want: filepath.Join("unicodé", "报告_005.json")},
		{path: filepath.Join("reports", "report.1.out"), number: 1000, want: filepath.Join("reports", "report.1_1000.out")},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, numberedReportPath(test.path, test.number), test.path)
	}
}

func TestReportOutputWrite(t *testing.T) {
	newOutput := func(t *testing.T, path string, mode reportMode) *reportOutput {
		t.Helper()
		out, err := newReportOutput(path, mode, "", nil, "", ".", "")
		require.NoError(t, err)
		return out
	}

	t.Run("overwrite replaces an existing destination", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("PRIOR"), 0o600))

		outputPath, err := newOutput(t, dest, reportOverwrite).write(&report.JsonReporter{}, []report.Finding{})
		require.NoError(t, err)
		assert.Equal(t, dest, outputPath)
		data, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
		if runtime.GOOS != "windows" {
			info, statErr := os.Stat(dest)
			require.NoError(t, statErr)
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("overwrite writes a new destination", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "report.json")
		outputPath, err := newOutput(t, dest, reportOverwrite).write(&report.JsonReporter{}, []report.Finding{})
		require.NoError(t, err)
		assert.Equal(t, dest, outputPath)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	})

	t.Run("reject writes an absent destination", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "report.json")
		outputPath, err := newOutput(t, dest, reportReject).write(&report.JsonReporter{}, []report.Finding{})
		require.NoError(t, err)
		assert.Equal(t, dest, outputPath)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	})

	t.Run("reject preserves a destination that appeared after validation", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		out := newOutput(t, dest, reportReject)
		require.NoError(t, os.WriteFile(dest, []byte("APPEARED"), 0o600))

		_, err := out.write(&stringReporter{payload: "new"}, nil)
		require.Error(t, err)
		assert.Equal(t, "APPEARED", string(mustReadFile(t, dest)))
		requireNoTemporaryReports(t, dir)
	})

	t.Run("aggregation writes the requested path when absent", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "report.json")
		outputPath, err := newOutput(t, dest, reportAggregate).write(&stringReporter{payload: "one"}, nil)
		require.NoError(t, err)
		assert.Equal(t, dest, outputPath)
		assert.Equal(t, "one", string(mustReadFile(t, dest)))
		if runtime.GOOS != "windows" {
			info, statErr := os.Stat(dest)
			require.NoError(t, statErr)
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("aggregation numbers successive reports", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		out := newOutput(t, dest, reportAggregate)

		for i, want := range []string{dest, filepath.Join(dir, "report_001.json"), filepath.Join(dir, "report_002.json")} {
			outputPath, err := out.write(&stringReporter{payload: fmt.Sprintf("report-%d", i)}, nil)
			require.NoError(t, err)
			assert.Equal(t, want, outputPath)
			assert.Equal(t, fmt.Sprintf("report-%d", i), string(mustReadFile(t, want)))
		}
	})

	t.Run("aggregation skips existing numbered files", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "report_001.json"), []byte("ONE"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "report_003.json"), []byte("THREE"), 0o600))

		outputPath, err := newOutput(t, dest, reportAggregate).write(&stringReporter{payload: "two"}, nil)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "report_002.json"), outputPath)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
		assert.Equal(t, "ONE", string(mustReadFile(t, filepath.Join(dir, "report_001.json"))))
		assert.Equal(t, "THREE", string(mustReadFile(t, filepath.Join(dir, "report_003.json"))))
		assert.Equal(t, "two", string(mustReadFile(t, outputPath)))
	})

	t.Run("aggregation preserves unrelated existing content", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte(`[{"not":"a report"}]`), 0o600))
		outputPath, err := newOutput(t, dest, reportAggregate).write(&stringReporter{payload: "new"}, nil)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "report_001.json"), outputPath)
		assert.Equal(t, `[{"not":"a report"}]`, string(mustReadFile(t, dest)))
	})

	t.Run("aggregation preserves nested unicode paths", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "nested", "unicodé")
		require.NoError(t, os.MkdirAll(dir, 0o700))
		dest := filepath.Join(dir, "报告.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		outputPath, err := newOutput(t, dest, reportAggregate).write(&stringReporter{payload: "new"}, nil)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "报告_001.json"), outputPath)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
		assert.Equal(t, "new", string(mustReadFile(t, outputPath)))
	})

	t.Run("writer failure preserves the destination", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("PRIOR"), 0o600))

		_, err := newOutput(t, dest, reportOverwrite).write(&stringReporter{
			payload: "partial",
			err:     errors.New("writer failed"),
		}, nil)
		require.Error(t, err)
		assert.Equal(t, "PRIOR", string(mustReadFile(t, dest)))
		requireNoTemporaryReports(t, dir)
	})

	t.Run("publication failure preserves the destination", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("PRIOR"), 0o600))
		out := newOutput(t, dest, reportOverwrite)
		out.rename = func(string, string) error { return errors.New("publish failed") }

		_, err := out.write(&stringReporter{payload: "new"}, nil)
		require.Error(t, err)
		assert.Equal(t, "PRIOR", string(mustReadFile(t, dest)))
		requireNoTemporaryReports(t, dir)
	})

	t.Run("aggregate publication failure removes the reservation", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
		out := newOutput(t, dest, reportAggregate)
		out.rename = func(string, string) error { return errors.New("publish failed") }

		_, err := out.write(&stringReporter{payload: "new"}, nil)
		require.Error(t, err)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
		assert.NoFileExists(t, filepath.Join(dir, "report_001.json"))
		requireNoTemporaryReports(t, dir)
	})

	t.Run("symlink appearing before publication is rejected", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		target := filepath.Join(dir, "target.json")
		require.NoError(t, os.WriteFile(target, []byte("TARGET"), 0o600))
		out := newOutput(t, dest, reportOverwrite)
		if err := os.Symlink(target, dest); err != nil {
			t.Skipf("symlinks are not available: %v", err)
		}

		_, err := out.write(&stringReporter{payload: "new"}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "symbolic link")
		assert.Equal(t, "TARGET", string(mustReadFile(t, target)))
		requireNoTemporaryReports(t, dir)
	})
}

func TestAggregateSupportsAllReporters(t *testing.T) {
	templatePath := filepath.Join(t.TempDir(), "report.tmpl")
	require.NoError(t, os.WriteFile(templatePath, []byte(`findings={{ len . }}`), 0o600))
	templateReporter, err := report.NewTemplateReporter(templatePath)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		reporter report.Reporter
	}{
		{name: "json", filename: "report.json", reporter: &report.JsonReporter{}},
		{name: "csv", filename: "report.csv", reporter: &report.CsvReporter{}},
		{name: "junit", filename: "report.xml", reporter: &report.JunitReporter{}},
		{name: "sarif", filename: "report.sarif", reporter: &report.SarifReporter{}},
		{name: "template", filename: "report.txt", reporter: templateReporter},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			dest := filepath.Join(dir, test.filename)
			require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
			out, err := newReportOutput(dest, reportAggregate, "", nil, "", ".", "")
			require.NoError(t, err)

			outputPath, err := out.write(test.reporter, []report.Finding{{RuleID: "test", Description: "test"}})
			require.NoError(t, err)
			assert.Equal(t, numberedReportPath(dest, 1), outputPath)
			assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
			assert.FileExists(t, outputPath)
		})
	}
}

func TestAggregateConcurrentAllocation(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
	out, err := newReportOutput(dest, reportAggregate, "", nil, "", ".", "")
	require.NoError(t, err)

	const workers = 16
	var (
		wg       sync.WaitGroup
		errs     = make(chan error, workers)
		paths    = make(chan string, workers)
		contents = make(chan string, workers)
	)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			content := fmt.Sprintf("report-%d", i)
			path, err := out.write(&stringReporter{payload: content}, nil)
			if err != nil {
				errs <- err
				return
			}
			paths <- path
			contents <- content
		}(i)
	}
	wg.Wait()
	close(errs)
	close(paths)
	close(contents)
	for err := range errs {
		require.NoError(t, err)
	}

	seenPaths := map[string]struct{}{}
	for path := range paths {
		require.NotEqual(t, dest, path)
		_, duplicate := seenPaths[path]
		require.False(t, duplicate, "duplicate aggregate path %q", path)
		seenPaths[path] = struct{}{}
	}
	seenContents := map[string]struct{}{}
	for content := range contents {
		_, duplicate := seenContents[content]
		require.False(t, duplicate, "duplicate aggregate content %q", content)
		seenContents[content] = struct{}{}
	}

	assert.Equal(t, workers, len(seenPaths))
	assert.Equal(t, workers, len(seenContents))
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
	for path := range seenPaths {
		assert.Contains(t, string(mustReadFile(t, path)), "report-")
	}
	requireNoTemporaryReports(t, dir)
}

func TestReportWriteUnwritableDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permission semantics differ on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses directory permissions")
	}

	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	require.NoError(t, err)
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	_, err = out.write(&report.JsonReporter{}, nil)
	require.Error(t, err)
	assert.NoFileExists(t, dest)
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

// containsPath handles zerolog's quoting of strings that contain path
// separators or other characters requiring an escape sequence.
func containsPath(output, path string) bool {
	return strings.Contains(output, path) || strings.Contains(output, strconv.Quote(path))
}

func requireNoTemporaryReports(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, entry := range entries {
		assert.NotContains(t, entry.Name(), ".tmp-")
	}
}

// TestReportCLIChild is the child side of the CLI integration tests.
func TestReportCLIChild(t *testing.T) {
	if os.Getenv("GITLEAKS_REPORT_CHILD") != "1" {
		return
	}

	var args []string
	if err := json.Unmarshal([]byte(os.Getenv("GITLEAKS_TEST_ARGS")), &args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Args = append([]string{"gitleaks"}, args...)
	Execute()
	// Prevent the Go test harness from writing PASS to stdout, which would
	// contaminate the machine-readable report when --report-path=- is used.
	os.Exit(0)
}

// TestReportPartialScanChild verifies that a propagated scan error prevents
// report publication even when the destination and writer are configured.
func TestReportPartialScanChild(t *testing.T) {
	if os.Getenv("GITLEAKS_PARTIAL_SCAN_CHILD") != "1" {
		return
	}

	dest := os.Getenv("GITLEAKS_TEST_REPORT")
	output, err := newReportOutput(dest, reportOverwrite, "", nil, "", ".", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	detector := &detect.Detector{
		ReportPath: dest,
		Reporter:   &report.JsonReporter{},
		ReportWriter: func(findings []report.Finding) (string, error) {
			return output.write(&report.JsonReporter{}, findings)
		},
	}
	diagnosticsManager = &DiagnosticsManager{}
	findingSummaryAndExit(detector, nil, 1, time.Now(), errors.New("synthetic scan failure"))
	os.Exit(0)
}

func runPartialScanChild(t *testing.T, dest string) reportCLIResult {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestReportPartialScanChild$")
	env := make([]string, 0, len(os.Environ())+2)
	for _, value := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(value), "GITLEAKS_") {
			continue
		}
		env = append(env, value)
	}
	env = append(env, "GITLEAKS_PARTIAL_SCAN_CHILD=1", "GITLEAKS_TEST_REPORT="+dest)
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := reportCLIResult{stdout: stdout.String(), stderr: stderr.String()}
	if err == nil {
		return result
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result
	}
	t.Fatalf("could not run partial scan child: %v", err)
	return result
}

func TestPartialScanDoesNotWriteReport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("PRIOR"), 0o600))

	result := runPartialScanChild(t, dest)
	require.Equal(t, 1, result.exitCode)
	assert.Equal(t, "PRIOR", string(mustReadFile(t, dest)))
	assert.Empty(t, numberedReports(t, dir, "report"))
	requireNoTemporaryReports(t, dir)
	assert.Contains(t, result.stderr, "report not written because the scan was incomplete")
}

type reportCLIResult struct {
	stdout   string
	stderr   string
	exitCode int
}

func runReportCLI(t *testing.T, args ...string) reportCLIResult {
	t.Helper()
	return runReportCLIWithEnv(t, nil, args...)
}

func runReportCLIWithEnv(t *testing.T, extraEnv map[string]string, args ...string) reportCLIResult {
	t.Helper()
	encoded, err := json.Marshal(args)
	require.NoError(t, err)

	cmd := exec.Command(os.Args[0], "-test.run=^TestReportCLIChild$")
	env := make([]string, 0, len(os.Environ())+2+len(extraEnv))
	extraKeys := make(map[string]struct{}, len(extraEnv))
	for name := range extraEnv {
		extraKeys[strings.ToUpper(name)] = struct{}{}
	}
	for _, value := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(value), "GITLEAKS_") {
			continue
		}
		name := strings.SplitN(value, "=", 2)[0]
		if _, replaced := extraKeys[strings.ToUpper(name)]; replaced {
			continue
		}
		env = append(env, value)
	}
	env = append(env, "GITLEAKS_REPORT_CHILD=1", "GITLEAKS_TEST_ARGS="+string(encoded))
	for name, value := range extraEnv {
		env = append(env, name+"="+value)
	}
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	result := reportCLIResult{stdout: stdout.String(), stderr: stderr.String()}
	if err == nil {
		return result
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result
	}
	t.Fatalf("could not run report CLI child: %v", err)
	return result
}

func reportPolicyConfig(t *testing.T, autoOverwrite, aggregation bool) string {
	t.Helper()
	content := fmt.Sprintf("title = \"report policy\"\nauto-overwrite = %t\naggregation = %t\n", autoOverwrite, aggregation)
	path := filepath.Join(t.TempDir(), "gitleaks.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func reportCommandArgs(source, dest, configPath string, extra ...string) []string {
	args := []string{
		"dir", source,
		"--no-banner",
		"--report-path", dest,
		"--report-format", "json",
	}
	if configPath != "" {
		args = append(args, "--config", configPath)
	}
	return append(args, extra...)
}

func cleanScanDir(t *testing.T) string {
	t.Helper()
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "readme.txt"), []byte("ordinary text\n"), 0o600))
	return source
}

func numberedReports(t *testing.T, dir, base string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, base+"_*.json"))
	require.NoError(t, err)
	return matches
}

func TestCLIReportConfigurationMatrix(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	configs := []struct {
		name          string
		autoOverwrite bool
		aggregation   bool
		existingMode  reportMode
	}{
		{name: "true false", autoOverwrite: true, aggregation: false, existingMode: reportOverwrite},
		{name: "true true", autoOverwrite: true, aggregation: true, existingMode: reportAggregate},
		{name: "false true", autoOverwrite: false, aggregation: true, existingMode: reportAggregate},
		{name: "false false", autoOverwrite: false, aggregation: false, existingMode: reportReject},
	}

	for _, testConfig := range configs {
		t.Run(testConfig.name, func(t *testing.T) {
			configPath := reportPolicyConfig(t, testConfig.autoOverwrite, testConfig.aggregation)

			t.Run("absent destination", func(t *testing.T) {
				source := cleanScanDir(t)
				dir := t.TempDir()
				dest := filepath.Join(dir, "report.json")

				result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
				require.Equal(t, 0, result.exitCode, result.stderr)
				assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
				assert.Empty(t, numberedReports(t, dir, "report"))
				assert.Contains(t, result.stderr, "report written")
				assert.True(t, containsPath(result.stderr, dest), result.stderr)
			})

			t.Run("existing destination", func(t *testing.T) {
				source := cleanScanDir(t)
				dir := t.TempDir()
				dest := filepath.Join(dir, "report.json")
				require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

				result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
				switch testConfig.existingMode {
				case reportOverwrite:
					require.Equal(t, 0, result.exitCode, result.stderr)
					assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
					assert.Empty(t, numberedReports(t, dir, "report"))
				case reportAggregate:
					require.Equal(t, 0, result.exitCode, result.stderr)
					assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
					numbered := numberedReports(t, dir, "report")
					require.Len(t, numbered, 1)
					assert.JSONEq(t, "[]", string(mustReadFile(t, numbered[0])))
					assert.True(t, containsPath(result.stderr, numbered[0]), result.stderr)
				case reportReject:
					require.NotEqual(t, 0, result.exitCode)
					assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
					assert.Empty(t, numberedReports(t, dir, "report"))
					assert.Contains(t, result.stderr, "already exists")
					assert.Contains(t, result.stderr, "--overwrite")
					assert.Contains(t, result.stderr, "--aggregate")
					assert.NotContains(t, result.stderr, "scanned ~")
				}
			})
		})
	}
}

func TestCLIReportOverrides(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	configs := []struct {
		name          string
		autoOverwrite bool
		aggregation   bool
	}{
		{name: "true false", autoOverwrite: true, aggregation: false},
		{name: "true true", autoOverwrite: true, aggregation: true},
		{name: "false true", autoOverwrite: false, aggregation: true},
		{name: "false false", autoOverwrite: false, aggregation: false},
	}

	for _, testConfig := range configs {
		t.Run(testConfig.name, func(t *testing.T) {
			configPath := reportPolicyConfig(t, testConfig.autoOverwrite, testConfig.aggregation)

			t.Run("--overwrite", func(t *testing.T) {
				source := cleanScanDir(t)
				dir := t.TempDir()
				dest := filepath.Join(dir, "report.json")
				require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

				args := reportCommandArgs(source, dest, configPath)
				args = append(args, "--overwrite")
				result := runReportCLI(t, args...)
				require.Equal(t, 0, result.exitCode, result.stderr)
				assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
				assert.Empty(t, numberedReports(t, dir, "report"))
			})

			t.Run("--aggregate", func(t *testing.T) {
				source := cleanScanDir(t)
				dir := t.TempDir()
				dest := filepath.Join(dir, "report.json")
				require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

				args := reportCommandArgs(source, dest, configPath, "--aggregate")
				result := runReportCLI(t, args...)
				require.Equal(t, 0, result.exitCode, result.stderr)
				assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
				numbered := numberedReports(t, dir, "report")
				require.Len(t, numbered, 1)
				assert.JSONEq(t, "[]", string(mustReadFile(t, numbered[0])))
			})
		})
	}

	t.Run("persistent flags before subcommand and report options", func(t *testing.T) {
		source := cleanScanDir(t)
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
		configPath := reportPolicyConfig(t, false, false)

		result := runReportCLI(t,
			"--overwrite",
			"dir", source,
			"--no-banner",
			"--config", configPath,
			"--report-path", dest,
			"--report-format", "json",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	})
}

func TestCLIReportFlagConflictAndStdout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("both flags fail before scanning and preserve the destination", func(t *testing.T) {
		source := cleanScanDir(t)
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
		configPath := reportPolicyConfig(t, true, false)

		args := reportCommandArgs(source, dest, configPath)
		args = append(args, "--overwrite", "--aggregate")
		result := runReportCLI(t, args...)
		require.NotEqual(t, 0, result.exitCode)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
		assert.Empty(t, numberedReports(t, dir, "report"))
		assert.Contains(t, result.stderr, "cannot be used together")
		assert.NotContains(t, result.stderr, "scanned ~")
	})

	t.Run("stdout report is not contaminated by diagnostics", func(t *testing.T) {
		source := cleanScanDir(t)
		configPath := reportPolicyConfig(t, false, false)

		result := runReportCLI(t,
			"dir", source,
			"--no-banner",
			"--config", configPath,
			"--report-path", "-",
			"--report-format", "json",
		)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", result.stdout)
		assert.NotContains(t, result.stdout, "report written")
		assert.Contains(t, result.stderr, "no leaks found")
	})
}

func TestCLIReportAggregationSequence(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := cleanScanDir(t)
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	configPath := reportPolicyConfig(t, false, true)

	for i, want := range []string{dest, filepath.Join(dir, "report_001.json"), filepath.Join(dir, "report_002.json")} {
		result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, want)))
		assert.True(t, containsPath(result.stderr, want), result.stderr)
		t.Logf("scan %d wrote %s", i+1, want)
	}
	assert.Len(t, numberedReports(t, dir, "report"), 2)
}

func TestCLIReportInTreeDestination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	t.Run("existing in-tree report is replaced by default", func(t *testing.T) {
		source := cleanScanDir(t)
		dest := filepath.Join(source, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		result := runReportCLI(t, reportCommandArgs(source, dest, "")...)
		require.Equal(t, 0, result.exitCode, result.stderr)
		assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	})

	t.Run("removed in-tree flag is rejected as unknown", func(t *testing.T) {
		source := cleanScanDir(t)
		dest := filepath.Join(source, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		args := reportCommandArgs(source, dest, "")
		args = append(args, "--allow-in-tree-report-overwrite")
		result := runReportCLI(t, args...)
		require.Equal(t, 126, result.exitCode)
		assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
		assert.Contains(t, result.stderr, "unknown flag")
	})
}

func TestCLIReportValidationPreservesDestination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	configs := []struct {
		name          string
		autoOverwrite bool
		aggregation   bool
	}{
		{name: "true false", autoOverwrite: true, aggregation: false},
		{name: "true true", autoOverwrite: true, aggregation: true},
		{name: "false true", autoOverwrite: false, aggregation: true},
		{name: "false false", autoOverwrite: false, aggregation: false},
	}

	failures := []struct {
		name string
		args func(t *testing.T, source, dest string) []string
	}{
		{
			name: "invalid format",
			args: func(t *testing.T, source, dest string) []string {
				return []string{"dir", source, "--no-banner", "--report-path", dest, "--report-format", "not-a-format"}
			},
		},
		{
			name: "invalid template",
			args: func(t *testing.T, source, dest string) []string {
				return []string{
					"dir", source, "--no-banner", "--report-path", dest,
					"--report-format", "template",
					"--report-template", filepath.Join(t.TempDir(), "missing.tmpl"),
				}
			},
		},
		{
			name: "template execution failure",
			args: func(t *testing.T, source, dest string) []string {
				templatePath := filepath.Join(t.TempDir(), "failing.tmpl")
				require.NoError(t, os.WriteFile(templatePath, []byte(`{{ fail "template failed" }}`), 0o600))
				return []string{
					"dir", source, "--no-banner", "--report-path", dest,
					"--report-format", "template",
					"--report-template", templatePath,
				}
			},
		},
	}

	for _, testConfig := range configs {
		t.Run(testConfig.name, func(t *testing.T) {
			configPath := reportPolicyConfig(t, testConfig.autoOverwrite, testConfig.aggregation)
			for _, failure := range failures {
				t.Run(failure.name, func(t *testing.T) {
					source := cleanScanDir(t)
					dir := t.TempDir()
					dest := filepath.Join(dir, "report.json")
					require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

					args := failure.args(t, source, dest)
					args = append(args, "--config", configPath)
					result := runReportCLI(t, args...)
					require.NotEqual(t, 0, result.exitCode)
					assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
					assert.Empty(t, numberedReports(t, dir, "report"))
					requireNoTemporaryReports(t, dir)
				})
			}
		})
	}
}

func TestCLIReportWithCustomRule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(source, "creds.txt"),
		[]byte(`aws_access_key_id = "AKIAQ7Z3R5T7Y9U2I4O6"`+"\n"),
		0o600,
	))
	configPath := filepath.Join(t.TempDir(), "rules.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
title = "report policy with rules"
auto-overwrite = true
aggregation = true

[[rules]]
id = "test-rule"
description = "test rule"
regex = '''AKIA[0-9A-Z]{16}'''
keywords = ["akia"]
`), 0o600))

	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

	result := runReportCLI(t, reportCommandArgs(source, dest, configPath)...)
	require.Equal(t, 1, result.exitCode, result.stderr)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
	numbered := numberedReports(t, dir, "report")
	require.Len(t, numbered, 1)
	assert.Contains(t, string(mustReadFile(t, numbered[0])), `"RuleID": "test-rule"`)
}

func TestCLIDefaultReportBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess tests in short mode")
	}

	source := cleanScanDir(t)
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")
	require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

	result := runReportCLI(t, reportCommandArgs(source, dest, "")...)
	require.Equal(t, 0, result.exitCode, result.stderr)
	assert.JSONEq(t, "[]", string(mustReadFile(t, dest)))
	assert.Empty(t, numberedReports(t, dir, "report"))
}
