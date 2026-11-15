package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLISARIFLevelIsOptIn(t *testing.T) {
	source := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(source, "secret.txt"), []byte("ALPHA_12345678\n"), 0o600))
	cfg := redactionConfig(t, "")

	for _, test := range []struct {
		name      string
		args      []string
		wantLevel bool
	}{
		{name: "default omission", args: nil, wantLevel: false},
		{name: "explicit warning", args: []string{"--sarif-level", "warning"}, wantLevel: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			reportPath := filepath.Join(t.TempDir(), "report.sarif")
			args := []string{
				"dir", source,
				"--config", cfg,
				"--exit-code", "0",
				"--no-banner",
				"--log-level", "error",
				"--report-path", reportPath,
				"--report-format", "sarif",
			}
			args = append(args, test.args...)
			result := runReportCLI(t, args...)
			require.Equal(t, 0, result.exitCode, result.stderr)
			raw := string(mustReadFile(t, reportPath))
			if test.wantLevel {
				assert.Contains(t, raw, `"level": "warning"`)
			} else {
				assert.NotContains(t, raw, `"level"`)
			}
		})
	}
}
