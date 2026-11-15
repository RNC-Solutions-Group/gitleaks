package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/report"
)

func TestGroupedReportPublicationFailureLeavesNoDestination(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "grouped.json")
	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", "", "")
	require.NoError(t, err)
	out.rename = func(string, string) error { return errors.New("injected grouped publish failure") }

	_, err = out.write(&report.GroupedReporter{Key: []byte("0123456789abcdef0123456789abcdef")}, []report.Finding{
		{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "ALPHA_12345678", Fingerprint: "a.txt:alpha:1"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "injected grouped publish failure")
	assert.NoFileExists(t, dest)
}

func TestGroupedReportRenderFailureLeavesNoDestination(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "grouped.json")
	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", "", "")
	require.NoError(t, err)

	_, err = out.write(&report.GroupedReporter{Key: []byte("short")}, []report.Finding{
		{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "ALPHA_12345678", Fingerprint: "a.txt:alpha:1"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least 32 bytes")
	assert.NoFileExists(t, dest)
}

func TestGroupedReportExistingDestinationPreservedOnRenderFailure(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "grouped.json")
	require.NoError(t, os.WriteFile(dest, []byte("SENTINEL"), 0o600))
	out, err := newReportOutput(dest, reportOverwrite, "", nil, "", "", "")
	require.NoError(t, err)

	_, err = out.write(&report.GroupedReporter{Key: []byte("short")}, []report.Finding{
		{RuleID: "alpha", File: "a.txt", StartLine: 1, EndLine: 1, Secret: "ALPHA_12345678", Fingerprint: "a.txt:alpha:1"},
	})
	require.Error(t, err)
	assert.Equal(t, "SENTINEL", string(mustReadFile(t, dest)))
}
