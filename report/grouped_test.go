package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupedReporterUsesKeyedIdentifiers(t *testing.T) {
	tmp, err := os.Create(filepath.Join(t.TempDir(), "grouped.json"))
	require.NoError(t, err)
	defer tmp.Close()

	reporter := GroupedReporter{Key: []byte("0123456789abcdef0123456789abcdef")}
	findings := []Finding{
		{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Commit: "c1", Fingerprint: "c1:a.txt:alpha:1", Secret: "ALPHA_12345678"},
		{RuleID: "alpha", File: "b.txt", StartLine: 2, EndLine: 2, Commit: "c2", Fingerprint: "c2:b.txt:alpha:2", Secret: "ALPHA_12345678"},
	}
	require.NoError(t, reporter.Write(tmp, findings))
	raw, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"count": 2`)
	assert.Contains(t, string(raw), `hmac-sha256:`)
	assert.NotContains(t, string(raw), "ALPHA_12345678")
	assert.NotContains(t, string(raw), "12345678")
}

func TestGroupedReporterUsesPreRedactionSecret(t *testing.T) {
	tmp, err := os.Create(filepath.Join(t.TempDir(), "grouped.json"))
	require.NoError(t, err)
	defer tmp.Close()
	finding := Finding{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "REDACTED", Fingerprint: "a.txt:alpha:1"}
	finding.SetGroupingSecret("ALPHA_12345678")
	reporter := GroupedReporter{Key: []byte("0123456789abcdef0123456789abcdef")}
	require.NoError(t, reporter.Write(tmp, []Finding{finding}))
	raw, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"count": 1`)
	assert.NotContains(t, string(raw), "ALPHA_12345678")
}
