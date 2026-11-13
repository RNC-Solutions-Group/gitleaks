package detect

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func Test_createScmLink(t *testing.T) {
	tests := map[string]struct {
		remote  *sources.RemoteInfo
		finding report.Finding
		want    string
	}{
		// None
		"no platform": {
			remote: &sources.RemoteInfo{
				Platform: scm.NoPlatform,
				Url:      "",
			},
			want: "",
		},

		// GitHub
		"github - single line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitHubPlatform,
				Url:      "https://github.com/gitleaks/test",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "metrics/% of sales/.env",
				StartLine: 25,
				EndLine:   25,
			},
			want: "https://github.com/gitleaks/test/blob/20553ad96a4a080c94a54d677db97eed8ce2560d/metrics/%25%20of%20sales/.env#L25",
		},
		"github - multi line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitHubPlatform,
				Url:      "https://github.com/gitleaks/test",
			},
			finding: report.Finding{
				Commit:    "7bad9f7654cf9701b62400281748c0e8efd97666",
				File:      "config.json",
				StartLine: 235,
				EndLine:   238,
			},
			want: "https://github.com/gitleaks/test/blob/7bad9f7654cf9701b62400281748c0e8efd97666/config.json#L235-L238",
		},
		"github - markdown": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitHubPlatform,
				Url:      "https://github.com/gitleaks/test",
			},
			finding: report.Finding{
				Commit:    "1fc8961d172f39ffb671766e472aa76f8d713e87",
				File:      "docs/guides/ecosystem/discordjs.MD",
				StartLine: 34,
				EndLine:   34,
			},
			want: "https://github.com/gitleaks/test/blob/1fc8961d172f39ffb671766e472aa76f8d713e87/docs/guides/ecosystem/discordjs.MD?plain=1#L34",
		},
		"github - jupyter notebook": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitHubPlatform,
				Url:      "https://github.com/gitleaks/test",
			},
			finding: report.Finding{
				Commit:    "8f56bd2369595bcadbb007e88ba294630fb05c7b",
				File:      "Cloud/IPYNB/Overlapping Recommendation algorithm _OCuLaR_.ipynb",
				StartLine: 293,
				EndLine:   293,
			},
			want: "https://github.com/gitleaks/test/blob/8f56bd2369595bcadbb007e88ba294630fb05c7b/Cloud/IPYNB/Overlapping%20Recommendation%20algorithm%20_OCuLaR_.ipynb?plain=1#L293",
		},

		// GitLab
		"gitlab - single line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitLabPlatform,
				Url:      "https://gitlab.com/example-org/example-group/gitleaks",
			},
			finding: report.Finding{
				Commit:    "213ffd1c9bfa906eb4c7731771132c58a4ca0139",
				File:      ".gitlab-ci.yml",
				StartLine: 41,
				EndLine:   41,
			},
			want: "https://gitlab.com/example-org/example-group/gitleaks/blob/213ffd1c9bfa906eb4c7731771132c58a4ca0139/.gitlab-ci.yml#L41",
		},
		"gitlab - multi line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GitLabPlatform,
				Url:      "https://gitlab.com/example-org/example-group/gitleaks",
			},
			finding: report.Finding{
				Commit:    "63410f74e23a4e51e1f60b9feb073b5d325af878",
				File:      ".vscode/launchSettings.json",
				StartLine: 6,
				EndLine:   8,
			},
			want: "https://gitlab.com/example-org/example-group/gitleaks/blob/63410f74e23a4e51e1f60b9feb073b5d325af878/.vscode/launchSettings.json#L6-8",
		},

		// Azure DevOps
		"azuredevops - single line": {
			remote: &sources.RemoteInfo{
				Platform: scm.AzureDevOpsPlatform,
				Url:      "https://dev.azure.com/exampleorganisation/exampleproject/_git/exampleRepository",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   25,
			},
			want: "https://dev.azure.com/exampleorganisation/exampleproject/_git/exampleRepository/commit/20553ad96a4a080c94a54d677db97eed8ce2560d?path=/examplefile.json&line=25&lineStartColumn=1&lineEndColumn=10000000&type=2&lineStyle=plain&_a=files",
		},

		// Azure DevOps
		"azuredevops - multi line": {
			remote: &sources.RemoteInfo{
				Platform: scm.AzureDevOpsPlatform,
				Url:      "https://dev.azure.com/exampleorganisation/exampleproject/_git/exampleRepository",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   30,
			},
			want: "https://dev.azure.com/exampleorganisation/exampleproject/_git/exampleRepository/commit/20553ad96a4a080c94a54d677db97eed8ce2560d?path=/examplefile.json&line=25&lineEnd=30&lineStartColumn=1&lineEndColumn=10000000&type=2&lineStyle=plain&_a=files",
		},

		// Gitea
		"gitea - single line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GiteaPlatform,
				Url:      "https://gitea.com/exampleorganisation/exampleproject",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   25,
			},
			want: "https://gitea.com/exampleorganisation/exampleproject/src/commit/20553ad96a4a080c94a54d677db97eed8ce2560d/examplefile.json#L25",
		},
		"gitea- multi line": {
			remote: &sources.RemoteInfo{
				Platform: scm.GiteaPlatform,
				Url:      "https://gitea.com/exampleorganisation/exampleproject",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   30,
			},
			want: "https://gitea.com/exampleorganisation/exampleproject/src/commit/20553ad96a4a080c94a54d677db97eed8ce2560d/examplefile.json#L25-L30",
		},
		"gitea - markdown": {
			remote: &sources.RemoteInfo{
				Platform: scm.GiteaPlatform,
				Url:      "https://gitea.com/exampleorganisation/exampleproject",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "Readme.md",
				StartLine: 34,
				EndLine:   34,
			},
			want: "https://gitea.com/exampleorganisation/exampleproject/src/commit/20553ad96a4a080c94a54d677db97eed8ce2560d/Readme.md?display=source#L34",
		},
		// bitbucket
		"bitbucket - single line": {
			remote: &sources.RemoteInfo{
				Platform: scm.BitbucketPlatform,
				Url:      "https://bitbucket.org/exampleorganisation/exampleproject",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   25,
			},
			want: "https://bitbucket.org/exampleorganisation/exampleproject/src/20553ad96a4a080c94a54d677db97eed8ce2560d/examplefile.json#lines-25",
		},
		"bitbucket- multi line": {
			remote: &sources.RemoteInfo{
				Platform: scm.BitbucketPlatform,
				Url:      "https://bitbucket.org/exampleorganisation/exampleproject",
			},
			finding: report.Finding{
				Commit:    "20553ad96a4a080c94a54d677db97eed8ce2560d",
				File:      "examplefile.json",
				StartLine: 25,
				EndLine:   30,
			},
			want: "https://bitbucket.org/exampleorganisation/exampleproject/src/20553ad96a4a080c94a54d677db97eed8ce2560d/examplefile.json#lines-25:30",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := createScmLink(tt.remote, tt.finding)
			assert.Equal(t, tt.want, actual)
		})
	}
}

// Test_printFinding_colorOutput guards the Lip Gloss v2 migration.
//
// Lip Gloss v2 always renders full-fidelity ANSI and leaves color down-sampling
// to the output layer. The finding display path therefore writes through a
// down-sampling writer: redirected or piped output must stay free of escape
// sequences, while --no-color must keep winning regardless of the environment.
func Test_printFinding_colorOutput(t *testing.T) {
	// Sanity check that makes the assertions below meaningful: v2 styles render
	// ANSI sequences even when the destination is not a terminal.
	rendered := lipgloss.NewStyle().Foreground(lipgloss.Color("#f05c07")).Render("secret")
	assert.Contains(t, rendered, "\x1b[", "Lip Gloss v2 should render ANSI sequences")

	t.Setenv("CLICOLOR_FORCE", "0")
	t.Setenv("NO_COLOR", "0")

	finding := report.Finding{
		RuleID:      "generic-api-key",
		Match:       `token = "abcdef123456"`,
		Secret:      "abcdef123456",
		Line:        `token = "abcdef123456"`,
		File:        "creds.txt",
		Fingerprint: "creds.txt:generic-api-key:1",
	}

	// captureStdout replaces os.Stdout with a pipe, i.e. a redirected target.
	redirected := captureStdout(func() { printFinding(finding, false) })
	assert.Contains(t, redirected, "generic-api-key")
	assert.Contains(t, redirected, "abcdef123456")
	assert.NotContains(t, redirected, "\x1b[", "redirected output must not contain ANSI escape sequences")

	noColor := captureStdout(func() { printFinding(finding, true) })
	assert.Contains(t, noColor, "abcdef123456")
	assert.NotContains(t, noColor, "\x1b[", "--no-color output must not contain ANSI escape sequences")
}

func TestFilterRedactsAdjacentSecrets(t *testing.T) {
	findings := []report.Finding{
		{
			RuleID:    "alpha",
			StartLine: 1,
			EndLine:   1,
			Match:     "ALPHA_12345678",
			Line:      "left=ALPHA_12345678 right=BETA_87654321",
			Secret:    "ALPHA_12345678",
		},
		{
			RuleID:    "beta",
			StartLine: 1,
			EndLine:   1,
			Match:     "BETA_87654321",
			Line:      "left=ALPHA_12345678 right=BETA_87654321",
			Secret:    "BETA_87654321",
		},
	}

	got := filter(findings, 100)
	if assert.Len(t, got, 2) {
		for _, finding := range got {
			assert.NotContains(t, finding.Line, "ALPHA_12345678")
			assert.NotContains(t, finding.Line, "BETA_87654321")
			assert.Equal(t, "REDACTED", finding.Secret)
			assert.Equal(t, "REDACTED", finding.Match)
		}
	}
}

func TestFilterNoRedactPreservesFinding(t *testing.T) {
	findings := []report.Finding{{
		RuleID:    "alpha",
		StartLine: 1,
		EndLine:   1,
		Match:     "ALPHA_12345678",
		Line:      "left=ALPHA_12345678",
		Secret:    "ALPHA_12345678",
	}}
	got := filter(findings, 0)
	if assert.Len(t, got, 1) {
		assert.Equal(t, "ALPHA_12345678", got[0].Secret)
		assert.Equal(t, "left=ALPHA_12345678", got[0].Line)
	}
}
