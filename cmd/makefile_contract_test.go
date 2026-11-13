package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakefileCommandContract guards the build/test boundary documented in the
// Makefile. Routine commands must not format, generate, tidy, or fetch from a
// remote; those operations have explicit maintenance targets.
func TestMakefileCommandContract(t *testing.T) {
	content, err := os.ReadFile("../Makefile")
	require.NoError(t, err)
	text := string(content)

	require.NotContains(t, text, "git fetch")
	assert.NotContains(t, text, "git tag |")

	build := makeTarget(t, text, "build")
	assert.Contains(t, build, "$(GO) build")
	for _, forbidden := range []string{"$(GO) fmt", "$(GO) mod tidy", "$(GO) generate", "git "} {
		assert.NotContains(t, build, forbidden, "build must not invoke %q", forbidden)
	}

	testTarget := makeTarget(t, text, "test")
	assert.Contains(t, testTarget, "$(GO) test")
	for _, forbidden := range []string{"$(GO) fmt", "$(GO) mod tidy", "$(GO) generate", "--race"} {
		assert.NotContains(t, testTarget, forbidden, "test must not invoke %q", forbidden)
	}
	assert.Contains(t, makeTarget(t, text, "test-race"), "-race")

	verify := makeTarget(t, text, "verify")
	assert.Contains(t, verify, "$(GO) vet")
	assert.Contains(t, verify, "TestGeneratedConfigMatchesTracked")
	assert.NotContains(t, verify, "$(GO) generate")

	assert.Contains(t, makeTarget(t, text, "format"), "$(GO) fmt")
	assert.Contains(t, makeTarget(t, text, "tidy"), "$(GO) mod tidy")
	assert.Contains(t, makeTarget(t, text, "generate"), "$(GO) generate")
}

func makeTarget(t *testing.T, content, name string) string {
	t.Helper()
	lines := strings.Split(content, "\n")
	for index, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line != name+":" {
			continue
		}
		var recipe strings.Builder
		for _, candidate := range lines[index+1:] {
			candidate = strings.TrimSuffix(candidate, "\r")
			if !strings.HasPrefix(candidate, "\t") {
				break
			}
			recipe.WriteString(candidate)
			recipe.WriteByte('\n')
		}
		return recipe.String()
	}
	t.Fatalf("target %q not found", name)
	return ""
}
