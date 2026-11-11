package detect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const readerSizeTestSecret = `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`

func TestDetectReaderHonorsMaxTargetSize(t *testing.T) {
	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	detector.MaxTargetMegaBytes = 1

	input := strings.Repeat("A", 1_000_001) + readerSizeTestSecret + "\n"
	findings, err := detector.DetectReader(strings.NewReader(input), 64)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestStreamDetectReaderHonorsMaxTargetSize(t *testing.T) {
	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	detector.MaxTargetMegaBytes = 1

	input := strings.Repeat("A", 1_000_001) + readerSizeTestSecret + "\n"
	findingsCh, errCh := detector.StreamDetectReader(strings.NewReader(input), 64)
	for range findingsCh {
		t.Fatal("oversized stream must be excluded before any finding is yielded")
	}
	require.NoError(t, <-errCh)
}

func TestDetectReaderScansAtLimit(t *testing.T) {
	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	detector.MaxTargetMegaBytes = 1

	input := strings.Repeat("A", 1_000_000-len(readerSizeTestSecret)) + readerSizeTestSecret
	findings, err := detector.DetectReader(strings.NewReader(input), 64)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, "AKIAABCDEFGHIJKLMNOP", findings[0].Secret)
}
