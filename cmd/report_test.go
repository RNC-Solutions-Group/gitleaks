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

type appearingReportWriter struct {
	path string
}

func (w *appearingReportWriter) Write(out io.WriteCloser, findings []report.Finding) error {
	if _, err := io.WriteString(out, "temporary"); err != nil {
		return err
	}
	return os.WriteFile(w.path, []byte("appeared"), 0o600)
}

func requireDestination(t *testing.T, reportPath, source string, allowInTree bool) reportDestination {
	t.Helper()
	dest, err := validateReportDestination(reportPath, source, "", "", ".", allowInTree)
	require.NoError(t, err)
	return dest
}

func TestValidateReportDestination(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	require.NoError(t, os.MkdirAll(source, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("secret"), 0o600))

	t.Run("new destination outside source is allowed", func(t *testing.T) {
		_, err := validateReportDestination(filepath.Join(root, "new.json"), source, "", "", ".", false)
		require.NoError(t, err)
	})

	t.Run("existing destination outside source is allowed", func(t *testing.T) {
		dest := filepath.Join(root, "existing.json")
		require.NoError(t, os.WriteFile(dest, []byte("not a report"), 0o600))
		_, err := validateReportDestination(dest, source, "", "", ".", false)
		require.NoError(t, err)
	})

	t.Run("new destination inside source is allowed", func(t *testing.T) {
		_, err := validateReportDestination(filepath.Join(source, "new.json"), source, "", "", ".", false)
		require.NoError(t, err)
	})

	t.Run("existing input inside source is rejected", func(t *testing.T) {
		_, err := validateReportDestination(filepath.Join(source, "secret.txt"), source, "", "", ".", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "allow-in-tree-report-overwrite")
	})

	t.Run("existing input inside source requires explicit authorization", func(t *testing.T) {
		dest := filepath.Join(source, "existing.json")
		require.NoError(t, os.WriteFile(dest, []byte("[]\n"), 0o600))
		_, err := validateReportDestination(dest, source, "", "", ".", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "allow-in-tree-report-overwrite")
	})

	t.Run("authorized in-tree replacement with a second input is allowed", func(t *testing.T) {
		dest := filepath.Join(source, "authorized.json")
		require.NoError(t, os.WriteFile(dest, []byte("unrelated"), 0o600))
		_, err := validateReportDestination(dest, source, "", "", ".", true)
		require.NoError(t, err)
	})

	t.Run("sole in-tree file cannot be replaced even with authorization", func(t *testing.T) {
		sole := filepath.Join(root, "sole")
		require.NoError(t, os.Mkdir(sole, 0o700))
		dest := filepath.Join(sole, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("[]\n"), 0o600))
		_, err := validateReportDestination(dest, sole, "", "", ".", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sole file")
	})

	t.Run("source file collision is rejected", func(t *testing.T) {
		sourceFile := filepath.Join(root, "single.txt")
		require.NoError(t, os.WriteFile(sourceFile, []byte("input"), 0o600))
		_, err := validateReportDestination(sourceFile, sourceFile, "", "", ".", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "collides with the scan target")
	})

	t.Run("directory destination is rejected", func(t *testing.T) {
		_, err := validateReportDestination(source, source, "", "", ".", true)
		require.Error(t, err)
	})

	t.Run("missing parent is rejected", func(t *testing.T) {
		_, err := validateReportDestination(filepath.Join(root, "missing", "report.json"), source, "", "", ".", true)
		require.Error(t, err)
	})

	t.Run("config input collision is rejected with or without authorization", func(t *testing.T) {
		configPath := filepath.Join(root, "gitleaks.toml")
		require.NoError(t, os.WriteFile(configPath, []byte("title = \"x\""), 0o600))
		for _, allow := range []bool{false, true} {
			_, err := validateReportDestination(configPath, source, configPath, "", ".", allow)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "collides with scan input")
		}
	})

	t.Run("baseline input collision is rejected with authorization", func(t *testing.T) {
		baselinePath := filepath.Join(root, "baseline.json")
		require.NoError(t, os.WriteFile(baselinePath, []byte("[]"), 0o600))
		_, err := validateReportDestination(baselinePath, source, "", baselinePath, ".", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "collides with scan input")
	})

	t.Run("symlink destination is rejected with and without authorization", func(t *testing.T) {
		target := filepath.Join(root, "target.json")
		link := filepath.Join(root, "link.json")
		require.NoError(t, os.WriteFile(target, []byte("target"), 0o600))
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlinks are not available: %v", err)
		}
		for _, allow := range []bool{false, true} {
			_, err := validateReportDestination(link, source, "", "", ".", allow)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "symbolic link")
		}
	})
}

// TestNoContentHeuristicsAuthorizeOverwrite verifies that report-shaped
// content does not turn an in-tree input into a replaceable destination.
func TestNoContentHeuristicsAuthorizeOverwrite(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	require.NoError(t, os.Mkdir(source, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("real input"), 0o600))

	cases := map[string]string{
		"unrelated JSON array":      `[{"name":"application-record","credential":"synthetic-credential-example"}]`,
		"JSON array unknown fields": `[{"RuleID":"not-a-gitleaks-finding","name":"record","credential":"synthetic"}]`,
		"empty JSON array":          `[]`,
		"valid prior Gitleaks JSON": `[{"RuleID":"rule","Description":"desc","File":"f","Secret":"AKIAQ7Z3R5T7Y9U2I4O6","StartLine":1,"EndLine":1,"StartColumn":1,"EndColumn":20}]`,
		"empty CSV":                 "",
		"CSV report-like header":    "RuleID,Commit,File,Secret\n",
		"JUnit-like XML":            `<testsuites><testsuite name="gitleaks"></testsuite></testsuites>`,
		"SARIF-like JSON":           `{"version":"2.1.0","runs":[]}`,
		"arbitrary text":            "not a report",
	}

	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			dest := filepath.Join(source, strings.ReplaceAll(name, " ", "-")+".out")
			require.NoError(t, os.WriteFile(dest, []byte(content), 0o600))
			_, err := validateReportDestination(dest, source, "", "", ".", false)
			require.Error(t, err, "content must not implicitly authorize replacement")
			assert.Contains(t, err.Error(), "allow-in-tree-report-overwrite")

			_, err = validateReportDestination(dest, source, "", "", ".", true)
			require.NoError(t, err, "explicit authorization should permit a non-sole in-tree replacement")
		})
	}
}

func TestWriteReportAtomic(t *testing.T) {
	t.Run("successful report replaces destination", func(t *testing.T) {
		dir := t.TempDir()
		destPath := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(destPath, []byte("prior"), 0o600))
		dest := requireDestination(t, destPath, "", false)

		require.NoError(t, writeReportAtomic(dest, &report.JsonReporter{}, []report.Finding{}))
		data, err := os.ReadFile(destPath)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
		if runtime.GOOS != "windows" {
			info, err := os.Stat(destPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
		}
	})

	t.Run("write failure preserves destination and removes temp", func(t *testing.T) {
		dir := t.TempDir()
		destPath := filepath.Join(dir, "report.json")
		require.NoError(t, os.WriteFile(destPath, []byte("prior"), 0o600))
		dest := requireDestination(t, destPath, "", false)

		err := writeReportAtomic(dest, &failingReportWriter{
			payload: "partial",
			err:     errors.New("writer failed"),
		}, nil)
		require.Error(t, err)

		data, readErr := os.ReadFile(destPath)
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
		destPath := filepath.Join(dir, "missing", "report.json")
		dest := reportDestination{path: destPath}
		err := writeReportAtomic(dest, &report.JsonReporter{}, nil)
		require.Error(t, err)
		_, statErr := os.Stat(destPath)
		assert.True(t, os.IsNotExist(statErr))
	})
}

func TestInTreeDestinationAppearingAfterValidationIsNotOverwritten(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
	destPath := filepath.Join(source, "report.json")
	dest := requireDestination(t, destPath, source, false)

	// A new in-tree destination is allowed. If it appears after validation and
	// during the writer's work, the second revalidation must preserve it.
	err := writeReportAtomic(dest, &appearingReportWriter{path: destPath}, nil)
	require.Error(t, err)
	data, readErr := os.ReadFile(destPath)
	require.NoError(t, readErr)
	assert.Equal(t, "appeared", string(data))
}

func TestAuthorizedInTreeReplacementRepeats(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
	destPath := filepath.Join(source, "report.json")
	require.NoError(t, os.WriteFile(destPath, []byte("unrelated"), 0o600))

	for i := 0; i < 2; i++ {
		dest, err := validateReportDestination(destPath, source, "", "", ".", true)
		require.NoError(t, err)
		require.NoError(t, writeReportAtomic(dest, &report.JsonReporter{}, []report.Finding{}))
		data, readErr := os.ReadFile(destPath)
		require.NoError(t, readErr)
		assert.JSONEq(t, "[]", string(data))
	}
}

// TestReportDestinationChild is the child side of the CLI destination tests.
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
	if os.Getenv("GITLEAKS_TEST_ALLOW") == "1" {
		args = append(args, "--allow-in-tree-report-overwrite")
	}
	os.Args = args
	Execute()
}

func TestReportDestinationCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	runChild := func(t *testing.T, source, dest, format, template string, allow bool) ([]byte, error) {
		t.Helper()
		cmd := exec.Command(os.Args[0], "-test.run=^TestReportDestinationChild$")
		env := make([]string, 0, len(os.Environ())+6)
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
			"GITLEAKS_TEST_ALLOW="+map[bool]string{true: "1", false: "0"}[allow],
		)
		cmd.Env = env
		_, err := cmd.CombinedOutput()
		data, readErr := os.ReadFile(dest)
		if readErr != nil {
			return nil, readErr
		}
		return data, err
	}

	t.Run("clean scan writes a new report", func(t *testing.T) {
		source := t.TempDir()
		dest := filepath.Join(t.TempDir(), "report.json")

		data, err := runChild(t, source, dest, "json", "", false)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
	})

	t.Run("new in-tree report is allowed", func(t *testing.T) {
		source := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
		dest := filepath.Join(source, "report.json")

		data, err := runChild(t, source, dest, "json", "", false)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
	})

	t.Run("existing in-tree report requires authorization", func(t *testing.T) {
		source := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
		dest := filepath.Join(source, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("UNRELATED"), 0o600))

		data, err := runChild(t, source, dest, "json", "", false)
		require.Error(t, err)
		assert.Equal(t, "UNRELATED", string(data))
	})

	t.Run("authorized in-tree replacement is explicit", func(t *testing.T) {
		source := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(source, "input.txt"), []byte("input"), 0o600))
		dest := filepath.Join(source, "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("UNRELATED"), 0o600))

		data, err := runChild(t, source, dest, "json", "", true)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(data))
	})

	t.Run("invalid format preserves existing report", func(t *testing.T) {
		source := t.TempDir()
		dest := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		data, err := runChild(t, source, dest, "not-a-format", "", false)
		require.Error(t, err)
		assert.Equal(t, "SENTINEL", string(data))
	})

	t.Run("invalid template preserves existing report", func(t *testing.T) {
		source := t.TempDir()
		dest := filepath.Join(t.TempDir(), "report.json")
		require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))

		data, err := runChild(t, source, dest, "template", filepath.Join(t.TempDir(), "missing.tmpl"), false)
		require.Error(t, err)
		assert.Equal(t, "SENTINEL", string(data))
	})

	t.Run("input collision is rejected and input preserved", func(t *testing.T) {
		source := t.TempDir()
		input := filepath.Join(source, "secret.txt")
		require.NoError(t, os.WriteFile(input, []byte("aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\"\n"), 0o600))

		data, err := runChild(t, source, input, "json", "", false)
		require.Error(t, err)
		assert.Contains(t, string(data), "AKIAQ7Z3R5T7Y9U2I4O6")
	})
}
