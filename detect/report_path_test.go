package detect

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func TestDetectSourceSkipsReportPath(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.json")
	otherPath := filepath.Join(dir, "other.txt")
	content := "aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\"\n"
	require.NoError(t, os.WriteFile(reportPath, []byte(content), 0o600))
	require.NoError(t, os.WriteFile(otherPath, []byte(content), 0o600))

	cfg := config.Config{
		Rules: map[string]config.Rule{
			"test-rule": {
				RuleID:   "test-rule",
				Regex:    regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
				Keywords: []string{"akia"},
			},
		},
		Keywords: map[string]struct{}{"akia": {}},
	}
	detector := NewDetector(cfg)
	detector.ReportPath = reportPath

	findings, err := detector.DetectSource(context.Background(), &sources.Files{
		Config: &cfg,
		Path:   dir,
		Sema:   detector.Sema,
	})
	require.NoError(t, err)
	require.Len(t, findings, 1)
	require.Equal(t, filepath.ToSlash(otherPath), findings[0].File)
}
