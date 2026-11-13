package report

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestRedact(t *testing.T) {
	finding := Finding{
		Match:  "line containing secret",
		Secret: "secret",
	}
	finding.Redact(100)
	assert.Equal(t, "REDACTED", finding.Secret)
	assert.Equal(t, "line containing REDACTED", finding.Match)
}

func TestMask(t *testing.T) {
	tests := map[string]struct {
		finding Finding
		percent uint
		expect  Finding
	}{
		"normal secret": {
			finding: Finding{Match: "line containing secret", Secret: "secret"},
			expect:  Finding{Match: "line containing se...", Secret: "se..."},
			percent: 75,
		},
		"empty secret": {
			finding: Finding{Match: "line containing", Secret: ""},
			expect:  Finding{Match: "line...", Secret: "line..."},
			percent: 75,
		},
		"short secret": {
			finding: Finding{Match: "line containing", Secret: "ss"},
			expect:  Finding{Match: "line containing", Secret: "..."},
			percent: 75,
		},
		"path-only empty secret": {
			finding: Finding{Match: "file detected: secret.txt", Secret: ""},
			expect:  Finding{Match: "file detected: secret.txt", Secret: ""},
			percent: 100,
		},
		"empty secret full redaction": {
			finding: Finding{Match: "EMPTYTOKEN=SUPERSECRET", Line: "EMPTYTOKEN=SUPERSECRET", Secret: ""},
			expect:  Finding{Match: "REDACTED", Line: "REDACTED", Secret: "REDACTED"},
			percent: 100,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := test.finding
			e := test.expect
			f.Redact(test.percent)
			assert.Equal(t, e.Secret, f.Secret)
			assert.Equal(t, e.Match, f.Match)
			assert.Equal(t, e.Line, f.Line)
		})
	}
}

func TestRedactZeroIsNoop(t *testing.T) {
	finding := Finding{Match: "line containing secret", Line: "line containing secret", Secret: "secret"}
	finding.Redact(0)
	assert.Equal(t, "line containing secret", finding.Match)
	assert.Equal(t, "line containing secret", finding.Line)
	assert.Equal(t, "secret", finding.Secret)
}

func TestMaskSecret(t *testing.T) {
	tests := map[string]struct {
		secret  string
		percent uint
		expect  string
	}{
		"normal masking":  {secret: "secret", percent: 75, expect: "se..."},
		"high masking":    {secret: "secret", percent: 90, expect: "s..."},
		"low masking":     {secret: "secret", percent: 10, expect: "secre..."},
		"invalid masking": {secret: "secret", percent: 1000, expect: "..."},
		"multibyte":       {secret: "pässwörd", percent: 50, expect: "päss..."},
		"mixed width":     {secret: "pässwörd", percent: 75, expect: "pä..."},
		"emoji":           {secret: "🔐🔑🧪🧬", percent: 75, expect: "🔐..."},
		"combining":       {secret: "e\u0301x", percent: 50, expect: "e\u0301..."},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := maskSecret(test.secret, test.percent)
			assert.Equal(t, test.expect, got)
			assert.True(t, utf8.ValidString(got), "masked secret must be valid UTF-8")
		})
	}
}

func TestRedactUnicodeKeepsValidUTF8(t *testing.T) {
	for _, percent := range []uint{10, 45, 50, 75, 90, 100} {
		finding := Finding{
			Match:  "token = 🔐sëcret🔑",
			Line:   "token = 🔐sëcret🔑",
			Secret: "🔐sëcret🔑",
		}
		finding.Redact(percent)
		assert.True(t, utf8.ValidString(finding.Match), "percent %d match", percent)
		assert.True(t, utf8.ValidString(finding.Line), "percent %d line", percent)
		assert.True(t, utf8.ValidString(finding.Secret), "percent %d secret", percent)
		assert.NotContains(t, finding.Match, "🔐sëcret🔑")
		assert.NotContains(t, finding.Line, "🔐sëcret🔑")
	}
}

func TestRedactWithMultipleSecrets(t *testing.T) {
	finding := Finding{
		Match:  "left=ALPHA_12345678 right=BETA_87654321",
		Line:   "left=ALPHA_12345678 right=BETA_87654321",
		Secret: "ALPHA_12345678",
	}
	finding.RedactWith(100, []string{"ALPHA_12345678", "BETA_87654321"})
	assert.NotContains(t, finding.Match, "ALPHA_12345678")
	assert.NotContains(t, finding.Match, "BETA_87654321")
	assert.NotContains(t, finding.Line, "ALPHA_12345678")
	assert.NotContains(t, finding.Line, "BETA_87654321")
	assert.Equal(t, "REDACTED", finding.Secret)
}

func TestRedactWithRepeatedSecrets(t *testing.T) {
	finding := Finding{
		Match:  "secret secret",
		Line:   "secret secret",
		Secret: "secret",
	}
	finding.RedactWith(100, nil)
	assert.Equal(t, "REDACTED REDACTED", finding.Match)
	assert.Equal(t, "REDACTED REDACTED", finding.Line)
}

func TestRedactWithOverlappingSecretsDoesNotPanicOrLeak(t *testing.T) {
	finding := Finding{
		Match:  "abcdef",
		Line:   "abcdef",
		Secret: "abc",
	}
	finding.RedactWith(100, []string{"abcdef", "abc"})
	assert.NotContains(t, finding.Match, "abcdef")
	assert.NotContains(t, finding.Line, "abcdef")
	assert.Equal(t, "REDACTED", finding.Secret)
}

func TestRedactRequiredFindings(t *testing.T) {
	finding := Finding{
		Match:  "PRIMARY_12345678",
		Line:   "PRIMARY_12345678",
		Secret: "PRIMARY_12345678",
	}
	finding.AddRequiredFindings([]*RequiredFinding{{
		RuleID:    "aux",
		StartLine: 2,
		Match:     "AUX_87654321",
		Line:      "AUX_87654321",
		Secret:    "AUX_87654321",
	}})

	finding.Redact(100)
	assert.Equal(t, "REDACTED", finding.Secret)
	assert.Equal(t, "REDACTED", finding.Match)
	assert.Equal(t, "REDACTED", finding.requiredFindings[0].Secret)
	assert.Equal(t, "REDACTED", finding.requiredFindings[0].Match)
	assert.Equal(t, "REDACTED", finding.requiredFindings[0].Line)
}

func TestRedactControlCharacters(t *testing.T) {
	finding := Finding{
		Match:  "\x1b[31msecret\x1b[0m",
		Line:   "\x1b[31msecret\x1b[0m",
		Secret: "secret",
	}
	finding.Redact(100)
	assert.NotContains(t, finding.Match, "secret")
	assert.NotContains(t, finding.Line, "secret")
	assert.Equal(t, "REDACTED", finding.Secret)
}

func TestRedactNewlineSecret(t *testing.T) {
	finding := Finding{
		Match:  "key=line\nsecret",
		Line:   "key=line\nsecret",
		Secret: "line\nsecret",
	}
	finding.Redact(100)
	assert.NotContains(t, finding.Match, "line\nsecret")
	assert.NotContains(t, finding.Line, "line\nsecret")
	assert.Equal(t, "REDACTED", finding.Secret)
}

func TestRedactEmptyFindingDoesNotPanic(t *testing.T) {
	var finding Finding
	finding.Redact(100)
	assert.Empty(t, finding.Match)
	assert.Empty(t, finding.Line)
	assert.Empty(t, finding.Secret)
}
