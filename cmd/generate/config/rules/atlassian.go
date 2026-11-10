package rules

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/utils"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/secrets"
	"github.com/zricethezav/gitleaks/v8/config"
)

func Atlassian() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Detected an Atlassian API token, posing a threat to project management and collaboration tool security and data confidentiality.",
		RuleID:      "atlassian-api-token",
		Regex: utils.MergeRegexps(
			utils.GenerateSemiGenericRegex(
				[]string{"(?-i:ATLASSIAN|[Aa]tlassian)", "(?-i:CONFLUENCE|[Cc]onfluence)", "(?-i:JIRA|[Jj]ira)"},
				`[a-z0-9]{20}[a-f0-9]{4}`, // The last 4 characters are an MD5 hash.
				true,
			),
			utils.GenerateUniqueTokenRegex(`ATATT3[A-Za-z0-9_\-=]{186}`, false),
		),
		Entropy:  3.5,
		Keywords: []string{"atlassian", "confluence", "jira", "atatt3"},
	}

	// validate
	// Generate against the declared entropy threshold so a regex-valid but
	// low-entropy candidate is rejected before the detector validator runs.
	sampleToken := secrets.MustGenerateSecret(secrets.Requirements{
		Regex:       utils.AlphaNumeric("20") + utils.Hex("4"),
		MinEntropy:  r.Entropy,
		Seed:        secrets.DefaultSeed,
		MaxAttempts: secrets.DefaultMaxAttempts,
		Label:       r.RuleID,
	})
	tps := utils.GenerateSampleSecrets("atlassian", sampleToken)
	tps = append(tps, utils.GenerateSampleSecrets("confluence", sampleToken)...)
	tps = append(tps, utils.GenerateSampleSecrets("jira", sampleToken)...)
	tps = append(tps, `JIRA_API_TOKEN=`+sampleToken)
	tps = append(tps, utils.GenerateSampleSecrets("jira", "ATATT3xFfGF0K3irG5tKKi-6u-wwaXQFeGwZ-IHR-hQ3CulkKtMSuteRQFfLZ6jihHThzZCg_UjnDt-4Wl_gIRf4zrZJs5JqaeuBhsfJ4W5GD6yGg3W7903gbvaxZPBjxIQQ7BgFDSkPS8oPispw4KLz56mdK-G6CIvLO6hHRrZHY0Q3tvJ6JxE=C63992E6")...)

	fps := []string{"getPagesInConfluenceSpace,searchConfluenceUsingCql"}

	return utils.Validate(r, tps, fps)
}
