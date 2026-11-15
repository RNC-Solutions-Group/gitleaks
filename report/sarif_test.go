package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/zricethezav/gitleaks/v8/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSarifArtifactURI(t *testing.T) {
	tests := map[string]string{
		"ApplicationCode/Java/sensitive data exposure/keystore.p12": "ApplicationCode/Java/sensitive%20data%20exposure/keystore.p12",
		"reports/a%b#c?.txt":     "reports/a%25b%23c%3F.txt",
		"src/秘密.py":              "src/%E7%A7%98%E5%AF%86.py",
		"/tmp/reports/a b.txt":   "/tmp/reports/a%20b.txt",
		`C:\Users\name\a b.txt`:  "file:///C:/Users/name/a%20b.txt",
		`\\server\share\a b.txt`: "file://server/share/a%20b.txt",
		"relative/clean.txt":     "relative/clean.txt",
	}
	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, want, sarifArtifactURI(input))
		})
	}
}

func TestGetLocationUsesCodepointColumns(t *testing.T) {
	locations := getLocation(Finding{
		File:                 "src/秘密.py",
		StartColumn:          9,
		EndColumn:            12,
		StartColumnCodepoint: 5,
		EndColumnCodepoint:   8,
	})
	require.Len(t, locations, 1)
	region := locations[0].PhysicalLocation.Region
	assert.Equal(t, 5, region.StartColumn)
	assert.Equal(t, 8, region.EndColumn)
}

func TestSarifResultLevel(t *testing.T) {
	findings := []Finding{{RuleID: "test-rule", File: "auth.py", StartLine: 1, EndLine: 1, StartColumn: 1, EndColumn: 2}}
	for _, test := range []struct {
		name      string
		level     string
		wantLevel bool
	}{
		{name: "historical omission", level: "", wantLevel: false},
		{name: "explicit warning", level: "warning", wantLevel: true},
		{name: "explicit note", level: "note", wantLevel: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			tmpfile, err := os.Create(filepath.Join(t.TempDir(), "report.sarif"))
			require.NoError(t, err)
			defer tmpfile.Close()
			reporter := SarifReporter{Level: test.level}
			require.NoError(t, reporter.Write(tmpfile, findings))
			raw, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			var document struct {
				Runs []struct {
					Results []struct {
						Level string `json:"level"`
					} `json:"results"`
				} `json:"runs"`
			}
			require.NoError(t, json.Unmarshal(raw, &document))
			require.Len(t, document.Runs, 1)
			require.Len(t, document.Runs[0].Results, 1)
			if test.wantLevel {
				assert.Equal(t, test.level, document.Runs[0].Results[0].Level)
			} else {
				assert.Empty(t, document.Runs[0].Results[0].Level)
			}
		})
	}
}

func TestWriteSarif(t *testing.T) {
	tests := []struct {
		findings       []Finding
		testReportName string
		expected       string
		wantEmpty      bool
		cfgName        string
	}{
		{
			cfgName:        "simple",
			testReportName: "simple",
			expected:       filepath.Join(expectPath, "report", "sarif_simple.sarif"),
			findings: []Finding{
				{

					RuleID:      "test-rule",
					Description: "A test rule",
					Match:       "line containing secret",
					Secret:      "a secret",
					StartLine:   1,
					EndLine:     2,
					StartColumn: 1,
					EndColumn:   2,
					Message:     "opps",
					File:        "auth.py",
					Commit:      "0000000000000000",
					Author:      "John Doe",
					Email:       "johndoe@gmail.com",
					Date:        "10-19-2003",
					Tags:        []string{"tag1", "tag2", "tag3"},
				},
			}},
	}

	for _, test := range tests {
		t.Run(test.cfgName, func(t *testing.T) {
			tmpfile, err := os.Create(filepath.Join(t.TempDir(), test.testReportName+".json"))
			require.NoError(t, err)
			defer tmpfile.Close()

			reporter := SarifReporter{
				OrderedRules: []config.Rule{
					{
						RuleID:      "aws-access-key",
						Description: "AWS Access Key",
					},
					{
						RuleID:      "pypi",
						Description: "PyPI upload token",
					},
				},
			}
			err = reporter.Write(tmpfile, test.findings)
			require.NoError(t, err)
			assert.FileExists(t, tmpfile.Name())

			got, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)

			if test.wantEmpty {
				assert.Empty(t, got)
				return
			}

			want, err := os.ReadFile(test.expected)
			require.NoError(t, err)

			wantStr := lineEndingReplacer.Replace(string(want))
			gotStr := lineEndingReplacer.Replace(string(got))
			assert.Equal(t, wantStr, gotStr)
		})
	}
}
