package detect

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func TestGlobalPathAllowlistRecordsTargetExclusion(t *testing.T) {
	allowlist := &config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`^skip\.txt$`)},
	}
	require.NoError(t, allowlist.Validate())
	collector := informative.New("directory")
	target := collector.Track(informative.KindFile, "skip.txt", "", 6, true)
	detector := NewDetector(config.Config{Allowlists: []*config.Allowlist{allowlist}})

	findings := detector.DetectContext(context.Background(), Fragment(sources.Fragment{
		FilePath:          "skip.txt",
		Raw:               "secret",
		InformativeTarget: target,
	}))
	assert.Empty(t, findings)

	document := collector.Document("", true)
	assert.Equal(t, 1, document.Summary.ObservedTargets)
	assert.Equal(t, 1, document.Summary.ExcludedTargets)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[informative.ReasonPathAllowlist])
	require.Len(t, document.Exclusions, 1)
	assert.Equal(t, "skip.txt", document.Exclusions[0].Path)
	assert.Equal(t, "global path allowlist", document.Exclusions[0].Detail)
}

func TestGlobalCommitAllowlistRecordsCommitReason(t *testing.T) {
	allowlist := &config.Allowlist{Commits: []string{"abc1234"}}
	require.NoError(t, allowlist.Validate())
	collector := informative.New("git")
	target := collector.Track(informative.KindGitDiff, "path.txt", "abc1234", 0, false)
	detector := NewDetector(config.Config{Allowlists: []*config.Allowlist{allowlist}})

	findings := detector.DetectContext(context.Background(), Fragment(sources.Fragment{
		FilePath:          "path.txt",
		CommitSHA:         "abc1234",
		Raw:               "secret",
		InformativeTarget: target,
	}))
	assert.Empty(t, findings)

	document := collector.Document("", true)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[informative.ReasonCommitAllowlist])
	require.Len(t, document.Exclusions, 1)
	assert.Equal(t, "abc1234", document.Exclusions[0].Commit)
	assert.Equal(t, "global commit allowlist", document.Exclusions[0].Detail)
}

func TestReportInputExclusionDoesNotCountAsInspected(t *testing.T) {
	collector := informative.New("directory")
	target := collector.Track(informative.KindFile, "report.json", "", 100, true)
	detector := NewDetector(config.Config{})
	detector.ReportPath = "report.json"

	findings := detector.DetectContext(context.Background(), Fragment(sources.Fragment{
		FilePath:          "report.json",
		Raw:               "{}",
		InformativeTarget: target,
	}))
	assert.Empty(t, findings)

	document := collector.Document("report.json", true)
	assert.Equal(t, 0, document.Summary.InspectedTargets)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[informative.ReasonExcludedInput])
}
