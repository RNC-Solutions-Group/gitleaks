package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreCommitHookManifest(t *testing.T) {
	content, err := os.ReadFile("../.pre-commit-hooks.yaml")
	require.NoError(t, err)

	// Avoid adding a YAML dependency for a small, fixed manifest. Split on the
	// list entries and verify the fields that determine filename forwarding.
	parts := strings.Split(string(content), "- id: ")
	require.GreaterOrEqual(t, len(parts), 4, "expected three hook definitions")
	hooks := make(map[string]string)
	for _, part := range parts[1:] {
		id, body, _ := strings.Cut(part, "\n")
		hooks[strings.TrimSpace(id)] = body
	}

	assert.Len(t, hooks, 3)
	for _, id := range []string{"gitleaks", "gitleaks-docker", "gitleaks-system"} {
		body, ok := hooks[id]
		require.True(t, ok, "missing hook %q", id)
		assert.Contains(t, body, "pass_filenames: false", "hook %q must not receive filenames", id)
	}
	assert.Contains(t, hooks["gitleaks"], "entry: gitleaks git --pre-commit --redact --staged --verbose")
	assert.Contains(t, hooks["gitleaks-docker"], "entry: zricethezav/gitleaks git --pre-commit --redact --staged --verbose")
	assert.Contains(t, hooks["gitleaks-system"], "entry: gitleaks git --pre-commit --redact --staged --verbose")
}

// TestPreCommitHelperScript runs the documented helper's regression suite when
// a Python 3 interpreter is available. The helper itself is Python, so this
// keeps its behavior covered by go test without adding a runtime dependency.
func TestPreCommitHelperScript(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	python := findPython(t)
	script, err := filepath.Abs("../scripts/test_pre_commit.py")
	require.NoError(t, err)
	cmd := exec.Command(python, script)
	cmd.Dir = filepath.Dir(script)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}

// TestPreCommitHookIntegration exercises the real pre-commit runner against the
// working-tree hook manifest. It is skipped when pre-commit is not installed;
// the metadata regression test above always runs.
func TestPreCommitHookIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	python := findPython(t)
	if err := exec.Command(python, "-m", "pre_commit", "--version").Run(); err != nil {
		t.Skip("pre-commit is not installed")
	}
	goBin, err := exec.LookPath("go")
	require.NoError(t, err)

	binary := filepath.Join(t.TempDir(), "gitleaks")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command(goBin, "build", "-o", binary, ".")
	build.Dir = ".."
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	script, err := filepath.Abs("../scripts/test_pre_commit_hooks.py")
	require.NoError(t, err)
	cmd := exec.Command(python, script)
	cmd.Dir = filepath.Dir(script)
	cmd.Env = append(os.Environ(), "GITLEAKS_BIN="+binary)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}

func findPython(t *testing.T) string {
	t.Helper()
	for _, name := range []string{"python3", "python"} {
		path, err := exec.LookPath(name)
		if err != nil {
			continue
		}
		if err := exec.Command(path, "--version").Run(); err != nil {
			continue
		}
		return path
	}
	t.Skip("python 3 is not available")
	return ""
}
