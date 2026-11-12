package detect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
)

func TestDetectReaderRecordsSizeExclusion(t *testing.T) {
	collector := informative.New("stream")
	detector := NewDetector(config.Config{})
	detector.Informative = collector
	detector.MaxTargetMegaBytes = 1

	findings, err := detector.DetectReader(strings.NewReader(strings.Repeat("A", 1_000_001)), 64)
	require.NoError(t, err)
	assert.Empty(t, findings)

	document := collector.Document("", true)
	assert.Equal(t, 1, document.Summary.ObservedTargets)
	assert.Equal(t, 1, document.Summary.ExcludedTargets)
	assert.Equal(t, 0, document.Summary.InspectedTargets)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[informative.ReasonSizeLimit])
	require.Len(t, document.Exclusions, 1)
	assert.Equal(t, informative.KindStream, document.Exclusions[0].Kind)
}

func TestStreamDetectReaderRecordsSizeExclusion(t *testing.T) {
	collector := informative.New("stream")
	detector := NewDetector(config.Config{})
	detector.Informative = collector
	detector.MaxTargetMegaBytes = 1

	findingsCh, errCh := detector.StreamDetectReader(strings.NewReader(strings.Repeat("A", 1_000_001)), 64)
	var findings []string
	for finding := range findingsCh {
		findings = append(findings, finding.RuleID)
	}
	require.NoError(t, <-errCh)
	assert.Empty(t, findings)

	document := collector.Document("", true)
	assert.Equal(t, 1, document.Summary.ObservedTargets)
	assert.Equal(t, 1, document.Summary.ExcludedTargets)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[informative.ReasonSizeLimit])
}
