package report

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/zricethezav/gitleaks/v8/config"
)

type SarifReporter struct {
	OrderedRules []config.Rule
}

var _ Reporter = (*SarifReporter)(nil)

func (r *SarifReporter) Write(w io.WriteCloser, findings []Finding) error {
	sarif := Sarif{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs:    r.getRuns(findings),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	return encoder.Encode(sarif)
}

func (r *SarifReporter) getRuns(findings []Finding) []Runs {
	return []Runs{
		{
			Tool:    r.getTool(),
			Results: getResults(findings),
		},
	}
}

func (r *SarifReporter) getTool() Tool {
	tool := Tool{
		Driver: Driver{
			Name:            driver,
			SemanticVersion: version,
			InformationUri:  "https://github.com/gitleaks/gitleaks",
			Rules:           r.getRules(),
		},
	}

	// if this tool has no rules, ensure that it is represented as [] instead of null/nil
	if hasEmptyRules(tool) {
		tool.Driver.Rules = make([]Rules, 0)
	}

	return tool
}

func hasEmptyRules(tool Tool) bool {
	return len(tool.Driver.Rules) == 0
}

func (r *SarifReporter) getRules() []Rules {
	// TODO	for _, rule := range cfg.Rules {
	var rules []Rules
	for _, rule := range r.OrderedRules {
		rules = append(rules, Rules{
			ID: rule.RuleID,
			Description: ShortDescription{
				Text: rule.Description,
			},
		})
	}
	return rules
}

func messageText(f Finding) string {
	if f.Commit == "" {
		return fmt.Sprintf("%s has detected secret for file %s.", f.RuleID, f.File)
	}

	return fmt.Sprintf("%s has detected secret for file %s at commit %s.", f.RuleID, f.File, f.Commit)

}

func getResults(findings []Finding) []Results {
	results := []Results{}
	for _, f := range findings {
		r := Results{
			Message: Message{
				Text: messageText(f),
			},
			RuleId:    f.RuleID,
			Locations: getLocation(f),
			// This information goes in partial fingerprings until revision
			// data can be added somewhere else
			PartialFingerPrints: PartialFingerPrints{
				CommitSha:     f.Commit,
				Email:         f.Email,
				CommitMessage: f.Message,
				Date:          f.Date,
				Author:        f.Author,
			},
			Properties: Properties{
				Tags: f.Tags,
			},
		}
		results = append(results, r)
	}
	return results
}

// sarifArtifactURI converts a scanner path into a URI reference suitable for
// SARIF artifactLocation.uri. Spaces, percent signs, fragments, queries and
// non-ASCII bytes are escaped exactly once. Windows drive and UNC separators
// are normalized before escaping; POSIX backslashes remain literal path bytes.
func sarifArtifactURI(uri string) string {
	if uri == "" {
		return uri
	}
	if isWindowsDrivePath(uri) {
		// A drive path cannot be emitted bare because the drive letter would be
		// parsed as a URI scheme. A file URI preserves its absolute identity.
		path := "/" + strings.ReplaceAll(uri, `\`, "/")
		return (&url.URL{Scheme: "file", Path: path}).String()
	}
	if strings.HasPrefix(uri, `\\`) || strings.HasPrefix(uri, "//") {
		normalized := strings.ReplaceAll(uri, `\`, "/")
		normalized = strings.TrimPrefix(normalized, "//")
		host, path, _ := strings.Cut(normalized, "/")
		return (&url.URL{Scheme: "file", Host: host, Path: "/" + path}).String()
	}
	return (&url.URL{Path: uri}).String()
}

func isWindowsDrivePath(path string) bool {
	if len(path) < 3 || path[1] != ':' || (path[2] != '/' && path[2] != '\\') {
		return false
	}
	c := path[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func getLocation(f Finding) []Locations {
	uri := f.File
	if f.SymlinkFile != "" {
		uri = f.SymlinkFile
	}
	uri = sarifArtifactURI(uri)
	startColumn := f.StartColumn
	if f.StartColumnCodepoint > 0 {
		startColumn = f.StartColumnCodepoint
	}
	endColumn := f.EndColumn
	if f.EndColumnCodepoint > 0 {
		endColumn = f.EndColumnCodepoint
	}
	return []Locations{
		{
			PhysicalLocation: PhysicalLocation{
				ArtifactLocation: ArtifactLocation{
					URI: uri,
				},
				Region: Region{
					StartLine:   f.StartLine,
					EndLine:     f.EndLine,
					StartColumn: startColumn,
					EndColumn:   endColumn,
					Snippet: Snippet{
						Text: f.Secret,
					},
				},
			},
		},
	}
}

type PartialFingerPrints struct {
	CommitSha     string `json:"commitSha"`
	Email         string `json:"email"`
	Author        string `json:"author"`
	Date          string `json:"date"`
	CommitMessage string `json:"commitMessage"`
}

type Sarif struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Runs `json:"runs"`
}

type ShortDescription struct {
	Text string `json:"text"`
}

type FullDescription struct {
	Text string `json:"text"`
}

type Rules struct {
	ID          string           `json:"id"`
	Description ShortDescription `json:"shortDescription"`
}

type Driver struct {
	Name            string  `json:"name"`
	SemanticVersion string  `json:"semanticVersion"`
	InformationUri  string  `json:"informationUri"`
	Rules           []Rules `json:"rules"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Message struct {
	Text string `json:"text"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

type Region struct {
	StartLine   int     `json:"startLine"`
	StartColumn int     `json:"startColumn"`
	EndLine     int     `json:"endLine"`
	EndColumn   int     `json:"endColumn"`
	Snippet     Snippet `json:"snippet"`
}

type Snippet struct {
	Text string `json:"text"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           Region           `json:"region"`
}

type Locations struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type Properties struct {
	Tags []string `json:"tags"`
}

type Results struct {
	Message             Message     `json:"message"`
	RuleId              string      `json:"ruleId"`
	Locations           []Locations `json:"locations"`
	PartialFingerPrints `json:"partialFingerprints"`
	Properties          Properties `json:"properties"`
}

type Runs struct {
	Tool    Tool      `json:"tool"`
	Results []Results `json:"results"`
}
