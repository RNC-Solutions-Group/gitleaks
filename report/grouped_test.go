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

func TestGroupedReporterKeyRotationChangesIDs(t *testing.T) {
	tmp, err := os.Create(filepath.Join(t.TempDir(), "a.json"))
	require.NoError(t, err)
	defer tmp.Close()
	finding := Finding{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "ALPHA_12345678", Fingerprint: "a.txt:alpha:1"}
	first := GroupedReporter{Key: []byte("0123456789abcdef0123456789abcdef")}
	second := GroupedReporter{Key: []byte("abcdef0123456789abcdef0123456789")}
	require.NoError(t, first.Write(tmp, []Finding{finding}))
	rawFirst, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	require.NoError(t, tmp.Truncate(0))
	_, err = tmp.Seek(0, 0)
	require.NoError(t, err)
	require.NoError(t, second.Write(tmp, []Finding{finding}))
	rawSecond, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.NotEqual(t, string(rawFirst), string(rawSecond))
}

func TestGroupedReporterIsDeterministicAndSkipsEmptySecrets(t *testing.T) {
	tmp, err := os.Create(filepath.Join(t.TempDir(), "a.json"))
	require.NoError(t, err)
	defer tmp.Close()
	key := []byte("0123456789abcdef0123456789abcdef")
	finding := Finding{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "a", Fingerprint: "a.txt:alpha:1"}
	reporter := GroupedReporter{Key: key}
	require.NoError(t, reporter.Write(tmp, []Finding{finding, finding, {RuleID: "empty", File: "empty.txt", Secret: ""}}))
	raw, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"count": 2`)
	assert.NotContains(t, string(raw), "empty.txt")
	require.NoError(t, tmp.Truncate(0))
	_, err = tmp.Seek(0, 0)
	require.NoError(t, err)
	require.NoError(t, reporter.Write(tmp, []Finding{finding, finding}))
	rawAgain, err := os.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Equal(t, string(raw), string(rawAgain))
}
