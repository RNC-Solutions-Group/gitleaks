package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/report"
)

func newReportPlanTestCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "report-plan-test"}
	flags := cmd.Flags()
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

func TestBuildReportPlanPrimaryAndGrouped(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	require.NoError(t, os.Mkdir(source, 0o700))
	reportPath := filepath.Join(dir, "report.json")
	groupedPath := filepath.Join(dir, "grouped.json")
	keyPath := filepath.Join(dir, "grouping.key")
	require.NoError(t, os.WriteFile(keyPath, []byte(strings.Repeat("k", 32)), 0o600))

	cmd := newReportPlanTestCommand()
	require.NoError(t, cmd.Flags().Set("report-path", reportPath))
	require.NoError(t, cmd.Flags().Set("report-format", "json"))
	require.NoError(t, cmd.Flags().Set("grouped-report", groupedPath))
	require.NoError(t, cmd.Flags().Set("grouping-key-file", keyPath))

	cfg := config.Config{
		AutoOverwrite: true,
		Rules: map[string]config.Rule{
			"test-rule": {RuleID: "test-rule", Description: "test"},
		},
		OrderedRules: []string{"test-rule"},
	}
	plan, err := buildReportPlan(cmd, cfg, source, "", "", nil)
	require.NoError(t, err)
	require.NotNil(t, plan.primary)
	require.NotNil(t, plan.grouped)

	d := detect.NewDetector(cfg)
	plan.apply(d)
	assert.Equal(t, reportPath, d.ReportPath)
	assert.Equal(t, groupedPath, d.GroupedReportPath)
	require.NotNil(t, d.ReportWriter)
	require.NotNil(t, d.GroupedReportWriter)

	findings := []report.Finding{{
		RuleID: "test-rule",
		File:   "a.txt",
		Secret: "synthetic-secret",
		Match:  "synthetic-secret",
	}}
	published, err := d.ReportWriter(findings)
	require.NoError(t, err)
	assert.Equal(t, reportPath, published)
	body, err := os.ReadFile(reportPath)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"RuleID": "test-rule"`)

	grouped, err := d.GroupedReportWriter(findings)
	require.NoError(t, err)
	assert.Equal(t, groupedPath, grouped)
	_, err = os.Stat(groupedPath)
	require.NoError(t, err)
}

func TestBuildReporterValidation(t *testing.T) {
	cfg := config.Config{Rules: map[string]config.Rule{}}

	tests := []struct {
		name       string
		path       string
		format     string
		template   string
		sarifLevel string
		wantError  string
	}{
		{name: "inferred json", path: "report.json"},
		{name: "inferred csv", path: "report.csv"},
		{name: "inferred sarif", path: "report.sarif"},
		{name: "unknown extension", path: "report.txt", wantError: "Unknown report format"},
		{name: "unknown format", path: "report.json", format: "xml", wantError: "unknown report format"},
		{name: "bad sarif level", path: "report.sarif", format: "sarif", sarifLevel: "debug", wantError: "invalid --sarif-level"},
		{name: "sarif level without sarif", path: "report.json", format: "json", sarifLevel: "note", wantError: "--sarif-level requires"},
		{name: "template mismatch", path: "report.json", format: "json", template: "template.tmpl", wantError: "Report format must be 'template'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildReporter(tt.path, tt.format, tt.template, tt.sarifLevel, cfg)
			if tt.wantError == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}
