package rules

import (
	"testing"

	"github.com/zricethezav/gitleaks/v8/config"
)

func TestAffectedProviderFixtures(t *testing.T) {
	tests := []struct {
		name     string
		generate func() *config.Rule
	}{
		{name: "atlassian-api-token", generate: Atlassian},
		{name: "artifactory-api-key", generate: ArtifactoryApiKey},
		{name: "artifactory-reference-token", generate: ArtifactoryReferenceToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := tt.generate()
			if rule == nil {
				t.Fatalf("%s returned a nil rule", tt.name)
			}
			if rule.Entropy <= 0 {
				t.Fatalf("%s entropy threshold = %f, want > 0", tt.name, rule.Entropy)
			}
		})
	}
}
