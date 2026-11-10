package secrets

import (
	"math"
	"regexp"
	"strings"
	"testing"
)

func TestShannonEntropy(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  float64
	}{
		{name: "empty", value: "", want: 0},
		{name: "single character", value: "a", want: 0},
		{name: "two symbols", value: "ab", want: 1},
		{name: "two symbols balanced", value: "aabb", want: 1},
		{name: "four symbols", value: "abcd", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShannonEntropy(tt.value)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Fatalf("ShannonEntropy(%q) = %f, want %f", tt.value, got, tt.want)
			}
		})
	}
}

func TestGenerateSecretReproducible(t *testing.T) {
	req := Requirements{
		Regex:       `[a-z0-9]{20}[a-f0-9]{4}`,
		MinEntropy:  3.5,
		Seed:        1,
		MaxAttempts: DefaultMaxAttempts,
		Label:       "atlassian-api-token",
	}

	first, err := GenerateSecret(req)
	if err != nil {
		t.Fatalf("first generation failed: %v", err)
	}
	second, err := GenerateSecret(req)
	if err != nil {
		t.Fatalf("second generation failed: %v", err)
	}
	if first != second {
		t.Fatalf("same seed produced different candidates")
	}
	if len(first) != 24 {
		t.Fatalf("generated token length = %d, want 24", len(first))
	}
	if !regexp.MustCompile(`^[a-z0-9]{20}[a-f0-9]{4}$`).MatchString(first) {
		t.Fatalf("generated token does not match the requested structure")
	}
	if got := ShannonEntropy(first); got <= req.MinEntropy {
		t.Fatalf("generated token entropy = %f, want > %f", got, req.MinEntropy)
	}
}

func TestGenerateSecretIncludesPrefixInEntropy(t *testing.T) {
	req := Requirements{
		Regex:       `[A-Za-z0-9]{59}`,
		Prefix:      "cmVmd",
		MinEntropy:  4.5,
		Seed:        1,
		MaxAttempts: DefaultMaxAttempts,
		Label:       "artifactory-reference-token",
	}

	token, err := GenerateSecret(req)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}
	if !strings.HasPrefix(token, req.Prefix) {
		t.Fatalf("generated token does not include the required prefix")
	}
	if len(token) != 64 {
		t.Fatalf("generated token length = %d, want 64", len(token))
	}
	if !regexp.MustCompile(`^cmVmd[A-Za-z0-9]{59}$`).MatchString(token) {
		t.Fatalf("generated token does not match the reference-token structure")
	}
	if got := ShannonEntropy(token); got <= req.MinEntropy {
		t.Fatalf("generated token entropy = %f, want > %f", got, req.MinEntropy)
	}
}

func TestGenerateSecretThresholdBoundary(t *testing.T) {
	// "aabb" has exactly 1.0 bits of Shannon entropy. The detector rejects
	// candidates at or below the configured threshold, so 1.0 must fail and
	// 0.999 must succeed.
	req := Requirements{
		Regex:       `aabb`,
		MinEntropy:  1.0,
		Seed:        1,
		MaxAttempts: 3,
		Label:       "boundary-at-threshold",
	}
	if _, err := GenerateSecret(req); err == nil {
		t.Fatal("expected entropy equal to the threshold to be rejected")
	}

	req.MinEntropy = 0.999
	token, err := GenerateSecret(req)
	if err != nil {
		t.Fatalf("expected candidate above threshold to succeed: %v", err)
	}
	if token != "aabb" {
		t.Fatalf("generated token = %q, want aabb", token)
	}
}

func TestGenerateSecretExhaustion(t *testing.T) {
	_, err := GenerateSecret(Requirements{
		Regex:       `a{10}`,
		MinEntropy:  1.0,
		Seed:        1,
		MaxAttempts: 5,
		Label:       "exhaustion",
	})
	if err == nil {
		t.Fatal("expected bounded generation to fail")
	}
	message := err.Error()
	if !strings.Contains(message, "after 5 attempts") {
		t.Fatalf("error does not report the attempt bound: %v", err)
	}
	if !strings.Contains(message, "entropy") {
		t.Fatalf("error does not report the entropy condition: %v", err)
	}
	if strings.Contains(message, "aaaaaaaaaa") {
		t.Fatalf("error unexpectedly contains the generated candidate: %v", err)
	}
}

func TestGenerateSecretSeedSweep(t *testing.T) {
	tests := []struct {
		name       string
		req        Requirements
		pattern    *regexp.Regexp
		wantLength int
	}{
		{
			name: "atlassian-api-token",
			req: Requirements{
				Regex:       `[a-z0-9]{20}[a-f0-9]{4}`,
				MinEntropy:  3.5,
				MaxAttempts: DefaultMaxAttempts,
				Label:       "atlassian-api-token",
			},
			pattern:    regexp.MustCompile(`^[a-z0-9]{20}[a-f0-9]{4}$`),
			wantLength: 24,
		},
		{
			name: "artifactory-reference-token",
			req: Requirements{
				Regex:       `[a-z0-9]{59}`,
				Prefix:      "cmVmd",
				MinEntropy:  4.5,
				MaxAttempts: DefaultMaxAttempts,
				Label:       "artifactory-reference-token",
			},
			pattern:    regexp.MustCompile(`^cmVmd[a-z0-9]{59}$`),
			wantLength: 64,
		},
		{
			name: "artifactory-api-key",
			req: Requirements{
				Regex:       `[a-z0-9]{69}`,
				Prefix:      "AKCp",
				MinEntropy:  4.5,
				MaxAttempts: DefaultMaxAttempts,
				Label:       "artifactory-api-key",
			},
			pattern:    regexp.MustCompile(`^AKCp[a-z0-9]{69}$`),
			wantLength: 73,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for seed := int64(1); seed <= 100; seed++ {
				tt.req.Seed = seed
				token, err := GenerateSecret(tt.req)
				if err != nil {
					t.Fatalf("seed %d failed: %v", seed, err)
				}
				if len(token) != tt.wantLength {
					t.Fatalf("seed %d length = %d, want %d", seed, len(token), tt.wantLength)
				}
				if !tt.pattern.MatchString(token) {
					t.Fatalf("seed %d token does not match the requested structure", seed)
				}
				if got := ShannonEntropy(token); got <= tt.req.MinEntropy {
					t.Fatalf("seed %d entropy = %f, want > %f", seed, got, tt.req.MinEntropy)
				}
			}
		})
	}
}

func TestNewSecretDeterministicWithoutEntropy(t *testing.T) {
	first := NewSecret(`[a-z]{8}`)
	second := NewSecret(`[a-z]{8}`)
	if first != second {
		t.Fatal("NewSecret should be deterministic for canonical generation")
	}
	if !regexp.MustCompile(`^[a-z]{8}$`).MatchString(first) {
		t.Fatalf("NewSecret produced an invalid candidate")
	}
}
