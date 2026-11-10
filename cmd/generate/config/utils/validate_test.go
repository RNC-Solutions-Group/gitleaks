package utils

import (
	"testing"

	"github.com/zricethezav/gitleaks/v8/config"
	gitleaksregexp "github.com/zricethezav/gitleaks/v8/regexp"
)

func TestDiagnoseTruePositiveRegexMismatch(t *testing.T) {
	r := &config.Rule{
		RuleID: "regex-mismatch",
		Regex:  gitleaksregexp.MustCompile(`[A-Z]{8}`),
	}

	diagnostic, _, hasEntropy := diagnoseTruePositive(r, "lowercase")
	if diagnostic != "regex/structure mismatch" {
		t.Fatalf("diagnostic = %q, want regex/structure mismatch", diagnostic)
	}
	if hasEntropy {
		t.Fatal("regex mismatch should not report entropy diagnostics")
	}
}

func TestDiagnoseTruePositiveEntropyRejection(t *testing.T) {
	r := &config.Rule{
		RuleID:  "entropy-rejection",
		Regex:   gitleaksregexp.MustCompile(`[a-z0-9]{24}`),
		Entropy: 3.5,
	}

	diagnostic, entropy, hasEntropy := diagnoseTruePositive(r, "aaaaaaaaaaaaaaaaaaaaaaaa")
	if diagnostic != "entropy rejection" {
		t.Fatalf("diagnostic = %q, want entropy rejection", diagnostic)
	}
	if !hasEntropy {
		t.Fatal("entropy rejection should include entropy diagnostics")
	}
	if entropy != 0 {
		t.Fatalf("entropy = %f, want 0", entropy)
	}
}

func TestDiagnoseTruePositiveAllowlistSuppression(t *testing.T) {
	allowlist := &config.Allowlist{
		Regexes: []*gitleaksregexp.Regexp{
			gitleaksregexp.MustCompile(`^[a-z0-9]{24}$`),
		},
	}
	if err := allowlist.Validate(); err != nil {
		t.Fatalf("failed to validate test allowlist: %v", err)
	}

	r := &config.Rule{
		RuleID:     "allowlist-suppression",
		Regex:      gitleaksregexp.MustCompile(`[a-z0-9]{24}`),
		Allowlists: []*config.Allowlist{allowlist},
	}

	diagnostic, _, hasEntropy := diagnoseTruePositive(r, "abcdefghijklmnopqrstuvwx")
	if diagnostic != "allowlist suppression" {
		t.Fatalf("diagnostic = %q, want allowlist suppression", diagnostic)
	}
	if hasEntropy {
		t.Fatal("allowlist suppression should not report entropy diagnostics")
	}
}
