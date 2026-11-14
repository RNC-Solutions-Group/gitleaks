package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReleaseWorkflowIsRNCScoped(t *testing.T) {
	raw, err := os.ReadFile("../.github/workflows/release.yml")
	require.NoError(t, err)
	content := string(raw)

	assert.Contains(t, content, "REGISTRY: ghcr.io")
	assert.Contains(t, content, "IMAGE_NAME: rnc-solutions-group/gitleaks")
	assert.Contains(t, content, "rnc-mod")
	assert.Contains(t, content, "sort -V")
	assert.Contains(t, content, "packages: write")
	assert.Contains(t, content, "contents: read")
	assert.NotContains(t, content, "zricethezav/gitleaks")
	assert.NotContains(t, content, "DOCKER_PASSWORD")
	assert.NotContains(t, content, "docker.io")
	assert.NotContains(t, content, "pull_request")
}
