package report

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/zricethezav/gitleaks/v8/sources"
)

// ColorWriter returns a writer for standard output that down-samples ANSI
// color sequences to the color profile supported by the current terminal.
//
// Lip Gloss v2 always renders full-fidelity ANSI and leaves color down-sampling
// to the output layer. Writing styled output through this writer preserves the
// v1 behavior where redirected or piped output contains no escape sequences.
func ColorWriter() io.Writer {
	return colorprofile.NewWriter(os.Stdout, os.Environ())
}

// Finding contains a whole bunch of information about a secret finding.
// Plenty of real estate in this bad boy so fillerup as needed.
type Finding struct {
	// Rule is the name of the rule that was matched
	RuleID      string
	Description string

	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int

	// SARIF defines columns in Unicode codepoints. These values are populated
	// from the corrected line start for SARIF output while the historical
	// byte-based StartColumn/EndColumn contract remains unchanged.
	StartColumnCodepoint int `json:"-"`
	EndColumnCodepoint   int `json:"-"`

	Line string `json:"-"`

	Match string

	// Captured secret
	Secret string

	// File is the name of the file containing the finding
	File        string
	SymlinkFile string
	Commit      string
	Link        string `json:",omitempty"`

	// Entropy is the shannon entropy of Value
	Entropy float32

	Author  string
	Email   string
	Date    string
	Message string
	Tags    []string

	// unique identifier
	Fingerprint string

	// Fragment used for multi-part rule checking, CEL filtering,
	// and eventually ML validation
	Fragment *sources.Fragment `json:",omitempty"`

	// TODO keeping private for now to during experimental phase
	requiredFindings []*RequiredFinding
}

type RequiredFinding struct {
	// contains a subset of the Finding fields
	// only used for reporting
	RuleID      string
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Line        string `json:"-"`
	Match       string
	Secret      string
}

func (f *Finding) AddRequiredFindings(afs []*RequiredFinding) {
	if f.requiredFindings == nil {
		f.requiredFindings = make([]*RequiredFinding, 0)
	}
	f.requiredFindings = append(f.requiredFindings, afs...)
}

// Redact removes sensitive information from a finding.
func (f *Finding) Redact(percent uint) {
	f.RedactWith(percent, nil)
}

// RedactWith removes the finding's secret and any other known secrets that
// share its line. Callers that have several findings for one line should pass
// all of their secrets so that a verbose finding does not expose a neighboring
// credential.
func (f *Finding) RedactWith(percent uint, secrets []string) {
	if percent == 0 {
		return
	}

	known := make([]string, 0, len(secrets)+1)
	known = append(known, secrets...)
	if f.Secret != "" {
		known = append(known, f.Secret)
	}
	for _, aux := range f.requiredFindings {
		if aux != nil && aux.Secret != "" {
			known = append(known, aux.Secret)
		}
	}
	known = uniqueSecrets(known)

	if f.Secret == "" {
		f.Line, f.Match, f.Secret = redactEmptySecret(f.Line, f.Match, percent, known)
	} else {
		hadSecretInMatch := strings.Contains(f.Match, f.Secret)
		f.Line = redactSecrets(f.Line, known, percent)
		f.Match = redactSecrets(f.Match, known, percent)
		f.Secret = maskedSecret(f.Secret, percent)

		// If a longer known secret subsumed the selected secret, the selected
		// mask may no longer be a substring of Match. Keep the rendered fields
		// consistent without reintroducing the original value.
		if hadSecretInMatch && !strings.Contains(f.Match, f.Secret) {
			f.Match = "REDACTED"
			f.Secret = "REDACTED"
		}
	}

	for _, aux := range f.requiredFindings {
		if aux == nil {
			continue
		}
		auxKnown := append([]string(nil), known...)
		if aux.Secret != "" {
			auxKnown = append(auxKnown, aux.Secret)
		}
		auxKnown = uniqueSecrets(auxKnown)

		if aux.Secret == "" {
			aux.Line, aux.Match, aux.Secret = redactEmptySecret(aux.Line, aux.Match, percent, auxKnown)
			continue
		}
		hadSecretInMatch := strings.Contains(aux.Match, aux.Secret)
		aux.Line = redactSecrets(aux.Line, auxKnown, percent)
		aux.Match = redactSecrets(aux.Match, auxKnown, percent)
		aux.Secret = maskedSecret(aux.Secret, percent)
		if hadSecretInMatch && !strings.Contains(aux.Match, aux.Secret) {
			aux.Match = "REDACTED"
			aux.Secret = "REDACTED"
		}
	}
}

// redactEmptySecret handles findings whose selected capture group is empty.
// Replacing an empty string would insert the mask between every character.
// Path-only findings have no secret and are left alone; content findings with
// an empty capture can still carry a credential in Match, so the whole match is
// masked rather than exposed.
func redactEmptySecret(line, match string, percent uint, known []string) (string, string, string) {
	if line == "" && strings.HasPrefix(match, "file detected: ") {
		return line, match, ""
	}
	if match == "" {
		line = redactSecrets(line, known, percent)
		if line != "" {
			if percent >= 100 {
				line = "REDACTED"
				return line, "REDACTED", "REDACTED"
			}
			line = maskSecret(line, percent)
		}
		return line, match, ""
	}

	masked := maskedSecret(match, percent)
	line = strings.ReplaceAll(line, match, masked)
	line = redactSecrets(line, known, percent)
	return line, masked, masked
}

func uniqueSecrets(secrets []string) []string {
	seen := make(map[string]struct{}, len(secrets))
	unique := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		if _, ok := seen[secret]; ok {
			continue
		}
		seen[secret] = struct{}{}
		unique = append(unique, secret)
	}
	// Longest first prevents a shorter overlapping secret from consuming part
	// of a longer one before the longer replacement is applied.
	sort.SliceStable(unique, func(i, j int) bool {
		return len(unique[i]) > len(unique[j])
	})
	return unique
}

func redactSecrets(value string, secrets []string, percent uint) string {
	if value == "" || len(secrets) == 0 {
		return value
	}
	replacements := make([]string, 0, len(secrets)*2)
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		replacements = append(replacements, secret, maskedSecret(secret, percent))
	}
	if len(replacements) == 0 {
		return value
	}
	return strings.NewReplacer(replacements...).Replace(value)
}

func maskedSecret(secret string, percent uint) string {
	if percent >= 100 {
		return "REDACTED"
	}
	return maskSecret(secret, percent)
}

func maskSecret(secret string, percent uint) string {
	if percent > 100 {
		percent = 100
	}
	if percent == 0 {
		return secret
	}
	// Count runes so a partial mask cannot split a multibyte character. This
	// rune-based masking follows gitleaks/gitleaks#2261 (@arpitjain099).
	runes := []rune(secret)
	length := float64(len(runes))
	if length <= 0 {
		return secret
	}
	prc := float64(100 - percent)
	lth := int64(math.RoundToEven(length * prc / float64(100)))
	if lth < 0 {
		lth = 0
	}
	if lth > int64(len(runes)) {
		lth = int64(len(runes))
	}

	return string(runes[:lth]) + "..."
}

func truncateRunes(value string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max-3]) + "..."
}

func (f *Finding) PrintRequiredFindings() {
	if len(f.requiredFindings) == 0 {
		return
	}

	out := ColorWriter()
	fmt.Fprintf(out, "%-12s ", "Required:")

	// Create orange style for secrets
	orangeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#bf9478"))

	for i, aux := range f.requiredFindings {
		auxSecret := strings.TrimSpace(aux.Secret)
		// Truncate long secrets for readability without splitting a rune.
		if len([]rune(auxSecret)) > 40 {
			auxSecret = truncateRunes(auxSecret, 40)
		}

		// Format: rule-id:line:secret
		if i == 0 {
			fmt.Fprintf(out, "%s:%d:%s\n", aux.RuleID, aux.StartLine, orangeStyle.Render(auxSecret))
		} else {
			fmt.Fprintf(out, "%-12s %s:%d:%s\n", "", aux.RuleID, aux.StartLine, orangeStyle.Render(auxSecret))
		}
	}
}
