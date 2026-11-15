package cmd

import (
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.Flags().Bool("follow-symlinks", false, "scan named files that are symlinks to other files")
}

var filesCmd = &cobra.Command{
	Use:   "files [flags] <file>...",
	Short: "scan only the explicit files named on the command line",
	Args:  cobra.MinimumNArgs(1),
	Run:   runFiles,
}

func runFiles(cmd *cobra.Command, args []string) {
	for _, path := range args {
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			logging.Fatal().Msgf("explicit file path %q is a directory; use dir for directory scans", path)
		}
	}

	initConfig(".")
	initDiagnostics()
	cfg := Config(cmd)
	start := time.Now()
	detector := Detector(cmd, cfg, ".")
	setInformativeSourceType(detector, "file")

	followSymlinks, err := cmd.Flags().GetBool("follow-symlinks")
	if err != nil {
		logging.Fatal().Err(err).Send()
	}
	detector.FollowSymlinks = followSymlinks

	exitCode, err := cmd.Flags().GetInt("exit-code")
	if err != nil {
		logging.Fatal().Err(err).Msg("could not get exit code")
	}

	findings, err := detector.DetectSource(
		cmd.Context(),
		&sources.ExplicitFiles{
			Paths:           args,
			Config:          &cfg,
			FollowSymlinks:  detector.FollowSymlinks,
			MaxFileSize:     detector.MaxTargetMegaBytes * 1_000_000,
			MaxArchiveDepth: detector.MaxArchiveDepth,
			Informative:     detector.Informative,
			Sema:            detector.Sema,
		},
	)
	if err != nil {
		logging.Error().Err(err).Msg("failed scan explicit files")
	}
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].StartLine != findings[j].StartLine {
			return findings[i].StartLine < findings[j].StartLine
		}
		if findings[i].StartColumn != findings[j].StartColumn {
			return findings[i].StartColumn < findings[j].StartColumn
		}
		return findings[i].RuleID < findings[j].RuleID
	})
	findingSummaryAndExit(detector, findings, exitCode, start, err)
}
