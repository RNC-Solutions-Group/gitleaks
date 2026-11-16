package detect

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/report"
)

func findingsForRule(findings []report.Finding, ruleID string) []report.Finding {
	var out []report.Finding
	for _, finding := range findings {
		if finding.RuleID == ruleID {
			out = append(out, finding)
		}
	}
	return out
}

func TestC1DefaultRuleRegressions(t *testing.T) {
	d, err := NewDetectorDefaultConfig()
	require.NoError(t, err)

	t.Run("Airtable PAT is detected without the product name", func(t *testing.T) {
		// #2257: the keyword prefilter previously required the unrelated
		// literal "airtable" and suppressed this self-signaling token.
		token := "patABCDEFGHIJKLMN." + strings.Repeat("a", 64)
		findings := d.Detect(Fragment{Raw: "credential=" + token})
		got := findingsForRule(findings, "airtable-personnal-access-token")
		require.Len(t, got, 1)
		assert.Equal(t, token, got[0].Secret)
	})

	t.Run("generic Base64 values with leading plus and slash", func(t *testing.T) {
		// #2327: both prefixed values are valid padded Base64 and must be
		// reported; the path controls below must remain negative.
		value := "hwnIHN4aRNJTOGA3MsdyoGFspvPQqmoYza1o7CXwbaD30TxxAfxnnCRrZLYvYeJ8H+C/ocyzk8n5UiBGBaelg=="
		for _, prefix := range []string{"+", "/"} {
			t.Run(prefix, func(t *testing.T) {
				findings := d.Detect(Fragment{Raw: `cookie_key = "` + prefix + value + `"`})
				got := findingsForRule(findings, "generic-api-key")
				require.Len(t, got, 1)
				assert.Equal(t, prefix+value, got[0].Secret)
			})
		}
	})

	t.Run("slash-led paths remain negative", func(t *testing.T) {
		for _, raw := range []string{
			`cookie_key = "/opts/config/alias/"`,
			`cookie_key = "/usr/local/bin/example-tool"`,
		} {
			findings := d.Detect(Fragment{Raw: raw})
			assert.Empty(t, findingsForRule(findings, "generic-api-key"), raw)
		}
	})
}

func TestC1Base64UnicodeFinding(t *testing.T) {
	// #2000: the decoded value contains printable Unicode before the token.
	// The detector must decode it and match the rule inside the decoded text.
	encoded := base64.StdEncoding.EncodeToString([]byte("café SECRET-12345"))
	rule := config.Rule{
		RuleID:      "unicode-secret",
		Description: "Unicode test secret",
		Regex:       regexp.MustCompile(`SECRET-[0-9]+`),
		Keywords:    []string{"token"},
	}
	cfg := config.Config{
		Rules:        map[string]config.Rule{rule.RuleID: rule},
		Keywords:     map[string]struct{}{"token": {}},
		OrderedRules: []string{rule.RuleID},
	}
	d := NewDetector(cfg)
	d.MaxDecodeDepth = 3

	findings := d.Detect(Fragment{Raw: "token: " + encoded})
	require.Len(t, findings, 1)
	assert.Equal(t, "SECRET-12345", findings[0].Secret)
	assert.Contains(t, findings[0].Tags, "decoded:base64")
}
