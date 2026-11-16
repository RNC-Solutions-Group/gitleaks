package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/version"
)

const banner = `
    ○
    │╲
    │ ○
    ○ ░
    ░    gitleaks

`

const configDescription = `config file path
order of precedence:
1. --config/-c
2. env var GITLEAKS_CONFIG
3. env var GITLEAKS_CONFIG_TOML with the file content
4. (target path)/.gitleaks.toml
If none of the four options are used, then gitleaks will use the default config`

var (
	rootCmd = &cobra.Command{
		Use:     "gitleaks",
		Short:   "Gitleaks scans code, past or present, for secrets",
		Version: version.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set the timeout for all the commands
			if timeout, err := cmd.Flags().GetInt("timeout"); err != nil {
				return err
			} else if timeout > 0 {
				ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeout)*time.Second)
				cmd.SetContext(ctx)
				cobra.OnFinalize(cancel)
			}
			return nil
		},
	}

	// diagnostics manager is global to ensure it can be started before a scan begins
	// and stopped after a scan completes
	diagnosticsManager *DiagnosticsManager
)

const (
	BYTE     = 1.0
	KILOBYTE = BYTE * 1000
	MEGABYTE = KILOBYTE * 1000
	GIGABYTE = MEGABYTE * 1000
)

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringP("config", "c", "", configDescription)
	rootCmd.PersistentFlags().Bool("no-local-config", false, "ignore an automatically discovered (source)/.gitleaks.toml; explicit --config and environment configuration still apply")
	rootCmd.PersistentFlags().Int("exit-code", 1, "exit code when leaks have been encountered")
	rootCmd.PersistentFlags().StringP("report-path", "r", "", "report file (use \"-\" for stdout)")
	rootCmd.PersistentFlags().StringP("report-format", "f", "", "output format (json, csv, junit, sarif, template)")
	rootCmd.PersistentFlags().StringP("report-template", "", "", "template file used to generate the report (implies --report-format=template)")
	rootCmd.PersistentFlags().String("sarif-level", "", "optional SARIF result level (none, note, warning, error); empty preserves historical SARIF")
	rootCmd.PersistentFlags().Bool("report-unused-ignore", false, "after a completed scan, report loaded .gitleaksignore entries that did not suppress a finding")
	rootCmd.PersistentFlags().Bool("strict-config", false, "reject unknown configuration fields and incorrect nesting; default loading is unchanged")
	rootCmd.PersistentFlags().String("grouped-report", "", "write a separate opt-in keyed grouped report; requires --grouping-key-file")
	rootCmd.PersistentFlags().String("grouping-key-file", "", "file containing the HMAC key for --grouped-report (never printed)")
	rootCmd.PersistentFlags().Bool("overwrite", false, "replace an existing report file (overrides configured report policy)")
	rootCmd.PersistentFlags().Bool("aggregate", false, "write a numbered report instead of replacing an existing one (overrides configured report policy)")
	rootCmd.PersistentFlags().StringP("baseline-path", "b", "", "path to baseline with issues that can be ignored")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (trace, debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "show verbose output from scan")
	rootCmd.PersistentFlags().BoolP("no-color", "", false, "turn off color for verbose output")
	rootCmd.PersistentFlags().Bool("file-links", false, "print opt-in local file URI links in console findings; machine reports are unchanged")
	rootCmd.PersistentFlags().Int("max-target-megabytes", 0, "files larger than this will be skipped")
	rootCmd.PersistentFlags().BoolP("ignore-gitleaks-allow", "", false, "ignore gitleaks:allow comments")
	rootCmd.PersistentFlags().Uint("redact", 0, "redact secrets from logs and stdout. To redact only parts of the secret just apply a percent value from 0..100. For example --redact=20 (default 100%)")
	rootCmd.Flag("redact").NoOptDefVal = "100"
	rootCmd.PersistentFlags().Bool("no-banner", false, "suppress banner")
	rootCmd.PersistentFlags().StringSlice("enable-rule", []string{}, "only enable specific rules by id")
	rootCmd.PersistentFlags().StringP("gitleaks-ignore-path", "i", ".", "path to .gitleaksignore file or folder containing one")
	rootCmd.PersistentFlags().Int("max-decode-depth", 5, "allow recursive decoding up to this depth")
	rootCmd.PersistentFlags().Int("max-archive-depth", 0, "allow scanning into nested archives up to this depth (default \"0\", no archive traversal is done)")
	rootCmd.PersistentFlags().Int("timeout", 0, "set a timeout for gitleaks commands in seconds (default \"0\", no timeout is set)")

	// Add diagnostics flags
	rootCmd.PersistentFlags().String("diagnostics", "", "enable diagnostics (http OR comma-separated list: cpu,mem,trace). cpu=CPU prof, mem=memory prof, trace=exec tracing, http=serve via net/http/pprof")
	rootCmd.PersistentFlags().String("diagnostics-dir", "", "directory to store diagnostics output files when not using http mode (defaults to current directory)")

	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		logging.Fatal().Msgf("err binding config %s", err.Error())
	}
}

var logLevel = zerolog.InfoLevel

func initLog() {
	ll, err := rootCmd.Flags().GetString("log-level")
	if err != nil {
		logging.Fatal().Msg(err.Error())
	}

	switch strings.ToLower(ll) {
	case "trace":
		logLevel = zerolog.TraceLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "err", "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	default:
		logging.Warn().Msgf("unknown log level: %s", ll)
	}
	logging.Logger = logging.Logger.Level(logLevel)
}

func initConfig(source string) {
	hideBanner, err := rootCmd.Flags().GetBool("no-banner")
	viper.SetConfigType("toml")

	if err != nil {
		logging.Fatal().Msg(err.Error())
	}
	if !hideBanner {
		_, _ = fmt.Fprint(os.Stderr, banner)
	}

	logging.Debug().Msgf("using %s regex engine", regexp.Version)

	cfgPath, err := rootCmd.Flags().GetString("config")
	if err != nil {
		logging.Fatal().Msg(err.Error())
	}
	noLocalConfig, flagErr := rootCmd.Flags().GetBool("no-local-config")
	if flagErr != nil {
		logging.Fatal().Msg(flagErr.Error())
	}
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
		logging.Debug().Msgf("using gitleaks config %s from `--config`", cfgPath)
	} else if os.Getenv("GITLEAKS_CONFIG") != "" {
		envPath := os.Getenv("GITLEAKS_CONFIG")
		viper.SetConfigFile(envPath)
		logging.Debug().Msgf("using gitleaks config from GITLEAKS_CONFIG env var: %s", envPath)
	} else if os.Getenv("GITLEAKS_CONFIG_TOML") != "" {
		configContent := []byte(os.Getenv("GITLEAKS_CONFIG_TOML"))
		if err := viper.ReadConfig(bytes.NewBuffer(configContent)); err != nil {
			logging.Fatal().Err(err).Msg("unable to load gitleaks config from GITLEAKS_CONFIG_TOML env var")
		}
		logging.Debug().Msg("using gitleaks config from GITLEAKS_CONFIG_TOML env var content")
		return
	} else if noLocalConfig {
		logging.Debug().Msg("ignoring automatically discovered source-local .gitleaks.toml due to --no-local-config")
		if err = viper.ReadConfig(strings.NewReader(config.DefaultConfig)); err != nil {
			logging.Fatal().Msgf("err reading default config toml %s", err.Error())
		}
		return
	} else {
		fileInfo, err := os.Stat(source)
		if err != nil {
			logging.Fatal().Msg(err.Error())
		}

		if !fileInfo.IsDir() {
			logging.Debug().Msgf("unable to load gitleaks config from %s since --source=%s is a file, using default config",
				filepath.Join(source, ".gitleaks.toml"), source)
			if err = viper.ReadConfig(strings.NewReader(config.DefaultConfig)); err != nil {
				logging.Fatal().Msgf("err reading toml %s", err.Error())
			}
			return
		}

		if _, err := os.Stat(filepath.Join(source, ".gitleaks.toml")); os.IsNotExist(err) {
			logging.Debug().Msgf("no gitleaks config found in path %s, using default gitleaks config", filepath.Join(source, ".gitleaks.toml"))

			if err = viper.ReadConfig(strings.NewReader(config.DefaultConfig)); err != nil {
				logging.Fatal().Msgf("err reading default config toml %s", err.Error())
			}
			return
		} else {
			logging.Debug().Msgf("using existing gitleaks config %s from `(--source)/.gitleaks.toml`", filepath.Join(source, ".gitleaks.toml"))
		}

		viper.AddConfigPath(source)
		viper.SetConfigName(".gitleaks")
	}
	if err := viper.ReadInConfig(); err != nil {
		logging.Fatal().Msgf("unable to load gitleaks config, err: %s", err)
	}
}

func initDiagnostics() {
	// Initialize diagnostics manager
	diagnosticsFlag, err := rootCmd.PersistentFlags().GetString("diagnostics")
	if err != nil {
		logging.Fatal().Err(err).Msg("Error getting diagnostics flag")
	}

	diagnosticsDir, err := rootCmd.PersistentFlags().GetString("diagnostics-dir")
	if err != nil {
		logging.Fatal().Err(err).Msg("Error getting diagnostics-dir flag")
	}

	var diagErr error
	diagnosticsManager, diagErr = NewDiagnosticsManager(diagnosticsFlag, diagnosticsDir)
	if diagErr != nil {
		logging.Fatal().Err(diagErr).Msg("Error initializing diagnostics")
	}

	if diagnosticsManager.Enabled {
		logging.Info().Msg("Starting diagnostics...")
		if diagErr := diagnosticsManager.StartDiagnostics(); diagErr != nil {
			logging.Fatal().Err(diagErr).Msg("Failed to start diagnostics")
		}
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), "unknown flag") {
			// exit code 126: Command invoked cannot execute
			os.Exit(126)
		}
		logging.Fatal().Msg(err.Error())
	}
}

func Config(cmd *cobra.Command) config.Config {
	cfg, err := loadConfig(cmd)
	if err != nil {
		logging.Fatal().Err(err).Msg("Failed to load config")
	}
	return cfg
}

func loadConfig(cmd *cobra.Command) (config.Config, error) {
	var vc config.ViperConfig
	if err := viper.Unmarshal(&vc); err != nil {
		return config.Config{}, err
	}
	vc.ConfigPath = viper.ConfigFileUsed()

	cfg, err := vc.Translate()
	if err != nil {
		return config.Config{}, err
	}
	flagPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return config.Config{}, err
	}
	cfg.Path = flagPath
	return cfg, nil
}

func Detector(cmd *cobra.Command, cfg config.Config, source string) *detect.Detector {
	detector, err := newDetector(cmd, cfg, source)
	if err != nil {
		logging.Fatal().Err(err).Send()
	}
	return detector
}

// newDetector constructs the detector and report plan without terminating the
// process. Detector keeps the CLI boundary fatal; tests and future callers can
// exercise the complete control plane through this function and receive the
// error instead.
func newDetector(cmd *cobra.Command, cfg config.Config, source string) (*detect.Detector, error) {
	if err := validateStrictConfiguration(cmd, cfg); err != nil {
		return nil, err
	}

	opts, err := loadDetectorOptions(cmd)
	if err != nil {
		return nil, err
	}

	detector := detect.NewDetectorContext(cmd.Context(), cfg)
	if cfg.InformativeReport {
		detector.Informative = informative.New("directory")
	}
	opts.apply(detector)

	detector.Config.Path, err = cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	// If config path is not set, use the {source}/.gitleaks.toml path. The file
	// may not exist; this only records the source-local path for reporting.
	if detector.Config.Path == "" {
		detector.Config.Path = filepath.Join(source, ".gitleaks.toml")
	}

	gitleaksIgnorePath, err := cmd.Flags().GetString("gitleaks-ignore-path")
	if err != nil {
		return nil, fmt.Errorf("could not get .gitleaksignore path: %w", err)
	}
	if err := loadGitleaksIgnorePath(detector, gitleaksIgnorePath, cmd.Flags().Changed("gitleaks-ignore-path")); err != nil {
		return nil, fmt.Errorf("could not call AddGitleaksIgnore: %w", err)
	}
	if err := loadGitleaksIgnoreFile(detector, filepath.Join(source, ".gitleaksignore")); err != nil {
		return nil, fmt.Errorf("could not call AddGitleaksIgnore: %w", err)
	}

	baselinePath, err := cmd.Flags().GetString("baseline-path")
	if err != nil {
		return nil, err
	}
	if baselinePath != "" {
		if err := detector.AddBaseline(baselinePath, source); err != nil {
			logging.Error().Msgf("Could not load baseline. The path must point of a gitleaks report generated using the default format: %s", err)
		}
	}

	rules, err := cmd.Flags().GetStringSlice("enable-rule")
	if err != nil {
		return nil, err
	}
	if len(rules) > 0 {
		logging.Info().Msg("Overriding enabled rules: " + strings.Join(rules, ", "))
		ruleOverride := make(map[string]config.Rule)
		for _, ruleName := range rules {
			rule, ok := cfg.Rules[ruleName]
			if !ok {
				return nil, fmt.Errorf("requested rule %s not found in rules", ruleName)
			}
			ruleOverride[ruleName] = rule
		}
		detector.Config.Rules = ruleOverride
	}

	plan, err := buildReportPlan(cmd, cfg, source, baselinePath, gitleaksIgnorePath, detector.Informative)
	if err != nil {
		return nil, err
	}
	plan.apply(detector)
	return detector, nil
}

func rejectGroupingKeyInSource(keyPath, source string) error {
	if keyPath == "" || source == "" || source == "-" {
		return nil
	}
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil
	}
	if !sourceInfo.IsDir() {
		if samePath(keyPath, source) {
			return fmt.Errorf("grouping key must not be the scan target")
		}
		return nil
	}
	root, err := canonicalDirectory(source)
	if err != nil {
		return fmt.Errorf("could not resolve scan target for grouping-key check: %w", err)
	}
	key, err := canonicalDirectoryPath(keyPath)
	if err != nil {
		return fmt.Errorf("could not resolve grouping-key path: %w", err)
	}
	rel, err := filepath.Rel(root, key)
	if err != nil {
		return fmt.Errorf("could not compare grouping-key path with scan target: %w", err)
	}
	if rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))) {
		return fmt.Errorf("grouping key file %q is inside the scan target; store it outside the scanned tree", keyPath)
	}
	return nil
}

func canonicalDirectory(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func canonicalDirectoryPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(abs))
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(parent, filepath.Base(abs))), nil
}

func loadGroupingKey(path string) ([]byte, error) {
	before, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("could not inspect grouping key file: %w", err)
	}
	if before.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("grouping key file must not be a symbolic link")
	}
	if !before.Mode().IsRegular() {
		return nil, fmt.Errorf("grouping key file must be a regular file")
	}
	if before.Size() < 32 {
		return nil, fmt.Errorf("grouping key must contain at least 32 bytes")
	}
	if before.Size() > 1<<20 {
		return nil, fmt.Errorf("grouping key file is unreasonably large")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open grouping key file: %w", err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not inspect opened grouping key file: %w", err)
	}
	if !os.SameFile(before, opened) {
		return nil, fmt.Errorf("grouping key file changed while it was being opened")
	}
	if !opened.Mode().IsRegular() {
		return nil, fmt.Errorf("grouping key file must be a regular file")
	}
	if opened.Size() < 32 {
		return nil, fmt.Errorf("grouping key must contain at least 32 bytes")
	}
	if runtime.GOOS != "windows" && opened.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("grouping key file must not be group/world accessible")
	}
	key, err := io.ReadAll(io.LimitReader(file, (1<<20)+1))
	if err != nil {
		return nil, fmt.Errorf("could not read grouping key file: %w", err)
	}
	if len(key) > 1<<20 {
		return nil, fmt.Errorf("grouping key file is unreasonably large")
	}
	if len(key) < 32 {
		return nil, fmt.Errorf("grouping key must contain at least 32 bytes")
	}
	return key, nil
}

func inlineConfigurationIsEffective(cmd *cobra.Command) bool {
	if cmd.Flags().Changed("config") {
		return false
	}
	if configPath := os.Getenv("GITLEAKS_CONFIG"); configPath != "" {
		return false
	}
	return true
}

func validateStrictConfiguration(cmd *cobra.Command, cfg config.Config) error {
	strict, err := cmd.Flags().GetBool("strict-config")
	if err != nil {
		return fmt.Errorf("could not read --strict-config: %w", err)
	}
	if !strict {
		return nil
	}
	files := append([]string(nil), cfg.ConfigFiles...)
	if len(files) == 0 && cfg.Path != "" {
		files = append(files, cfg.Path)
	}

	// An inline configuration remains an effective source even when it extends
	// a file, which puts that extension in ConfigFiles. Validate the inline
	// root as well as every file selected by precedence/extension resolution.
	if inline := os.Getenv("GITLEAKS_CONFIG_TOML"); inline != "" && inlineConfigurationIsEffective(cmd) {
		if err := config.ValidateStrictTOML([]byte(inline)); err != nil {
			return fmt.Errorf("strict inline configuration validation failed: %w", err)
		}
	}
	for _, path := range files {
		if err := config.ValidateStrictFile(path); err != nil {
			return fmt.Errorf("strict configuration validation failed: %w", err)
		}
	}
	return nil
}

func setInformativeSourceType(detector *detect.Detector, sourceType string) {
	if detector == nil || detector.Informative == nil {
		return
	}
	detector.Informative.SetSourceType(sourceType)
}

func bytesConvert(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "bytes"
	case bytes == 0:
		return "0"
	}

	stringValue := strings.TrimSuffix(
		fmt.Sprintf("%.2f", value), ".00",
	)

	return fmt.Sprintf("%s %s", stringValue, unit)
}

func findingSummaryAndExit(detector *detect.Detector, findings []report.Finding, exitCode int, start time.Time, err error) {
	if diagnosticsManager.Enabled {
		logging.Debug().Msg("Finalizing diagnostics...")
		diagnosticsManager.StopDiagnostics()
	}

	totalBytes := detector.TotalBytes.Load()
	bytesMsg := fmt.Sprintf("scanned ~%d bytes (%s)", totalBytes, bytesConvert(totalBytes))
	if err == nil {
		logging.Info().Msgf("%s in %s", bytesMsg, FormatDuration(time.Since(start)))
		if len(findings) != 0 {
			logging.Warn().Msgf("leaks found: %d", len(findings))
		} else {
			logging.Info().Msg("no leaks found")
		}
	} else {
		logging.Warn().Msg(bytesMsg)
		logging.Warn().Msgf("partial scan completed in %s", FormatDuration(time.Since(start)))
		if len(findings) != 0 {
			logging.Warn().Msgf("%d leaks found in partial scan", len(findings))
		} else {
			logging.Warn().Msg("no leaks found in partial scan")
		}
	}

	if err == nil && detector.ReportUnusedIgnore {
		unused := detector.UnusedGitleaksIgnore()
		logging.Warn().Msgf("%d loaded .gitleaksignore entries did not suppress a finding in this completed scan", len(unused))
		for _, entry := range unused {
			logging.Warn().Str("fingerprint", entry).Msg("ignore entry not triggered")
		}
		if len(unused) > 0 {
			logging.Warn().Msg("not-triggered is not proof of obsolescence: the relevant path, commit or revision may not have been examined")
		}
	}

	// Write the report only for a scan that completed without a propagated
	// error. The destination is never touched during validation, and an
	// existing report is preserved if the scan or the write fails.
	if detector.Reporter != nil {
		if err != nil {
			logging.Warn().Msg("report not written because the scan was incomplete")
		} else {
			if detector.ReportPath == report.StdoutReportPath {
				if reportErr := detector.Reporter.Write(os.Stdout, findings); reportErr != nil {
					logging.Fatal().Err(reportErr).Msg("failed to write report")
				}
				detector.PublishedReportPath = "stdout"
			} else {
				if detector.ReportWriter == nil {
					logging.Fatal().Msg("failed to write report: report writer is not configured")
				}
				outputPath, reportErr := detector.ReportWriter(findings)
				if reportErr != nil {
					logging.Fatal().Err(reportErr).Msg("failed to write report")
				}
				logging.Info().Str("path", outputPath).Msg("report written")
			}
		}
	}

	if err == nil && detector.GroupedReportWriter != nil {
		groupedPath, groupedErr := detector.GroupedReportWriter(findings)
		if groupedErr != nil {
			if detector.PublishedReportPath != "" {
				logging.Fatal().Err(groupedErr).Msgf("failed to write grouped report after primary report %q was published; multi-file publication is not crash-atomic", detector.PublishedReportPath)
			}
			logging.Fatal().Err(groupedErr).Msg("failed to write grouped report")
		}
		logging.Info().Str("path", groupedPath).Msg("grouped report written")
	}

	// Print the informative console summary after the primary report has been
	// handled. When a companion file is configured, the summary lives in that
	// file instead. Console output always uses stderr so a stdout report remains
	// machine-readable.
	if err == nil && detector.Informative != nil && detector.InformativeReportPath == "" {
		primaryReport := ""
		if detector.ReportPath != "" {
			primaryReport = detector.ReportPath
		}
		document := detector.Informative.Document(primaryReport, true)
		if consoleErr := informative.WriteConsole(os.Stderr, document, informative.DefaultConsoleDetails); consoleErr != nil {
			logging.Warn().Err(consoleErr).Msg("could not write informative console summary")
		}
	}

	if err != nil {
		os.Exit(1)
	}

	if len(findings) != 0 {
		os.Exit(exitCode)
	}
}

func loadGitleaksIgnorePath(detector *detect.Detector, path string, required bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) && !required {
			return nil
		}
		return fmt.Errorf("could not inspect .gitleaksignore path %q: %w", path, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		info, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not resolve .gitleaksignore path %q: %w", path, err)
		}
	}

	if info.IsDir() {
		return loadGitleaksIgnoreFile(detector, filepath.Join(path, ".gitleaksignore"))
	}
	return detector.AddGitleaksIgnore(path)
}

func loadGitleaksIgnoreFile(detector *detect.Detector, path string) error {
	// An auto-discovered ignore file is a sibling of the scan target. When the
	// target is a file, filepath.Join(target, ".gitleaksignore") is not a
	// meaningful path; on Unix Lstat reports ENOTDIR rather than ENOENT.
	if parent, parentErr := os.Stat(filepath.Dir(path)); parentErr == nil && !parent.IsDir() {
		return nil
	}

	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("could not inspect .gitleaksignore file %q: %w", path, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		info, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not resolve .gitleaksignore file %q: %w", path, err)
		}
	}

	// A directory named .gitleaksignore is not an ignore file. Preserve the
	// historical absence behavior for auto-discovered source files.
	if info.IsDir() {
		return nil
	}
	return detector.AddGitleaksIgnore(path)
}

func FormatDuration(d time.Duration) string {
	scale := 100 * time.Second
	// look for the max scale that is smaller than d
	for scale > d {
		scale = scale / 10
	}
	return d.Round(scale / 100).String()
}

func mustGetBoolFlag(cmd *cobra.Command, name string) bool {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		logging.Fatal().Err(err).Msgf("could not get flag: %s", name)
	}
	return value
}

func mustGetIntFlag(cmd *cobra.Command, name string) int {
	value, err := cmd.Flags().GetInt(name)
	if err != nil {
		logging.Fatal().Err(err).Msgf("could not get flag: %s", name)
	}
	return value
}

func mustGetStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		logging.Fatal().Err(err).Msgf("could not get flag: %s", name)
	}
	return value
}
