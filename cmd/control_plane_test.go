package cmd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func newDetectorTestCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "detector-test"}
	flags := cmd.Flags()
	flags.String("config", "", "")
	flags.Bool("strict-config", false, "")
	flags.Int("max-decode-depth", 5, "")
	flags.Int("max-archive-depth", 0, "")
	flags.Bool("no-color", false, "")
	flags.Bool("file-links", false, "")
	flags.Bool("verbose", false, "")
	flags.Uint("redact", 0, "")
	flags.Int("max-target-megabytes", 0, "")
	flags.Bool("report-unused-ignore", false, "")
	flags.Bool("ignore-gitleaks-allow", false, "")
	flags.String("gitleaks-ignore-path", "", "")
	flags.String("baseline-path", "", "")
	flags.StringSlice("enable-rule", nil, "")
	flags.String("report-path", "", "")
	flags.String("report-format", "", "")
	flags.String("report-template", "", "")
	flags.String("sarif-level", "", "")
	flags.String("grouped-report", "", "")
	flags.String("grouping-key-file", "", "")
	flags.Bool("overwrite", false, "")
	flags.Bool("aggregate", false, "")
	return cmd
}

func TestNewDetectorReturnsErrorsWithoutExiting(t *testing.T) {
	source := t.TempDir()
	cmd := newDetectorTestCommand()
	require.NoError(t, cmd.Flags().Set("gitleaks-ignore-path", source))
	cfg := config.Config{
		AutoOverwrite: true,
		Rules: map[string]config.Rule{
			"test-rule": {
				RuleID: "test-rule",
				Regex:  regexp.MustCompile("synthetic"),
			},
		},
		OrderedRules: []string{"test-rule"},
	}

	detector, err := newDetector(cmd, cfg, source)
	require.NoError(t, err)
	require.NotNil(t, detector)
	assert.Equal(t, 5, detector.MaxDecodeDepth)
	assert.Equal(t, filepath.Join(source, ".gitleaks.toml"), detector.Config.Path)

	require.NoError(t, cmd.Flags().Set("enable-rule", "missing-rule"))
	_, err = newDetector(cmd, cfg, source)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requested rule missing-rule not found")
}
