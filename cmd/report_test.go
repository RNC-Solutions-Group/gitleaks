package cmd

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/report"
)

type failingReportWriter struct {
	payload string
	err     error
}

func (w *failingReportWriter) Write(out io.WriteCloser, findings []report.Finding) error {
	if w.payload != "" {
		if _, err := io.WriteString(out, w.payload); err != nil {
			return err
		}
	}
	return w.err
}

func TestValidateReportDestination(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	require.NoError(t, os.MkdirAll(source, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("secret"), 0o600))

	t.Run("new destination outside source is allowed", func(t *testing.T) {
		err := validateReportDestination(filepath.Join(root, "new.json"), "json", source, "", "", ".")
		require.NoError(t, err)
	})

	t.Run("existing destination outside source is allowed", func(t *testing.T) {
		dest := filepath.Join(root, "existing.json")
		require.NoError(t, os.WriteFile(dest, []byte("not a report"), 0o600))
		err := validateReportDestination(dest, "json", source, "", "", ".")
		require.NoError(t, err)
	})

	t.Run("existing input inside source is rejected", func(t *testing.T) {
		err := validateReportDestination(filepath.Join(source, "secret.txt"), "json", source, "", "", ".")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not look like a prior")
	})

	t.Run("valid prior report inside source is allowed", func(t *testing.T) {
		dest := filepath.Join(source, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("[]\n"), 0o600))
		err := validateReportDestination(dest, "json", source, "", "", ".")
		require.NoError(t, err)
	})

	t.Run("source file collision is rejected", func(t *testing.T) {
		sourceFile := filepath.Join(root, "single.txt")
		require.NoError(t, os.WriteFile(sourceFile, []byte("input"), 0o600))
		err := validateReportDestination(sourceFile, "json", sourceFile, "", "", ".")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "collides with the scan target")
	})

	t.Run("directory destination is rejected", func(t *testing.T) {
		err := validateReportDestination(source, "json", source, "", "", ".")
		require.Error(t, err)
	})

	t.Run("missing parent is rejected", func(t *testing.T) {
		err := validateReportDestination(filepath.Join(root, "missing", "report.json"), "json", source, "", "", ".")
		require.Error(t, err)
	})

	t.Run("config input collision is rejected", func(t *testing.T) {
		configPath := filepath.Join(root, "gitleaks.toml")
		require.NoError(t, os.WriteFile(configPath, []byte("title = \"x\""), 0o600))
		err := validateReportDestination(configPath, "json", source, configPath, "", ".")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "collides with scan input")
	})

	t.Run("symlink destination is rejected", func(t *testing.T) {
		target := filepath.Join(root, "target.json")
		link := filepath.Join(root, "link.json")
		require.NoError(t, os.WriteFile(target, []byte("target"), 0o600))
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlinks are not available: %v", err)
		}
		err := validateReportDestination(link, "json", source, "", "", ".")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "symbolic link")
	})
}

func TestLooksLikePriorReport(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		name    string
		format  string
		content string
		want    bool
	}{
		{name: "json empty report", format: "json", content: "[]\n", want: true},
		{name: "json object", format: "json", content: "{}\n", want: false},
		{name: "csv header", format: "csv", content: "RuleID,Commit,File,Secret\n", want: true},
		{name: "csv empty output", format: "csv", content: "", want: true},
		{name: "junit empty", format: "junit", content: `<testsuites><testsuite name="gitleaks" failures="0" tests="0" time=""></testsuite></testsuites>`, want: true},
		{name: "sarif empty", format: "sarif", content: `{"version":"2.1.0","runs":[]}`, want: true},
		{name: "template unknown", format: "template", content: "anything", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, strings.ReplaceAll(tc.name, " ", "-"))
			require.NoError(t, os.WriteFile(path, []byte(tc.content), 0o600))
			got, err := looksLikePriorReport(path, tc.format)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWriteReportAtomic(t *testing.T) {
	t.Run("successful report replaces destination", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("prior"), 0o600))

		require.NoError(t, writeReportAtomic(dest, &report.JsonReporter{}, []report.Finding{}))
		data, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
		if runtime.GOOS != "windows" {
			info, err := os.Stat(dest)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("write failure preserves destination and removes temp", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("prior"), 0o600))

		err := writeReportAtomic(dest, &failingReportWriter{
			payload: "partial",
			err:     errors.New("writer failed"),
		}, nil)
		require.Error(t, err)

		data, readErr := os.ReadFile(dest)
		require.NoError(t, readErr)
		assert.Equal(t, "prior", string(data))

		entries, readDirErr := os.ReadDir(dir)
		require.NoError(t, readDirErr)
		for _, entry := range entries {
			assert.NotContains(t, entry.Name(), ".report.json.tmp-")
		}
	})

	t.Run("missing directory fails without creating a destination", func(t *testing.T) {
		dir := t.TempDir()
		dest := filepath.Join(dir, "missing", "report.json")
		err := writeReportAtomic(dest, &report.JsonReporter{}, nil)
		require.Error(t, err)
		_, statErr := os.Stat(dest)
		assert.True(t, os.IsNotExist(statErr))
	})
}

// TestReportDestinationChild is the child side of the CLI preservation tests.
// It must remain a no-op in a normal test run.
func TestReportDestinationChild(t *testing.T) {
	if os.Getenv("GITLEAKS_REPORT_CHILD") != "1" {
		return
	}

	args := []string{
		"gitleaks", "dir", os.Getenv("GITLEAKS_TEST_SOURCE"),
		"--no-banner",
		"--report-path", os.Getenv("GITLEAKS_TEST_REPORT"),
	}
	if format := os.Getenv("GITLEAKS_TEST_FORMAT"); format != "" {
		args = append(args, "--report-format", format)
	}
	if template := os.Getenv("GITLEAKS_TEST_TEMPLATE"); template != "" {
		args = append(args, "--report-template", template)
	}
	os.Args = args
	Execute()
}

func TestReportDestinationPreservedOnCLIError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	runChild := func(t *testing.T, source, dest, format, template string) ([]byte, error) {
		t.Helper()
		cmd := exec.Command(os.Args[0], "-test.run=^TestReportDestinationChild$")
		env := make([]string, 0, len(os.Environ())+5)
		for _, value := range os.Environ() {
			if strings.HasPrefix(strings.ToUpper(value), "GITLEAKS_") {
				continue
			}
			env = append(env, value)
		}
		env = append(env,
			"GITLEAKS_REPORT_CHILD=1",
			"GITLEAKS_TEST_SOURCE="+source,
			"GITLEAKS_TEST_REPORT="+dest,
			"GITLEAKS_TEST_FORMAT="+format,
			"GITLEAKS_TEST_TEMPLATE="+template,
		)
		cmd.Env = env
		_, err := cmd.CombinedOutput()
		data, readErr := os.ReadFile(dest)
		if readErr != nil {
			return nil, err
		}
		return data, err
	}

	t.Run("invalid format preserves existing report", func(t *testing.T) {
		source := t.TempDir()
		dest := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		data, err := runChild(t, source, dest, "not-a-format", "")
		require.Error(t, err)
		assert.Equal(t, "SENTINEL", string(data))
	})

	t.Run("invalid template preserves existing report", func(t *testing.T) {
		source := t.TempDir()
		dest := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		data, err := runChild(t, source, dest, "template", filepath.Join(t.TempDir(), "missing.tmpl"))
		require.Error(t, err)
		assert.Equal(t, "SENTINEL", string(data))
	})

	t.Run("input collision is rejected and input preserved", func(t *testing.T) {
		source := t.TempDir()
		input := filepath.Join(source, "secret.txt")
		require.NoError(t, os.WriteFile(input, []byte("aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\"\n"), 0o600))

		data, err := runChild(t, source, input, "json", "")
		require.Error(t, err)
		assert.Contains(t, string(data), "AKIAQ7Z3R5T7Y9U2I4O6")
	})
}
