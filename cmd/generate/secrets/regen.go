// Package reggen generates text based on regex definitions.
// This is a slightly altered version of https://github.com/lucasjones/reggen.
package secrets

import (
	"fmt"
	"math"

	"github.com/lucasjones/reggen"
	"github.com/zricethezav/gitleaks/v8/logging"
)

const (
	// DefaultSeed makes generated validation fixtures reproducible across runs.
	DefaultSeed int64 = 1
	// DefaultMaxAttempts bounds rejection sampling for entropy-aware fixtures.
	DefaultMaxAttempts = 1000
)

// Requirements describes the candidate-generation contract for one fixture.
//
// Regex generates only the variable portion of the secret. Prefix is prepended
// before entropy is calculated, so callers can model patterns such as
// "cmVmd" + [A-Za-z0-9]{59}, where the detector evaluates the full match.
type Requirements struct {
	Regex       string
	Prefix      string
	MinEntropy  float64
	Seed        int64
	MaxAttempts int
	Label       string
}

// ShannonEntropy mirrors detect.shannonEntropy exactly. The generator must use
// the same full-secret content and the same byte-length denominator as the
// detector; otherwise a fixture can pass generation and still be rejected by
// the production entropy check.
func ShannonEntropy(data string) float64 {
	if data == "" {
		return 0
	}

	charCounts := make(map[rune]int)
	for _, char := range data {
		charCounts[char]++
	}

	invLength := 1.0 / float64(len(data))
	var entropy float64
	for _, count := range charCounts {
		freq := float64(count) * invLength
		entropy -= freq * math.Log2(freq)
	}
	return entropy
}

// GenerateSecret generates a deterministic candidate from Requirements.Regex,
// prepends Requirements.Prefix, and rejects candidates whose full-secret
// Shannon entropy does not exceed Requirements.MinEntropy. The detector treats
// entropy equal to the threshold as a rejection, so this helper does the same.
//
// Only candidate-generation requirements are enforced here. The caller must
// still run the existing detector validator; this function does not retry
// arbitrary detector failures.
func GenerateSecret(req Requirements) (string, error) {
	g, err := reggen.NewGenerator(req.Regex)
	if err != nil {
		return "", fmt.Errorf("generate secret: invalid regex %q: %w", req.Regex, err)
	}

	attempts := req.MaxAttempts
	if attempts <= 0 {
		attempts = DefaultMaxAttempts
	}

	g.SetSeed(req.Seed)
	lastEntropy := 0.0
	for attempt := 1; attempt <= attempts; attempt++ {
		candidate := req.Prefix + g.Generate(1)
		lastEntropy = ShannonEntropy(candidate)
		// Entropy 0 disables the detector check. Preserve that behavior for
		// rules that do not declare an entropy threshold.
		if req.MinEntropy <= 0 || lastEntropy > req.MinEntropy {
			return candidate, nil
		}
	}

	label := req.Label
	if label == "" {
		label = "secret"
	}
	return "", fmt.Errorf(
		"%s: no candidate exceeded entropy %.4f after %d attempts (last candidate entropy %.4f)",
		label, req.MinEntropy, attempts, lastEntropy,
	)
}

// MustGenerateSecret is the rule-fixture entry point. It keeps generation
// failures visible and fatal, but does not print the candidate value.
func MustGenerateSecret(req Requirements) string {
	secret, err := GenerateSecret(req)
	if err != nil {
		logging.Fatal().Err(err).Msg("failed to generate positive fixture")
	}
	return secret
}

// NewSecret preserves the existing unqualified generator for callers that do
// not declare an entropy condition. It is deterministic by default so repeated
// canonical generation runs are reproducible.
func NewSecret(regex string) string {
	secret, err := GenerateSecret(Requirements{
		Regex:       regex,
		Seed:        DefaultSeed,
		MaxAttempts: 1,
	})
	if err != nil {
		panic(err)
	}
	return secret
}
