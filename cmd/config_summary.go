package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/logging"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [flags] [path]",
	Short: "print a safe summary of the effective configuration",
	Long: `Print a safe summary of the configuration that a scan would use after
precedence, extension resolution, rule selection and report policy are applied.

The summary intentionally reports counts and paths, not rule regexes, allowlist
patterns or other configuration content that may contain sensitive material.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runConfigSummary,
}

func runConfigSummary(cmd *cobra.Command, args []string) {
	source := "."
	if len(args) == 1 {
		source = args[0]
	}
	initConfig(source)
	cfg := Config(cmd)

	fmt.Fprintln(cmd.OutOrStdout(), "RNC-Mod effective configuration")
	fmt.Fprintf(cmd.OutOrStdout(), "Target: %s\n", source)
	fmt.Fprintf(cmd.OutOrStdout(), "Configuration source: %s\n", configurationSource(cmd, cfg))
	fmt.Fprintln(cmd.OutOrStdout(), "Loaded files:")
	if len(cfg.ConfigFiles) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (embedded default)")
	} else {
		files := append([]string(nil), cfg.ConfigFiles...)
		sort.Strings(files)
		for _, path := range files {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", path)
		}
	}

	disabled := len(cfg.Extend.DisabledRules)
	fmt.Fprintf(cmd.OutOrStdout(), "Rules active: %d\n", len(cfg.Rules))
	fmt.Fprintf(cmd.OutOrStdout(), "Rules disabled by extension: %d\n", disabled)
	fmt.Fprintf(cmd.OutOrStdout(), "Ordered rules: %d\n", len(cfg.OrderedRules))
	fmt.Fprintf(cmd.OutOrStdout(), "Global allowlists: %d\n", len(cfg.Allowlists))
	ruleAllowlists := 0
	for _, rule := range cfg.Rules {
		ruleAllowlists += len(rule.Allowlists)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Rule allowlists: %d\n", ruleAllowlists)
	fmt.Fprintf(cmd.OutOrStdout(), "Minimum version: %s\n", displayValue(cfg.MinVersion))
	fmt.Fprintf(cmd.OutOrStdout(), "Title: %s\n", displayValue(cfg.Title))
	fmt.Fprintf(cmd.OutOrStdout(), "Report policy: auto-overwrite=%t aggregation=%t informative-report=%t console-only=%t\n",
		cfg.AutoOverwrite, cfg.Aggregation, cfg.InformativeReport, cfg.ConsoleOnly)

	for _, name := range []string{"max-target-megabytes", "max-decode-depth", "max-archive-depth", "timeout"} {
		value, err := cmd.Flags().GetInt(name)
		if err != nil {
			logging.Fatal().Err(err).Msgf("could not read --%s", name)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Limit %s: %d\n", name, value)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Sensitive rule/allowlist patterns are intentionally not printed.")
}

func configurationSource(cmd *cobra.Command, cfg config.Config) string {
	if cmd.Flags().Changed("config") {
		return "explicit --config"
	}
	if os.Getenv("GITLEAKS_CONFIG") != "" {
		return "GITLEAKS_CONFIG"
	}
	if os.Getenv("GITLEAKS_CONFIG_TOML") != "" {
		return "GITLEAKS_CONFIG_TOML"
	}
	if len(cfg.ConfigFiles) > 0 {
		return "target-local .gitleaks.toml"
	}
	return "embedded default"
}

func displayValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "(not set)"
	}
	return value
}
