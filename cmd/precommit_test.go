package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPreCommitHelperScript runs the documented helper's regression suite when
// a Python 3 interpreter is available. The helper itself is Python, so this
// keeps its behavior covered by go test without adding a runtime dependency.
func TestPreCommitHelperScript(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	var python string
	for _, name := range []string{"python3", "python"} {
		path, err := exec.LookPath(name)
		if err != nil {
			continue
		}
		if err := exec.Command(path, "--version").Run(); err != nil {
			continue
		}
		python = path
		break
	}
	if python == "" {
		t.Skip("python 3 is not available")
	}

	script, err := filepath.Abs("../scripts/test_pre_commit.py")
	require.NoError(t, err)
	cmd := exec.Command(python, script)
	cmd.Dir = filepath.Dir(script)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}
