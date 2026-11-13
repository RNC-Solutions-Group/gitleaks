// == WARNING ==
// These functions are used to generate GitLeak's default config.
// You are free to use these in your own project, HOWEVER, no API stability is guaranteed.

package utils

import (
	"strings"

	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/base"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/secrets"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/logging"
)

func Validate(rule config.Rule, truePositives []string, falsePositives []string) *config.Rule {
	r := &rule
	d := createSingleRuleDetector(r)
	for _, tp := range truePositives {
		if len(d.DetectString(tp)) < 1 {
			diagnostic, entropy, hasEntropy := diagnoseTruePositive(r, tp)
			event := logging.Fatal().
				Str("rule", r.RuleID).
				Str("diagnostic", diagnostic).
				Str("regex", ruleRegex(r))
			if hasEntropy {
				event = event.
					Float64("entropy", entropy).
					Float64("required-entropy", r.Entropy)
			}
			event.Msg("Failed to Validate. True positive was not detected by detector.")
		}
	}
	for _, fp := range falsePositives {
		findings := d.DetectString(fp)
		if len(findings) != 0 {
			logging.Fatal().
				Str("rule", r.RuleID).
				Str("regex", ruleRegex(r)).
				Msg("Failed to Validate. False positive was detected by detector.")
		}
	}
	return r
}

func ValidateWithPaths(rule config.Rule, truePositives map[string]string, falsePositives map[string]string) *config.Rule {
	r := &rule
	d := createSingleRuleDetector(r)
	for path, tp := range truePositives {
		f := detect.Fragment{Raw: tp, FilePath: path}
		if len(d.Detect(f)) != 1 {
			logging.Fatal().
				Str("rule", r.RuleID).
				Str("regex", ruleRegex(r)).
				Str("path", rulePath(r)).
				Msg("Failed to Validate. True positive was not detected by regex and/or path.")
		}
	}
	for path, fp := range falsePositives {
		f := detect.Fragment{Raw: fp, FilePath: path}
		if len(d.Detect(f)) != 0 {
			logging.Fatal().
				Str("rule", r.RuleID).
				Str("regex", ruleRegex(r)).
				Str("path", rulePath(r)).
				Msg("Failed to Validate. False positive was detected by regex and/or path.")
		}
	}
	return r
}

// diagnoseTruePositive classifies a detector miss without changing detector
// behavior. It evaluates the same regex/capture-group/entropy pipeline that
// detect.detectRule uses, then checks whether removing allowlists makes the
// sample detectable. The returned fields are safe for diagnostics and never
// include the full candidate value.
func diagnoseTruePositive(r *config.Rule, value string) (string, float64, bool) {
	if r.Regex == nil {
		return "rule has no regex", 0, false
	}

	matches := r.Regex.FindAllString(value, -1)
	if len(matches) == 0 {
		return "regex/structure mismatch", 0, false
	}

	var (
		maxEntropy   float64
		entropyPass  bool
		invalidGroup bool
	)
	for _, match := range matches {
		secret := strings.Trim(match, "\n")
		groups := r.Regex.FindStringSubmatch(secret)
		if len(groups) >= 2 {
			if r.SecretGroup > 0 {
				if len(groups) <= r.SecretGroup {
					invalidGroup = true
					continue
				}
				secret = groups[r.SecretGroup]
			} else {
				for _, group := range groups[1:] {
					if len(group) > 0 {
						secret = group
						break
					}
				}
			}
		}

		entropy := secrets.ShannonEntropy(secret)
		if entropy > maxEntropy {
			maxEntropy = entropy
		}
		if r.Entropy == 0 || entropy > r.Entropy {
			entropyPass = true
			break
		}
	}

	if invalidGroup && !entropyPass {
		return "invalid secret group", 0, false
	}
	if r.Entropy != 0 && !entropyPass {
		return "entropy rejection", maxEntropy, true
	}

	// The only difference between this detector and the canonical detector is
	// the absence of global and rule allowlists. If the value is detected here,
	// an allowlist (including an inline gitleaks:allow comment) suppressed it.
	if len(createSingleRuleDetectorWithoutAllowlists(r).DetectString(value)) > 0 {
		return "allowlist suppression", 0, false
	}
	return "detector validation failure", 0, false
}

func ruleRegex(r *config.Rule) string {
	if r.Regex == nil {
		return ""
	}
	return r.Regex.String()
}

func rulePath(r *config.Rule) string {
	if r.Path == nil {
		return ""
	}
	return r.Path.String()
}

func createSingleRuleDetector(r *config.Rule) *detect.Detector {
	return buildSingleRuleDetector(r, true)
}

func createSingleRuleDetectorWithoutAllowlists(r *config.Rule) *detect.Detector {
	ruleCopy := *r
	ruleCopy.Allowlists = nil
	return buildSingleRuleDetector(&ruleCopy, false)
}

func buildSingleRuleDetector(r *config.Rule, includeAllowlists bool) *detect.Detector {
	// normalize keywords like in the config package
	var (
		uniqueKeywords = make(map[string]struct{})
		keywords       []string
	)
	for _, keyword := range r.Keywords {
		k := strings.ToLower(keyword)
		if _, ok := uniqueKeywords[k]; ok {
			continue
		}
		keywords = append(keywords, k)
		uniqueKeywords[k] = struct{}{}
	}
	r.Keywords = keywords

	rules := map[string]config.Rule{
		r.RuleID: *r,
	}
	cfg := base.CreateGlobalConfig()
	if !includeAllowlists {
		cfg.Allowlists = nil
	}
	cfg.Rules = rules
	cfg.Keywords = uniqueKeywords
	for _, a := range cfg.Allowlists {
		if err := a.Validate(); err != nil {
			logging.Fatal().Err(err).Msg("invalid global allowlist")
		}
	}
	return detect.NewDetector(cfg)
}
