package informative

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorCountsAndReasons(t *testing.T) {
	collector := New("directory")
	inspected := collector.Track(KindFile, "/scan/readme.txt", "", 10, true)
	_ = inspected
	empty := collector.Track(KindFile, "/scan/empty.txt", "", 0, true)
	empty.Exclude(ReasonEmptyFile, ExclusionOptions{Detail: "empty file"})
	large := collector.Track(KindFile, "/scan/large.txt", "", 2_000_000, true)
	large.Exclude(ReasonSizeLimit, ExclusionOptions{
		LimitBytes: int64Pointer(1_000_000),
		Detail:     "maximum file size 1000000 bytes",
	})

	document := collector.Document("report.json", true)
	assert.Equal(t, SchemaVersion, document.SchemaVersion)
	assert.Equal(t, "directory", document.Scan.SourceType)
	assert.True(t, document.Scan.Completed)
	assert.Equal(t, "report.json", document.Scan.PrimaryReport)
	assert.Equal(t, 3, document.Summary.ObservedTargets)
	assert.Equal(t, 1, document.Summary.InspectedTargets)
	assert.Equal(t, 2, document.Summary.ExcludedTargets)
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[ReasonEmptyFile])
	assert.Equal(t, 1, document.Summary.ExclusionsByReason[ReasonSizeLimit])
	require.Len(t, document.Exclusions, 2)
	assert.Equal(t, "/scan/empty.txt", document.Exclusions[0].Path)
	assert.Equal(t, ReasonEmptyFile, document.Exclusions[0].Reason)
	assert.Equal(t, "/scan/large.txt", document.Exclusions[1].Path)
	require.NotNil(t, document.Exclusions[1].SizeBytes)
	assert.EqualValues(t, 2_000_000, *document.Exclusions[1].SizeBytes)
	require.NotNil(t, document.Exclusions[1].LimitBytes)
	assert.EqualValues(t, 1_000_000, *document.Exclusions[1].LimitBytes)
	assert.False(t, document.DetailsTruncated)
	assert.Zero(t, document.OmittedExclusions)
}

func TestExcludeIsFirstReasonOnly(t *testing.T) {
	collector := New("directory")
	target := collector.Track(KindFile, "file.bin", "", 1, true)
	target.Exclude(ReasonBinaryContent, ExclusionOptions{Detail: "first"})
	target.Exclude(ReasonSizeLimit, ExclusionOptions{Detail: "second"})

	document := collector.Document("", true)
	require.Len(t, document.Exclusions, 1)
	assert.Equal(t, ReasonBinaryContent, document.Exclusions[0].Reason)
	assert.Equal(t, "first", document.Exclusions[0].Detail)
	assert.Equal(t, 1, document.Summary.ExcludedTargets)
}

func TestCollectorTrackIsConcurrentSafe(t *testing.T) {
	collector := New("git")
	const workers = 32
	const perWorker = 100

	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for index := 0; index < perWorker; index++ {
				path := fmt.Sprintf("commit/%02d/file-%03d.txt", worker, index)
				target := collector.Track(KindGitDiff, path, "abc1234", 0, false)
				if index%2 == 0 {
					target.Exclude(ReasonPathAllowlist, ExclusionOptions{Detail: "global path allowlist"})
				}
			}
		}(worker)
	}
	wg.Wait()

	document := collector.Document("", true)
	assert.Equal(t, workers*perWorker, document.Summary.ObservedTargets)
	assert.Equal(t, workers*perWorker/2, document.Summary.ExcludedTargets)
	assert.Equal(t, workers*perWorker/2, document.Summary.InspectedTargets)
	assert.Equal(t, workers*perWorker/2, document.Summary.ExclusionsByReason[ReasonPathAllowlist])
}

func TestDetailsAreBoundedAndDeterministic(t *testing.T) {
	collector := New("directory")
	for index := MaxExclusionDetails + 5; index >= 0; index-- {
		target := collector.Track(KindFile, fmt.Sprintf("file-%08d.txt", index), "", 1, true)
		target.Exclude(ReasonSizeLimit, ExclusionOptions{Detail: "limit"})
	}

	document := collector.Document("", true)
	require.Len(t, document.Exclusions, MaxExclusionDetails)
	assert.True(t, document.DetailsTruncated)
	assert.Equal(t, 6, document.OmittedExclusions)
	assert.Equal(t, MaxExclusionDetails+6, document.Summary.ExcludedTargets)
	assert.Equal(t, "file-00000000.txt", document.Exclusions[0].Path)
	assert.Equal(t, fmt.Sprintf("file-%08d.txt", MaxExclusionDetails-1), document.Exclusions[len(document.Exclusions)-1].Path)
}

func TestMarshalDocument(t *testing.T) {
	collector := New("stdin")
	target := collector.Track(KindStdin, "", "", 0, false)
	target.Exclude(ReasonSizeLimit, ExclusionOptions{
		SizeBytes:  int64Pointer(1_000_001),
		LimitBytes: int64Pointer(1_000_000),
		Detail:     "maximum file size 1000000 bytes",
	})

	data, err := Marshal(collector.Document("-", true))
	require.NoError(t, err)
	assert.True(t, bytes.HasSuffix(data, []byte("\n")))

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.EqualValues(t, SchemaVersion, decoded["schema_version"])
	summary := decoded["summary"].(map[string]any)
	assert.EqualValues(t, 1, summary["observed_targets"])
	assert.EqualValues(t, 0, summary["inspected_targets"])
	assert.EqualValues(t, 1, summary["excluded_targets"])
	exclusions := decoded["exclusions"].([]any)
	require.Len(t, exclusions, 1)
	exclusion := exclusions[0].(map[string]any)
	assert.Equal(t, "stdin", exclusion["kind"])
	assert.Equal(t, "size_limit", exclusion["reason"])
}

func TestDocumentGeneratedAtIsUTCStableAndConsistent(t *testing.T) {
	fixed := time.Date(2026, 11, 12, 20, 4, 12, 987_000_000, time.FixedZone("EST", -5*60*60))
	nowCalls := 0
	collector := New("directory")
	collector.now = func() time.Time {
		nowCalls++
		return fixed
	}

	first := collector.Document("report.json", true)
	second := collector.Document("report.json", true)
	assert.Equal(t, 1, nowCalls, "the timestamp must be captured once per collector")
	assert.Equal(t, "2026-11-13T01:04:12Z", first.Scan.GeneratedAt)
	assert.Equal(t, first.Scan.GeneratedAt, second.Scan.GeneratedAt)
	_, err := time.Parse(time.RFC3339, first.Scan.GeneratedAt)
	require.NoError(t, err)

	var console bytes.Buffer
	require.NoError(t, WriteConsole(&console, first, 20))
	assert.Contains(t, console.String(), "Generated at (UTC):")
	assert.Contains(t, console.String(), first.Scan.GeneratedAt)

	data, err := Marshal(first)
	require.NoError(t, err)
	var decoded struct {
		Scan struct {
			GeneratedAt string `json:"generated_at"`
		} `json:"scan"`
	}
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, first.Scan.GeneratedAt, decoded.Scan.GeneratedAt)
}

func TestDocumentGeneratedAtIsAbsentWhenCollectionIsDisabled(t *testing.T) {
	var collector *Collector
	document := collector.Document("", true)
	assert.Empty(t, document.Scan.GeneratedAt)
	data, err := Marshal(document)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "generated_at")
}

func TestConcurrentCollectorsKeepIndependentGeneratedAt(t *testing.T) {
	const workers = 16
	base := time.Date(2026, 11, 12, 20, 4, 12, 0, time.UTC)
	collectors := make([]*Collector, workers)
	for index := range collectors {
		collectors[index] = New("directory")
		offset := time.Duration(index) * time.Second
		collectors[index].now = func() time.Time { return base.Add(offset) }
	}

	var wg sync.WaitGroup
	for index, collector := range collectors {
		wg.Add(1)
		go func(index int, collector *Collector) {
			defer wg.Done()
			document := collector.Document("", true)
			assert.Equal(t, base.Add(time.Duration(index)*time.Second).Format(time.RFC3339), document.Scan.GeneratedAt)
		}(index, collector)
	}
	wg.Wait()
}

func TestWriteConsole(t *testing.T) {
	collector := New("directory")
	collector.Track(KindFile, "readme.txt", "", 10, true)
	empty := collector.Track(KindFile, "empty.txt", "", 0, true)
	empty.Exclude(ReasonEmptyFile, ExclusionOptions{Detail: "empty file"})
	binary := collector.Track(KindFile, "image.bin", "", 12, true)
	binary.Exclude(ReasonBinaryContent, ExclusionOptions{Detail: "MIME type application/octet-stream"})

	var output bytes.Buffer
	require.NoError(t, WriteConsole(&output, collector.Document("", true), 20))
	text := output.String()
	assert.Contains(t, text, "INFORMATIVE SCAN SUMMARY")
	assert.Contains(t, text, "Observed targets:")
	assert.Contains(t, text, "Excluded targets:           2")
	assert.Contains(t, text, "Binary content:")
	assert.Contains(t, text, "Empty file:")
	assert.Contains(t, text, "EXCLUDED TARGETS")
	assert.Contains(t, text, "empty.txt")
	assert.Contains(t, text, "binary_content")
	assert.NotContains(t, text, "readme.txt")
}

func TestWriteConsoleNoExclusions(t *testing.T) {
	collector := New("directory")
	collector.Track(KindFile, "readme.txt", "", 10, true)

	var output bytes.Buffer
	require.NoError(t, WriteConsole(&output, collector.Document("", true), 20))
	text := output.String()
	assert.Contains(t, text, "No excluded targets observed.")
	assert.NotContains(t, text, "EXCLUDED TARGETS")
}

func TestWriteConsoleBoundsDetailsAndSanitizesPaths(t *testing.T) {
	collector := New("directory")
	for index := 0; index < 5; index++ {
		target := collector.Track(KindFile, fmt.Sprintf("file-%d\nname.txt", index), "", 1, true)
		target.Exclude(ReasonPathAllowlist, ExclusionOptions{Detail: "global path allowlist"})
	}

	var output bytes.Buffer
	require.NoError(t, WriteConsole(&output, collector.Document("", true), 2))
	text := output.String()
	assert.Contains(t, text, "file-0 name.txt")
	assert.NotContains(t, text, "file-0\nname.txt")
	assert.Contains(t, text, "3 additional excluded target(s) omitted")
}

func TestNilCollectorAndTargetAreSafe(t *testing.T) {
	var collector *Collector
	target := collector.Track(KindFile, "file", "", 1, true)
	target.Exclude(ReasonEmptyFile, ExclusionOptions{})
	document := collector.Document("", true)
	assert.Equal(t, SchemaVersion, document.SchemaVersion)
	assert.Zero(t, document.Summary.ObservedTargets)
	require.NotNil(t, document.Exclusions)
}

func TestSanitizeConsoleText(t *testing.T) {
	assert.Equal(t, "a b c", SanitizeConsoleText("a\nb\tc", 10))
	assert.Equal(t, "abc", SanitizeConsoleText("abc", 3))
	assert.Equal(t, "abcd...", SanitizeConsoleText("abcdefgh", 7))
	assert.Equal(t, "ab", SanitizeConsoleText("abcdef", 2))
	assert.Equal(t, "a?b", SanitizeConsoleText("a\x00b", 10))
}

func int64Pointer(value int64) *int64 {
	return &value
}

func TestExclusionKeyStable(t *testing.T) {
	exclusion := Exclusion{
		Path:   "a",
		Commit: "b",
		Kind:   KindFile,
		Reason: ReasonSizeLimit,
		Detail: "d",
	}
	assert.Equal(t, "a\x00b\x00file\x00size_limit\x00d", exclusionKey(exclusion))
	assert.True(t, strings.HasPrefix(exclusionKey(exclusion), "a\x00"))
}
