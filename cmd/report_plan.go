package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/report"
)

// reportPlan owns every report destination and writer for one scan. It is
// built and validated before scanning, then applied to the detector. The plan
// never opens a destination during construction; publication remains owned by
// reportOutput and happens only after a complete scan.
type reportPlan struct {
	reportPath    string
	reporter      report.Reporter
	primary       *reportOutput
	companionPath string

	groupedPath string
	grouped     *reportOutput
	groupingKey []byte
}

func reportConfigFiles(cfg config.Config) []string {
	files := append([]string(nil), cfg.ConfigFiles...)
	if len(files) == 0 && cfg.Path != "" {
		files = append(files, cfg.Path)
	}
	return files
}

// buildReportPlan resolves report flags and configuration policy, validates
// the complete destination set, and loads the optional grouping key. It returns
// errors instead of terminating the process so the control plane can be tested
// without a subprocess.
func buildReportPlan(cmd *cobra.Command, cfg config.Config, source string, baselinePath, gitleaksIgnorePath string, collector *informative.Collector) (*reportPlan, error) {
	mode, err := resolveReportMode(cfg, cmd)
	if err != nil {
		return nil, fmt.Errorf("invalid report options: %w", err)
	}
	plan := &reportPlan{}

	reportPath, err := cmd.Flags().GetString("report-path")
	if err != nil {
		return nil, err
	}
	reportFormat, err := cmd.Flags().GetString("report-format")
	if err != nil {
		return nil, err
	}
	reportFormat = strings.TrimSpace(strings.ToLower(reportFormat))
	reportTemplate, err := cmd.Flags().GetString("report-template")
	if err != nil {
		return nil, err
	}
	sarifLevel, err := cmd.Flags().GetString("sarif-level")
	if err != nil {
		return nil, err
	}
	groupedPath, err := cmd.Flags().GetString("grouped-report")
	if err != nil {
		return nil, err
	}
	groupingKeyPath, err := cmd.Flags().GetString("grouping-key-file")
	if err != nil {
		return nil, err
	}

	configFiles := reportConfigFiles(cfg)

	if groupedPath != "" || groupingKeyPath != "" {
		if groupedPath == "" || groupingKeyPath == "" {
			return nil, fmt.Errorf("--grouped-report and --grouping-key-file must be used together")
		}
		if groupedPath == report.StdoutReportPath {
			return nil, fmt.Errorf("--grouped-report requires a file path, not stdout")
		}
		if err := rejectGroupingKeyInSource(groupingKeyPath, source); err != nil {
			return nil, fmt.Errorf("invalid grouping key location: %w", err)
		}
		plan.groupingKey, err = loadGroupingKey(groupingKeyPath)
		if err != nil {
			return nil, fmt.Errorf("invalid grouping key: %w", err)
		}
		plan.groupedPath = groupedPath
	}

	if reportPath != "" {
		reporter, err := buildReporter(reportPath, reportFormat, reportTemplate, sarifLevel, cfg)
		if err != nil {
			return nil, err
		}
		plan.reportPath = reportPath
		plan.reporter = reporter

		if reportPath != report.StdoutReportPath {
			output, err := newReportOutput(reportPath, mode, source, configFiles, baselinePath, gitleaksIgnorePath, reportTemplate)
			if err != nil {
				return nil, fmt.Errorf("invalid report destination: %w", err)
			}
			plan.primary = output
			if groupingKeyPath != "" {
				if err := output.protectInput(groupingKeyPath); err != nil {
					return nil, fmt.Errorf("grouping key collides with the primary report destination: %w", err)
				}
			}
			if collector != nil && !cfg.ConsoleOnly {
				if err := output.enableInformative(func(primaryPath string) ([]byte, error) {
					return informative.Marshal(collector.Document(primaryPath, true))
				}); err != nil {
					return nil, fmt.Errorf("invalid informative report destination: %w", err)
				}
				// Prevent a pre-existing companion from being scanned as source
				// content. The actual aggregate filename is allocated later; it
				// cannot exist before publication unless it is an older report.
				plan.companionPath = output.companion.path
			}
		}
	}

	if groupingKeyPath != "" {
		output, err := newReportOutput(groupedPath, mode, source, configFiles, baselinePath, gitleaksIgnorePath, reportTemplate)
		if err != nil {
			return nil, fmt.Errorf("invalid grouped report destination: %w", err)
		}
		if err := output.protectInput(groupingKeyPath); err != nil {
			return nil, fmt.Errorf("grouping key collides with the grouped report destination: %w", err)
		}
		if reportPath != "" && reportPath != report.StdoutReportPath {
			if err := output.protectInput(reportPath); err != nil {
				return nil, fmt.Errorf("primary report collides with the grouped report destination: %w", err)
			}
		}
		if plan.companionPath != "" {
			if err := output.protectInput(plan.companionPath); err != nil {
				return nil, fmt.Errorf("informative companion collides with the grouped report destination: %w", err)
			}
		}
		if err := validateReportDestinationSet(reportPath, plan.companionPath, groupedPath, mode); err != nil {
			return nil, fmt.Errorf("invalid report destination set: %w", err)
		}
		if plan.primary != nil {
			// Make the grouped destination a protected peer of the primary
			// output. The selected primary/companion path is checked against it
			// after reservation and immediately before publication.
			if err := plan.primary.protectInput(groupedPath); err != nil {
				return nil, fmt.Errorf("grouped report collides with the primary report destination: %w", err)
			}
		}
		plan.grouped = output
	}

	return plan, nil
}

func buildReporter(reportPath, reportFormat, reportTemplate, sarifLevel string, cfg config.Config) (report.Reporter, error) {
	if reportFormat == "" {
		ext := strings.ToLower(filepath.Ext(reportPath))
		switch ext {
		case ".csv":
			reportFormat = "csv"
		case ".json":
			reportFormat = "json"
		case ".sarif":
			reportFormat = "sarif"
		default:
			return nil, fmt.Errorf("Unknown report format: %s", reportFormat)
		}
		logging.Debug().Msgf("No report format specified, inferred %q from %q", reportFormat, ext)
	}

	var reporter report.Reporter
	switch reportFormat {
	case "csv":
		reporter = &report.CsvReporter{}
	case "json":
		reporter = &report.JsonReporter{}
	case "junit":
		reporter = &report.JunitReporter{}
	case "sarif":
		switch sarifLevel {
		case "", "none", "note", "warning", "error":
		default:
			return nil, fmt.Errorf("invalid --sarif-level %q (expected none, note, warning, error)", sarifLevel)
		}
		reporter = &report.SarifReporter{
			OrderedRules: cfg.GetOrderedRules(),
			Level:        sarifLevel,
		}
	case "template":
		templateReporter, err := report.NewTemplateReporter(reportTemplate)
		if err != nil {
			return nil, fmt.Errorf("Invalid report template: %w", err)
		}
		reporter = templateReporter
	default:
		return nil, fmt.Errorf("unknown report format %s", reportFormat)
	}

	if reportTemplate != "" && reportFormat != "template" {
		return nil, fmt.Errorf("Report format must be 'template' if --report-template is specified")
	}
	if reportFormat != "sarif" && sarifLevel != "" {
		return nil, fmt.Errorf("--sarif-level requires --report-format=sarif")
	}
	return reporter, nil
}

func (p *reportPlan) apply(d *detect.Detector) {
	if p.reportPath != "" {
		d.ReportPath = p.reportPath
		d.Reporter = p.reporter
		if p.primary != nil {
			d.InformativeReportPath = p.companionPath
			output := p.primary
			d.ReportWriter = func(findings []report.Finding) (string, error) {
				path, writeErr := output.write(p.reporter, findings)
				if writeErr == nil {
					d.PublishedReportPath = path
				}
				return path, writeErr
			}
		}
	}

	if p.grouped != nil {
		d.GroupedReportPath = p.groupedPath
		output := p.grouped
		key := p.groupingKey
		d.GroupedReportWriter = func(findings []report.Finding) (string, error) {
			return output.write(&report.GroupedReporter{Key: key}, findings)
		}
	}
}
