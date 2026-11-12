// Package informative collects optional, non-sensitive scan coverage
// information for the operator-facing informative report.
//
// This package never observes finding content. It records only target paths,
// exclusion reasons, and coarse metadata that the scanner already has while
// making its normal exclusion decisions.
package informative

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

// SchemaVersion is the version of the informative JSON document. The schema is
// additive within a major version; consumers must ignore unknown fields.
const SchemaVersion = 1

// MaxExclusionDetails bounds retained per-target detail records. Aggregate
// counts remain exact when the cap is reached, and the document states how many
// records were omitted. Retaining the lexicographically smallest keys keeps the
// bounded output deterministic across concurrent workers.
const MaxExclusionDetails = 10_000

// Reason is a machine-readable exclusion reason.
type Reason string

const (
	// ReasonSizeLimit means a file, archive member, or decompressed stream
	// exceeded the configured size policy.
	ReasonSizeLimit Reason = "size_limit"
	// ReasonBinaryContent means content was classified as binary.
	ReasonBinaryContent Reason = "binary_content"
	// ReasonArchiveDepth means archive traversal was disabled or the maximum
	// archive depth was reached.
	ReasonArchiveDepth Reason = "archive_depth"
	// ReasonPathAllowlist means a target path matched a path-level allowlist.
	ReasonPathAllowlist Reason = "path_allowlist"
	// ReasonCommitAllowlist means a commit matched a commit-level allowlist.
	ReasonCommitAllowlist Reason = "commit_allowlist"
	// ReasonSymlinkNotFollowed means a symlink was not followed under the
	// active traversal policy.
	ReasonSymlinkNotFollowed Reason = "symlink_not_followed"
	// ReasonEmptyFile means an empty target was skipped.
	ReasonEmptyFile Reason = "empty_file"
	// ReasonExcludedInput means a target is a report/configuration input that
	// must not be scanned as source content.
	ReasonExcludedInput Reason = "excluded_input"
)

// Kind identifies the type of target in an exclusion record.
type Kind string

const (
	KindFile          Kind = "file"
	KindArchiveMember Kind = "archive_member"
	KindDirectory     Kind = "directory"
	KindStdin         Kind = "stdin"
	KindStream        Kind = "stream"
	KindGitDiff       Kind = "git_diff"
	KindGitBlob       Kind = "git_blob"
)

var reasonOrder = []Reason{
	ReasonSizeLimit,
	ReasonBinaryContent,
	ReasonArchiveDepth,
	ReasonPathAllowlist,
	ReasonCommitAllowlist,
	ReasonSymlinkNotFollowed,
	ReasonEmptyFile,
	ReasonExcludedInput,
}

var reasonLabels = map[Reason]string{
	ReasonSizeLimit:          "Size limit",
	ReasonBinaryContent:      "Binary content",
	ReasonArchiveDepth:       "Archive depth",
	ReasonPathAllowlist:      "Path allowlist",
	ReasonCommitAllowlist:    "Commit allowlist",
	ReasonSymlinkNotFollowed: "Symlink not followed",
	ReasonEmptyFile:          "Empty file",
	ReasonExcludedInput:      "Excluded input",
}

// ReasonLabel returns the console label for a reason. Unknown values are
// preserved so newly added reasons remain visible to older builds.
func ReasonLabel(reason Reason) string {
	if label, ok := reasonLabels[reason]; ok {
		return label
	}
	if reason == "" {
		return "Unknown"
	}
	return string(reason)
}

// ExclusionOptions carries metadata that is known at the exclusion site.
// SizeBytes overrides the target's recorded size when the effective size is
// discovered later (for example, after bounded decompression).
type ExclusionOptions struct {
	SizeBytes  *int64
	LimitBytes *int64
	// OmitSize suppresses the target's recorded size when the exclusion site
	// only observed a lower bound (for example, a bounded stream read).
	OmitSize bool
	Detail   string
}

// Exclusion is one intentionally excluded target. It never contains source
// content, finding text, credentials, or Git diagnostics.
type Exclusion struct {
	Path       string `json:"path,omitempty"`
	Kind       Kind   `json:"kind"`
	Reason     Reason `json:"reason"`
	Commit     string `json:"commit,omitempty"`
	SizeBytes  *int64 `json:"size_bytes,omitempty"`
	LimitBytes *int64 `json:"limit_bytes,omitempty"`
	Detail     string `json:"detail,omitempty"`
}

// Scan describes the completed scan that produced the document.
type Scan struct {
	SourceType    string `json:"source_type"`
	Completed     bool   `json:"completed"`
	PrimaryReport string `json:"primary_report,omitempty"`
}

// Summary contains exact aggregate counts. InspectedTargets is the count of
// observed targets that were not intentionally excluded.
type Summary struct {
	ObservedTargets    int            `json:"observed_targets"`
	InspectedTargets   int            `json:"inspected_targets"`
	ExcludedTargets    int            `json:"excluded_targets"`
	ExclusionsByReason map[Reason]int `json:"exclusions_by_reason"`
}

// Document is the versioned informative JSON document.
type Document struct {
	SchemaVersion     int         `json:"schema_version"`
	Scan              Scan        `json:"scan"`
	Summary           Summary     `json:"summary"`
	Exclusions        []Exclusion `json:"exclusions"`
	DetailsTruncated  bool        `json:"details_truncated,omitempty"`
	OmittedExclusions int         `json:"omitted_exclusions,omitempty"`
}

// Collector accumulates exclusion decisions for one scan. It is safe for
// concurrent use. A nil Collector is disabled and all methods are nil-safe.
type Collector struct {
	sourceType string

	observed atomic.Int64

	mu       sync.Mutex
	excluded int
	byReason map[Reason]int
	details  detailHeap
	omitted  int
}

// New creates an enabled collector for sourceType.
func New(sourceType string) *Collector {
	return &Collector{
		sourceType: sourceType,
		byReason:   make(map[Reason]int),
	}
}

// SetSourceType records the scan source type. It is normally set once by the
// command layer before scanning starts.
func (c *Collector) SetSourceType(sourceType string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.sourceType = sourceType
	c.mu.Unlock()
}

// Target is one observed scan target. The source owns exactly one Target per
// filesystem entry, archive member, Git diff file, or stdin stream. Exclusions
// are recorded once per target; a later exclusion attempt is ignored.
type Target struct {
	collector *Collector
	path      string
	kind      Kind
	commit    string
	size      int64
	sizeKnown bool

	// excluded is protected by collector.mu.
	excluded bool
}

// Track observes a target and returns its exclusion handle. path and commit are
// reporting metadata only; callers must not pass content or credentials.
func (c *Collector) Track(kind Kind, path, commit string, size int64, sizeKnown bool) *Target {
	if c == nil {
		return nil
	}
	c.observed.Add(1)
	return &Target{
		collector: c,
		path:      normalizePath(path),
		kind:      kind,
		commit:    commit,
		size:      size,
		sizeKnown: sizeKnown,
	}
}

// Exclude records the first exclusion reason for the target.
func (t *Target) Exclude(reason Reason, options ExclusionOptions) {
	if t == nil || t.collector == nil {
		return
	}
	c := t.collector

	c.mu.Lock()
	defer c.mu.Unlock()
	if t.excluded {
		return
	}
	t.excluded = true
	c.excluded++
	if c.byReason == nil {
		c.byReason = make(map[Reason]int)
	}
	c.byReason[reason]++

	exclusion := Exclusion{
		Path:   t.path,
		Kind:   t.kind,
		Reason: reason,
		Commit: t.commit,
		Detail: options.Detail,
	}
	if options.SizeBytes != nil {
		size := *options.SizeBytes
		exclusion.SizeBytes = &size
	} else if t.sizeKnown && !options.OmitSize {
		size := t.size
		exclusion.SizeBytes = &size
	}
	if options.LimitBytes != nil {
		limit := *options.LimitBytes
		exclusion.LimitBytes = &limit
	}
	if c.details.add(exclusion) {
		c.omitted++
	}
}

// SetSize updates the reporting-only size metadata for a target before it is
// excluded. It is used when a reader's size becomes known after observation.
func (t *Target) SetSize(size int64, known bool) {
	if t == nil || t.collector == nil {
		return
	}
	t.collector.mu.Lock()
	t.size = size
	t.sizeKnown = known
	t.collector.mu.Unlock()
}

// Document returns the completed informative document. It must be called after
// all source and detector workers have finished.
func (c *Collector) Document(primaryReport string, completed bool) Document {
	if c == nil {
		return Document{
			SchemaVersion: SchemaVersion,
			Exclusions:    []Exclusion{},
			Summary: Summary{
				ExclusionsByReason: map[Reason]int{},
			},
		}
	}

	observed := int(c.observed.Load())

	c.mu.Lock()
	excluded := c.excluded
	byReason := make(map[Reason]int, len(c.byReason))
	for reason, count := range c.byReason {
		byReason[reason] = count
	}
	exclusions := c.details.sorted()
	omitted := c.omitted
	sourceType := c.sourceType
	c.mu.Unlock()

	inspected := observed - excluded
	if inspected < 0 {
		// Defensive only: an exclusion is always associated with an observed
		// target, so this cannot occur under the documented collection model.
		inspected = 0
	}
	if exclusions == nil {
		exclusions = []Exclusion{}
	}

	return Document{
		SchemaVersion: SchemaVersion,
		Scan: Scan{
			SourceType:    sourceType,
			Completed:     completed,
			PrimaryReport: primaryReport,
		},
		Summary: Summary{
			ObservedTargets:    observed,
			InspectedTargets:   inspected,
			ExcludedTargets:    excluded,
			ExclusionsByReason: byReason,
		},
		Exclusions:        exclusions,
		DetailsTruncated:  omitted > 0,
		OmittedExclusions: omitted,
	}
}

// Marshal renders the versioned JSON document with stable indentation and a
// trailing newline.
func Marshal(document Document) ([]byte, error) {
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not encode informative report: %w", err)
	}
	return append(data, '\n'), nil
}

// DefaultConsoleDetails is the number of target detail lines shown in the
// console summary. The complete aggregate counts are always shown.
const DefaultConsoleDetails = 20

// WriteConsole writes a bounded, human-readable summary to w. The caller must
// pass stderr when primary report data is written to stdout.
func WriteConsole(w io.Writer, document Document, maxDetails int) error {
	if maxDetails <= 0 {
		maxDetails = DefaultConsoleDetails
	}

	var b strings.Builder
	b.WriteString("INFORMATIVE SCAN SUMMARY\n")
	fmt.Fprintf(&b, "  %-27s %d\n", "Observed targets:", document.Summary.ObservedTargets)
	fmt.Fprintf(&b, "  %-27s %d\n", "Inspected targets:", document.Summary.InspectedTargets)
	fmt.Fprintf(&b, "  %-27s %d\n", "Excluded targets:", document.Summary.ExcludedTargets)
	if document.Scan.PrimaryReport != "" {
		fmt.Fprintf(&b, "  %-27s %s\n", "Primary report:", SanitizeConsoleText(document.Scan.PrimaryReport, 160))
	}

	if document.Summary.ExcludedTargets == 0 {
		b.WriteString("  No excluded targets observed.\n")
		_, err := io.WriteString(w, b.String())
		return err
	}

	for _, reason := range reasonOrder {
		if count := document.Summary.ExclusionsByReason[reason]; count > 0 {
			fmt.Fprintf(&b, "  %-27s %d\n", ReasonLabel(reason)+":", count)
		}
	}

	if len(document.Exclusions) > 0 {
		b.WriteString("\nEXCLUDED TARGETS\n")
		shown := len(document.Exclusions)
		if shown > maxDetails {
			shown = maxDetails
		}
		for _, exclusion := range document.Exclusions[:shown] {
			path := exclusion.Path
			if path == "" {
				path = "(unnamed target)"
			}
			fmt.Fprintf(&b, "  %-54s %s\n", SanitizeConsoleText(path, 54), exclusion.Reason)
		}
		omitted := document.OmittedExclusions + len(document.Exclusions) - shown
		if omitted > 0 {
			fmt.Fprintf(&b, "  ... %d additional excluded target(s) omitted\n", omitted)
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	return filepath.ToSlash(path)
}

// SanitizeConsoleText keeps a path on one line and bounds its display width.
// Paths are operator metadata, not content, but newline injection would make
// the console summary ambiguous.
func SanitizeConsoleText(value string, maxRunes int) string {
	var b strings.Builder
	for _, r := range value {
		switch r {
		case '\n', '\r', '\t':
			b.WriteByte(' ')
		default:
			if r < 0x20 || r == 0x7f {
				b.WriteByte('?')
			} else {
				b.WriteRune(r)
			}
		}
	}
	value = b.String()
	if maxRunes > 0 && len([]rune(value)) > maxRunes {
		runes := []rune(value)
		if maxRunes <= 3 {
			return string(runes[:maxRunes])
		}
		value = string(runes[:maxRunes-3]) + "..."
	}
	return value
}

type detailItem struct {
	exclusion Exclusion
	key       string
}

// detailHeap is a max-heap ordered by the deterministic detail key. Keeping the
// smallest MaxExclusionDetails keys avoids scheduling-dependent truncation.
type detailHeap []detailItem

func (h detailHeap) Len() int           { return len(h) }
func (h detailHeap) Less(i, j int) bool { return h[i].key > h[j].key }
func (h detailHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *detailHeap) Push(x any)        { *h = append(*h, x.(detailItem)) }
func (h *detailHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// add returns true when one exclusion detail is omitted from the retained set.
func (h *detailHeap) add(exclusion Exclusion) (omitted bool) {
	item := detailItem{exclusion: exclusion, key: exclusionKey(exclusion)}
	if h.Len() < MaxExclusionDetails {
		heap.Push(h, item)
		return false
	}
	// The heap root is the largest retained key. Replace it only with a
	// smaller key so the retained set is the deterministic minimum. Either way
	// one detail record is no longer retained.
	if item.key < (*h)[0].key {
		heap.Pop(h)
		heap.Push(h, item)
	}
	return true
}

func (h detailHeap) sorted() []Exclusion {
	if len(h) == 0 {
		return []Exclusion{}
	}
	items := append([]detailItem(nil), h...)
	sort.Slice(items, func(i, j int) bool { return items[i].key < items[j].key })
	exclusions := make([]Exclusion, len(items))
	for i := range items {
		exclusions[i] = items[i].exclusion
	}
	return exclusions
}

func exclusionKey(exclusion Exclusion) string {
	return strings.Join([]string{
		exclusion.Path,
		exclusion.Commit,
		string(exclusion.Kind),
		string(exclusion.Reason),
		exclusion.Detail,
	}, "\x00")
}
