package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/logging"
)

// detectorOptions is the set of CLI values that tune a detector after it has
// been constructed from configuration. Keeping extraction separate from
// construction makes flag failures ordinary errors and gives the control-plane
// tests a value to inspect without starting a scan.
type detectorOptions struct {
	maxDecodeDepth      int
	maxArchiveDepth     int
	noColor             bool
	fileLinks           bool
	verbose             bool
	redact              uint
	maxTargetMegaBytes  int
	reportUnusedIgnore  bool
	ignoreGitleaksAllow bool
}

func loadDetectorOptions(cmd *cobra.Command) (detectorOptions, error) {
	var (
		opts detectorOptions
		err  error
	)
	if opts.maxDecodeDepth, err = cmd.Flags().GetInt("max-decode-depth"); err != nil {
		return opts, err
	}
	if opts.maxArchiveDepth, err = cmd.Flags().GetInt("max-archive-depth"); err != nil {
		return opts, err
	}
	if opts.noColor, err = cmd.Flags().GetBool("no-color"); err != nil {
		return opts, err
	}
	if opts.fileLinks, err = cmd.Flags().GetBool("file-links"); err != nil {
		return opts, err
	}
	if opts.verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
		return opts, err
	}
	if opts.redact, err = cmd.Flags().GetUint("redact"); err != nil {
		return opts, err
	}
	if opts.maxTargetMegaBytes, err = cmd.Flags().GetInt("max-target-megabytes"); err != nil {
		return opts, err
	}
	if opts.reportUnusedIgnore, err = cmd.Flags().GetBool("report-unused-ignore"); err != nil {
		return opts, err
	}
	if opts.ignoreGitleaksAllow, err = cmd.Flags().GetBool("ignore-gitleaks-allow"); err != nil {
		return opts, err
	}
	return opts, nil
}

func (o detectorOptions) apply(d *detect.Detector) {
	d.MaxDecodeDepth = o.maxDecodeDepth
	d.MaxArchiveDepth = o.maxArchiveDepth
	d.NoColor = o.noColor
	d.FileLinks = o.fileLinks
	d.Verbose = o.verbose
	d.Redact = o.redact
	d.MaxTargetMegaBytes = o.maxTargetMegaBytes
	d.ReportUnusedIgnore = o.reportUnusedIgnore
	d.IgnoreGitleaksAllow = o.ignoreGitleaksAllow

	// Preserve the historical behavior: the console logger is rebuilt only
	// when color is disabled, and at the level selected before the scan.
	if d.NoColor {
		logging.Logger = log.Output(zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: d.NoColor,
		}).Level(logLevel)
	}
}
