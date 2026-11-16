package rules

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/utils"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/secrets"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func AirtableApiKey() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Uncovered a possible Airtable API Key, potentially compromising database access and leading to data leakage or alteration.",
		RuleID:      "airtable-api-key",
		Regex:       utils.GenerateSemiGenericRegex([]string{"airtable"}, utils.AlphaNumeric("17"), true),
		Keywords:    []string{"airtable"},
	}

	// validate
	tps := utils.GenerateSampleSecrets("airtable", secrets.NewSecret(utils.AlphaNumeric("17")))
	return utils.Validate(r, tps, nil)
}

func AirtablePersonalAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Uncovered a possible Airtable Personal AccessToken, potentially compromising database access and leading to data leakage or alteration.",
		RuleID:      "airtable-personnal-access-token",
		Regex:       regexp.MustCompile(`\b(pat[[:alnum:]]{14}\.[a-f0-9]{64})\b`),
		// The token prefix itself is the only literal guaranteed to appear in a
		// match. Requiring the unrelated product name "airtable" made the
		// keyword prefilter suppress bare PATs before the regex ran (#2257).
		Keywords: []string{"pat"},
	}

	// validate
	token := "pat" + secrets.NewSecret(utils.AlphaNumeric("14")+"\\."+utils.Hex("64"))
	tps := utils.GenerateSampleSecrets("airtable", token)
	tps = append(tps, token)
	return utils.Validate(r, tps, nil)
}
