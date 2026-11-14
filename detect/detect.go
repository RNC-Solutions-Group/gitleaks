package detect

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect/codec"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"github.com/fatih/semgroup"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

const (
	gitleaksAllowSignature = "gitleaks:allow"
	// SlowWarningThreshold is the amount of time to wait before logging that a file is slow.
	// This is useful for identifying problematic files and tuning the allowlist.
	SlowWarningThreshold = 5 * time.Second
)

var (
	newLineRegexp = regexp.MustCompile("\n")
)

// Detector is the main detector struct
type Detector struct {
	// Config is the configuration for the detector
	Config config.Config

	// Redact is a flag to redact findings. This is exported
	// so users using gitleaks as a library can set this flag
	// without calling `detector.Start(cmd *cobra.Command)`
	Redact uint

	// verbose is a flag to print findings
	Verbose bool

	// MaxDecodeDepths limits how many recursive decoding passes are allowed
	MaxDecodeDepth int

	// MaxArchiveDepth limits how deep the sources will explore nested archives
	MaxArchiveDepth int

	// files larger than this will be skipped
	MaxTargetMegaBytes int

	// followSymlinks is a flag to enable scanning symlink files
	FollowSymlinks bool

	// NoColor is a flag to disable color output
	NoColor bool

	// FileLinks enables opt-in local file URI links in console findings.
	// Machine-readable reports are never modified.
	FileLinks bool

	// IgnoreGitleaksAllow is a flag to ignore gitleaks:allow comments.
	IgnoreGitleaksAllow bool

	// commitMutex is to prevent concurrent access to the
	// commit map when adding commits
	commitMutex *sync.Mutex

	// commitMap is used to keep track of commits that have been scanned.
	// This is only used for logging purposes and git scans.
	commitMap map[string]bool

	// findingMutex is to prevent concurrent access to the
	// findings slice when adding findings.
	findingMutex *sync.Mutex

	// findings is a slice of report.Findings. This is the result
	// of the detector's scan which can then be used to generate a
	// report.
	findings []report.Finding

	// prefilter is a ahocorasick struct used for doing efficient string
	// matching given a set of words (keywords from the rules in the config)
	prefilter ahocorasick.Trie

	// a list of known findings that should be ignored
	baseline []report.Finding

	// path to baseline
	baselinePath string

	// gitleaksIgnore
	gitleaksIgnore map[string]struct{}
	// gitleaksIgnoreUsed records entries that suppressed a finding during this
	// run. It is protected by ignoreMutex.
	gitleaksIgnoreUsed map[string]struct{}
	ignoreMutex        *sync.Mutex

	// Sema (https://github.com/fatih/semgroup) controls the concurrency
	Sema *semgroup.Group

	// ReportUnusedIgnore enables a completed-scan diagnostic listing loaded
	// ignore entries that did not suppress a finding.
	ReportUnusedIgnore bool

	// report-related settings.
	ReportPath string
	Reporter   report.Reporter
	// ReportWriter, when set, validates and writes the report for a completed
	// scan, returning the actual output filename. The CLI sets it so
	// destination policy and atomic publication stay in the command layer.
	ReportWriter func([]report.Finding) (string, error)

	// Informative is the optional per-scan coverage collector. Nil disables
	// informative reporting and imposes no collection work.
	Informative *informative.Collector
	// InformativeReportPath is the companion file destination. It prevents a
	// previous informative report from being scanned as source input.
	InformativeReportPath string

	TotalBytes atomic.Uint64
}

// Fragment is an alias for sources.Fragment for backwards compatibility
//
// Deprecated: This will be replaced with sources.Fragment in v9
type Fragment sources.Fragment

// fragmentPaths holds the path representations used for matching. The original
// reported path is always evaluated first. Directory scans add the
// scan-root-relative form, plus its native-separator spelling on Windows so
// existing backslash allowlist patterns keep working.
type fragmentPaths struct {
	file               string
	windows            string
	rootRelative       string
	rootRelativeNative string
}

func newFragmentPaths(fragment Fragment) fragmentPaths {
	paths := fragmentPaths{
		file:    fragment.FilePath,
		windows: fragment.WindowsFilePath,
	}
	rel := sources.RootRelativePath(fragment.Root, fragment.FilePath)
	if rel == "" || rel == fragment.FilePath {
		return paths
	}
	paths.rootRelative = rel
	if native := filepath.FromSlash(rel); native != rel {
		paths.rootRelativeNative = native
	}
	return paths
}

func (p fragmentPaths) allowed(a *config.Allowlist) bool {
	if a == nil {
		return false
	}
	if a.PathAllowed(p.file) {
		return true
	}
	if p.windows != "" && a.PathAllowed(p.windows) {
		return true
	}
	if p.rootRelative != "" && a.PathAllowed(p.rootRelative) {
		return true
	}
	return p.rootRelativeNative != "" && a.PathAllowed(p.rootRelativeNative)
}

// NewDetector creates a new detector with the given config
func NewDetector(cfg config.Config) *Detector {
	return NewDetectorContext(context.Background(), cfg)
}

// NewDetectorContext is the same as NewDetector but supports passing in a
// context to use for timeouts
func NewDetectorContext(ctx context.Context, cfg config.Config) *Detector {
	return &Detector{
		commitMap:          make(map[string]bool),
		gitleaksIgnore:     make(map[string]struct{}),
		gitleaksIgnoreUsed: make(map[string]struct{}),
		ignoreMutex:        &sync.Mutex{},
		findingMutex:       &sync.Mutex{},
		commitMutex:        &sync.Mutex{},
		findings:           make([]report.Finding, 0),
		Config:             cfg,
		prefilter:          *ahocorasick.NewTrieBuilder().AddStrings(maps.Keys(cfg.Keywords)).Build(),
		Sema:               semgroup.NewGroup(ctx, 40),
	}
}

// NewDetectorDefaultConfig creates a new detector with the default config
func NewDetectorDefaultConfig() (*Detector, error) {
	viper.SetConfigType("toml")
	err := viper.ReadConfig(strings.NewReader(config.DefaultConfig))
	if err != nil {
		return nil, err
	}
	var vc config.ViperConfig
	err = viper.Unmarshal(&vc)
	if err != nil {
		return nil, err
	}
	cfg, err := vc.Translate()
	if err != nil {
		return nil, err
	}
	return NewDetector(cfg), nil
}

func (d *Detector) AddGitleaksIgnore(gitleaksIgnorePath string) error {
	logging.Debug().Str("path", gitleaksIgnorePath).Msgf("found .gitleaksignore file")
	file, err := openIgnoreFile(gitleaksIgnorePath)
	if err != nil {
		return err
	}

	info, statErr := file.Stat()
	if statErr != nil {
		_ = file.Close()
		return fmt.Errorf("could not stat .gitleaksignore %q: %w", gitleaksIgnorePath, statErr)
	}
	if info.Size() > maxGitleaksIgnoreBytes {
		_ = file.Close()
		return fmt.Errorf(".gitleaksignore %q is %d bytes, exceeds limit of %d bytes", gitleaksIgnorePath, info.Size(), maxGitleaksIgnoreBytes)
	}

	entries, parseErr := parseGitleaksIgnore(file)
	closeErr := file.Close()
	if parseErr != nil {
		return fmt.Errorf("could not read .gitleaksignore %q: %w", gitleaksIgnorePath, parseErr)
	}
	if closeErr != nil {
		return fmt.Errorf("could not close .gitleaksignore %q: %w", gitleaksIgnorePath, closeErr)
	}

	// Merge only after the complete file has been parsed successfully so an
	// incomplete read cannot leave a partially loaded suppression set.
	for entry := range entries {
		d.gitleaksIgnore[entry] = struct{}{}
	}
	return nil
}

const (
	maxGitleaksIgnoreBytes   = 16 << 20
	maxGitleaksIgnoreEntries = 100_000
	maxGitleaksIgnoreLine    = 1 << 20
)

type countingReader struct {
	reader io.Reader
	bytes  int64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.bytes += int64(n)
	return n, err
}

func parseGitleaksIgnore(reader io.Reader) (map[string]struct{}, error) {
	entries := make(map[string]struct{})
	counted := &countingReader{reader: reader}
	scanner := bufio.NewScanner(counted)
	scanner.Buffer(make([]byte, 64*1024), maxGitleaksIgnoreLine)
	replacer := strings.NewReplacer("\\", "/")
	entryCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip lines that start with a comment
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entryCount++
		if entryCount > maxGitleaksIgnoreEntries {
			return nil, fmt.Errorf("contains more than %d entries", maxGitleaksIgnoreEntries)
		}
		if counted.bytes > maxGitleaksIgnoreBytes {
			return nil, fmt.Errorf("exceeds limit of %d bytes", maxGitleaksIgnoreBytes)
		}

		// Normalize the path.
		// TODO: Make this a breaking change in v9.
		s := strings.Split(line, ":")
		switch len(s) {
		case 3:
			// Global fingerprint.
			// `file:rule-id:start-line`
			s[0] = replacer.Replace(s[0])
		case 4:
			// Commit fingerprint.
			// `commit:file:rule-id:start-line`
			s[1] = replacer.Replace(s[1])
		default:
			logging.Warn().Str("fingerprint", line).Msg("Invalid .gitleaksignore entry")
		}
		entries[strings.Join(s, ":")] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if counted.bytes > maxGitleaksIgnoreBytes {
		return nil, fmt.Errorf("exceeds limit of %d bytes", maxGitleaksIgnoreBytes)
	}
	return entries, nil
}

// DetectBytes scans the given bytes and returns a list of findings
func (d *Detector) DetectBytes(content []byte) []report.Finding {
	return d.DetectString(string(content))
}

// DetectString scans the given string and returns a list of findings
func (d *Detector) DetectString(content string) []report.Finding {
	return d.Detect(Fragment{
		Raw: content,
	})
}

// DetectSource scans the given source and returns a list of findings
func (d *Detector) DetectSource(ctx context.Context, source sources.Source) ([]report.Finding, error) {
	collector := scanerr.New(10)
	err := source.Fragments(ctx, func(fragment sources.Fragment, err error) error {
		logContext := logging.With()

		if len(fragment.FilePath) > 0 {
			logContext = logContext.Str("path", fragment.FilePath)
		}

		if len(fragment.CommitSHA) > 6 {
			logContext = logContext.Str("commit", fragment.CommitSHA[:7])
			d.addCommit(fragment.CommitSHA)
		} else if len(fragment.CommitSHA) > 0 {
			logContext = logContext.Str("commit", fragment.CommitSHA)
			d.addCommit(fragment.CommitSHA)
			logger := logContext.Logger()
			logger.Warn().Msg("commit SHAs should be >= 7 characters long")
		}

		logger := logContext.Logger()

		if err != nil {
			// Log the error but keep scanning so findings from readable input
			// are preserved. The error is returned after the source completes.
			logger.Error().Err(err).Send()
			collector.Add(err)
			return nil
		}

		// both the fragment's content and path should be empty for it to be
		// considered empty at this point because of path based matches
		if len(fragment.Raw) == 0 && len(fragment.FilePath) == 0 {
			fragment.InformativeTarget.Exclude(informative.ReasonEmptyFile, informative.ExclusionOptions{
				Detail: "empty fragment",
			})
			logger.Trace().Msg("skipping empty fragment")
			return nil
		}

		var timer *time.Timer
		// Only start the timer in debug mode
		if logger.GetLevel() <= zerolog.DebugLevel {
			timer = time.AfterFunc(SlowWarningThreshold, func() {
				logger.Debug().Msgf("Taking longer than %s to inspect fragment", SlowWarningThreshold.String())
			})
		}

		for _, finding := range d.DetectContext(ctx, Fragment(fragment)) {
			if d.knownRelativeToRoot(fragment, finding) {
				continue
			}
			d.AddFinding(finding)
		}

		// Stop the timer if it was created
		if timer != nil {
			timer.Stop()
		}

		return nil
	})
	if err != nil {
		collector.Add(err)
	}

	// Sources may return before their semgroup workers have finished. Waiting
	// here both joins those workers and preserves any errors they accumulated.
	if d.Sema != nil {
		collector.Add(d.Sema.Wait())
	}

	if _, isGit := source.(*sources.Git); isGit {
		logging.Info().Msgf("%d commits scanned.", len(d.commitMap))
		logging.Debug().Msg("Note: this number might be smaller than expected due to commits with no additions")
	}

	return d.Findings(), collector.Err()
}

// knownRelativeToRoot checks .gitleaksignore fingerprints and baseline entries
// against the finding's scan-root-relative path. The reported path was already
// checked by AddFinding and remains the primary identity; this only adds the
// equivalent spelling a directory scan can produce.
func (d *Detector) knownRelativeToRoot(fragment sources.Fragment, finding report.Finding) bool {
	rel := sources.RootRelativePath(fragment.Root, finding.File)
	if rel == "" || rel == finding.File {
		return false
	}
	logger := logging.With().Str("path", finding.File).Str("rule-id", finding.RuleID).Logger()

	if finding.Commit == "" && len(d.gitleaksIgnore) > 0 {
		fingerprint := fmt.Sprintf("%s:%s:%d", rel, finding.RuleID, finding.StartLine)
		if _, ok := d.gitleaksIgnore[fingerprint]; ok {
			d.markGitleaksIgnoreUsed(fingerprint)
			logger.Debug().
				Str("fingerprint", fingerprint).
				Msg("skipping finding: global fingerprint relative to the scan root")
			return true
		}
	}

	if d.baseline != nil {
		relative := finding
		relative.File = rel
		if !IsNew(relative, d.Redact, d.baseline) {
			logger.Debug().
				Str("file", rel).
				Msg("skipping finding: baseline entry relative to the scan root")
			return true
		}
	}
	return false
}

// configOrBaselinePath identifies the configuration or baseline input under
// either the reported spelling or the scan-root-relative spelling. Baseline
// paths are stored relative to the source, while config paths are joined onto
// the source; fragment paths keep the source spelling.
func (d *Detector) configOrBaselinePath(fragment Fragment, paths fragmentPaths) (string, bool) {
	if d.baselinePath != "" &&
		(samePath(fragment.FilePath, d.baselinePath) || samePathAtRoot(fragment.Root, paths.rootRelative, d.baselinePath)) {
		return "baseline file", true
	}
	if d.Config.Path != "" &&
		(samePath(fragment.FilePath, d.Config.Path) || samePathAtRoot(fragment.Root, paths.rootRelative, d.Config.Path)) {
		return "configuration file", true
	}
	return "", false
}

func samePathAtRoot(root, rootRelative, target string) bool {
	if root == "" || rootRelative == "" || target == "" {
		return false
	}
	if filepath.ToSlash(filepath.Clean(target)) == rootRelative {
		return true
	}
	rel := sources.RootRelativePath(root, target)
	return rel != "" && rel == rootRelative
}

// Detect scans the given fragment and returns a list of findings
func (d *Detector) Detect(fragment Fragment) []report.Finding {
	return d.DetectContext(context.Background(), fragment)
}

// DetectContext is the same as Detect but supports passing in a
// context to use for timeouts
func (d *Detector) DetectContext(ctx context.Context, fragment Fragment) []report.Finding {
	if fragment.Bytes == nil {
		d.TotalBytes.Add(uint64(len(fragment.Raw)))
	}
	d.TotalBytes.Add(uint64(len(fragment.Bytes)))

	var (
		findings []report.Finding
		logger   = func() zerolog.Logger {
			l := logging.With().Str("path", fragment.FilePath)
			if fragment.CommitSHA != "" {
				l = l.Str("commit", fragment.CommitSHA)
			}
			return l.Logger()
		}()
	)

	target := fragment.InformativeTarget
	paths := newFragmentPaths(fragment)

	// check if filepath is allowed
	if fragment.FilePath != "" {
		// Do not scan the requested report destination. Scanning it would
		// re-report findings embedded in a previous report. Aggregated reports
		// use newly allocated filenames and are not existing scan inputs.
		if d.ReportPath != "" && d.ReportPath != report.StdoutReportPath &&
			sameFileOrPath(fragment.FilePath, d.ReportPath) {
			target.Exclude(informative.ReasonExcludedInput, informative.ExclusionOptions{
				Detail: "primary report destination",
			})
			logging.Trace().Msg("skipping file: matches report path")
			return findings
		}
		if d.InformativeReportPath != "" &&
			sameFileOrPath(fragment.FilePath, d.InformativeReportPath) {
			target.Exclude(informative.ReasonExcludedInput, informative.ExclusionOptions{
				Detail: "informative report destination",
			})
			logging.Trace().Msg("skipping file: matches informative report path")
			return findings
		}
		// is the path our config or baseline file?
		if detail, ok := d.configOrBaselinePath(fragment, paths); ok {
			target.Exclude(informative.ReasonExcludedInput, informative.ExclusionOptions{
				Detail: detail,
			})
			logging.Trace().Msg("skipping file: matches config or baseline path")
			return findings
		}
	}
	// check if commit or filepath is allowed.
	if isAllowed, event := checkCommitOrPathAllowed(logger, fragment, paths, d.Config.Allowlists); isAllowed {
		event.Msg("skipping file: global allowlist")
		reason := informative.ReasonPathAllowlist
		detail := "global path allowlist"
		if fragment.CommitSHA != "" && commitAllowed(fragment.CommitSHA, d.Config.Allowlists) {
			reason = informative.ReasonCommitAllowlist
			detail = "global commit allowlist"
		}
		target.Exclude(reason, informative.ExclusionOptions{Detail: detail})
		return findings
	}

	// setup variables to handle different decoding passes
	currentRaw := fragment.Raw
	encodedSegments := []*codec.EncodedSegment{}
	currentDecodeDepth := 0
	decoder := codec.NewDecoder()

ScanLoop:
	for {
		select {
		case <-ctx.Done():
			break ScanLoop
		default:
			// build keyword map for prefiltering rules
			keywords := make(map[string]bool)
			normalizedRaw := strings.ToLower(currentRaw)
			matches := d.prefilter.MatchString(normalizedRaw)
			for _, m := range matches {
				keywords[normalizedRaw[m.Pos():int(m.Pos())+len(m.Match())]] = true
			}

			for _, rule := range d.Config.Rules {
				select {
				case <-ctx.Done():
					break ScanLoop
				default:
					if len(rule.Keywords) == 0 {
						// if no keywords are associated with the rule always scan the
						// fragment using the rule
						findings = append(findings, d.detectRuleWithPaths(fragment, currentRaw, rule, encodedSegments, paths)...)
						continue
					}

					// check if keywords are in the fragment
					for _, k := range rule.Keywords {
						if _, ok := keywords[strings.ToLower(k)]; ok {
							findings = append(findings, d.detectRuleWithPaths(fragment, currentRaw, rule, encodedSegments, paths)...)
							break
						}
					}
				}
			}

			// increment the depth by 1 as we start our decoding pass
			currentDecodeDepth++

			// stop the loop if we've hit our max decoding depth
			if currentDecodeDepth > d.MaxDecodeDepth {
				break ScanLoop
			}

			// decode the currentRaw for the next pass
			currentRaw, encodedSegments = decoder.Decode(currentRaw, encodedSegments)

			// stop the loop when there's nothing else to decode
			if len(encodedSegments) == 0 {
				break ScanLoop
			}
		}
	}

	return filter(findings, d.Redact)
}

// detectRule scans the given fragment for the given rule and returns a list of findings
func (d *Detector) detectRule(fragment Fragment, currentRaw string, r config.Rule, encodedSegments []*codec.EncodedSegment) []report.Finding {
	return d.detectRuleWithPaths(fragment, currentRaw, r, encodedSegments, newFragmentPaths(fragment))
}

func (d *Detector) detectRuleWithPaths(fragment Fragment, currentRaw string, r config.Rule, encodedSegments []*codec.EncodedSegment, paths fragmentPaths) []report.Finding {
	var (
		findings []report.Finding
		logger   = func() zerolog.Logger {
			l := logging.With().Str("rule-id", r.RuleID).Str("path", fragment.FilePath)
			if fragment.CommitSHA != "" {
				l = l.Str("commit", fragment.CommitSHA)
			}
			return l.Logger()
		}()
	)

	if r.SkipReport && !fragment.InheritedFromFinding {
		return findings
	}

	// check if commit or file is allowed for this rule.
	if isAllowed, event := checkCommitOrPathAllowed(logger, fragment, paths, r.Allowlists); isAllowed {
		event.Msg("skipping file: rule allowlist")
		return findings
	}

	if r.Path != nil {
		if r.Regex == nil && len(encodedSegments) == 0 {
			// Path _only_ rule
			if r.Path.MatchString(fragment.FilePath) || (fragment.WindowsFilePath != "" && r.Path.MatchString(fragment.WindowsFilePath)) {
				finding := report.Finding{
					Commit:      fragment.CommitSHA,
					RuleID:      r.RuleID,
					Description: r.Description,
					File:        fragment.FilePath,
					SymlinkFile: fragment.SymlinkFile,
					Match:       "file detected: " + fragment.FilePath,
					Tags:        r.Tags,
				}
				if fragment.CommitInfo != nil {
					finding.Author = fragment.CommitInfo.AuthorName
					finding.Date = fragment.CommitInfo.Date
					finding.Email = fragment.CommitInfo.AuthorEmail
					finding.Link = createScmLink(fragment.CommitInfo.Remote, finding)
					finding.Message = fragment.CommitInfo.Message
				}
				return append(findings, finding)
			}
		} else {
			// if path is set _and_ a regex is set, then we need to check both
			// so if the path does not match, then we should return early and not
			// consider the regex
			if !(r.Path.MatchString(fragment.FilePath) || (fragment.WindowsFilePath != "" && r.Path.MatchString(fragment.WindowsFilePath))) {
				return findings
			}
		}
	}

	// if path only rule, skip content checks
	if r.Regex == nil {
		return findings
	}

	// if flag configure and raw data size bigger then the flag
	if d.MaxTargetMegaBytes > 0 {
		rawLength := len(currentRaw) / 1_000_000
		if rawLength > d.MaxTargetMegaBytes {
			logger.Debug().
				Int("size", rawLength).
				Int("max-size", d.MaxTargetMegaBytes).
				Msg("skipping fragment: size")
			return findings
		}
	}

	matches := r.Regex.FindAllStringIndex(currentRaw, -1)
	if len(matches) == 0 {
		return findings
	}

	// TODO profile this, probably should replace with something more efficient
	newlineIndices := newLineRegexp.FindAllStringIndex(fragment.Raw, -1)

	// use currentRaw instead of fragment.Raw since this represents the current
	// decoding pass on the text
	for _, matchIndex := range matches {
		// Extract secret from match
		secret := strings.Trim(currentRaw[matchIndex[0]:matchIndex[1]], "\r\n")

		// For any meta data from decoding
		var metaTags []string
		currentLine := ""

		// Check if the decoded portions of the segment overlap with the match
		// to see if its potentially a new match
		if len(encodedSegments) > 0 {
			segments := codec.SegmentsWithDecodedOverlap(encodedSegments, matchIndex[0], matchIndex[1])
			if len(segments) == 0 {
				// This item has already been added to a finding
				continue
			}

			matchIndex = codec.AdjustMatchIndex(segments, matchIndex)
			metaTags = append(metaTags, codec.Tags(segments)...)
			currentLine = codec.CurrentLine(segments, currentRaw)
		} else {
			// Fixes: https://github.com/gitleaks/gitleaks/issues/1352
			// removes the incorrectly following line that was detected by regex expression '\n'
			matchIndex[1] = matchIndex[0] + len(secret)
		}

		// determine location of match. Note that the location
		// in the finding will be the line/column numbers of the _match_
		// not the _secret_, which will be different if the secretGroup
		// value is set for this rule
		loc := location(newlineIndices, fragment.Raw, matchIndex)

		if matchIndex[1] > loc.endLineIndex {
			loc.endLineIndex = matchIndex[1]
		}

		finding := report.Finding{
			Commit:               fragment.CommitSHA,
			RuleID:               r.RuleID,
			Description:          r.Description,
			StartLine:            fragment.StartLine + loc.startLine,
			EndLine:              fragment.StartLine + loc.endLine,
			StartColumn:          loc.startColumn,
			EndColumn:            loc.endColumn,
			StartColumnCodepoint: loc.startColumnCodepoint,
			EndColumnCodepoint:   loc.endColumnCodepoint,
			Line:                 fragment.Raw[loc.startLineIndex:loc.endLineIndex],
			Match:                secret,
			Secret:               secret,
			File:                 fragment.FilePath,
			SymlinkFile:          fragment.SymlinkFile,
			Tags:                 append(r.Tags, metaTags...),
		}
		if fragment.CommitInfo != nil {
			finding.Author = fragment.CommitInfo.AuthorName
			finding.Date = fragment.CommitInfo.Date
			finding.Email = fragment.CommitInfo.AuthorEmail
			finding.Link = createScmLink(fragment.CommitInfo.Remote, finding)
			finding.Message = fragment.CommitInfo.Message
		}
		if !d.IgnoreGitleaksAllow && strings.Contains(finding.Line, gitleaksAllowSignature) {
			logger.Trace().
				Msg("skipping finding: 'gitleaks:allow' signature")
			continue
		}

		if currentLine == "" {
			currentLine = finding.Line
		}

		// Set the value of |secret|, if the pattern contains at least one capture group.
		// (The first element is the full match, hence we check >= 2.)
		groups := r.Regex.FindStringSubmatch(finding.Secret)
		if len(groups) >= 2 {
			if r.SecretGroup > 0 {
				if len(groups) <= r.SecretGroup {
					// Config validation should prevent this
					continue
				}
				finding.Secret = groups[r.SecretGroup]
			} else {
				// If |secretGroup| is not set, we will use the first suitable capture group.
				for _, s := range groups[1:] {
					if len(s) > 0 {
						finding.Secret = s
						break
					}
				}
			}
		}

		// Canonicalize line endings before entropy and reporting. The detector
		// uses a trimmed byte span; a CRLF checkout must produce the same
		// Match/Secret/entropy as an LF checkout even though Finding.Line
		// retains its source bytes.
		finding.Match = normalizeBaselineLineEndings(finding.Match)
		finding.Secret = normalizeBaselineLineEndings(finding.Secret)

		// check entropy
		entropy := shannonEntropy(finding.Secret)
		finding.Entropy = float32(entropy)
		if r.Entropy != 0.0 {
			// entropy is too low, skip this finding
			if entropy <= r.Entropy {
				logger.Trace().
					Float32("entropy", finding.Entropy).
					Msg("skipping finding: low entropy")
				continue
			}
		}

		// check if the result matches any of the global allowlists.
		if isAllowed, event := checkFindingAllowed(logger, finding, fragment, paths, currentLine, d.Config.Allowlists); isAllowed {
			event.Msg("skipping finding: global allowlist")
			continue
		}

		// check if the result matches any of the rule allowlists.
		if isAllowed, event := checkFindingAllowed(logger, finding, fragment, paths, currentLine, r.Allowlists); isAllowed {
			event.Msg("skipping finding: rule allowlist")
			continue
		}
		findings = append(findings, finding)
	}

	// Handle required rules (multi-part rules)
	if fragment.InheritedFromFinding || len(r.RequiredRules) == 0 {
		return findings
	}

	// Process required rules and create findings with auxiliary findings
	return d.processRequiredRules(fragment, currentRaw, r, encodedSegments, findings, logger, paths)
}

// processRequiredRules handles the logic for multi-part rules with auxiliary findings
func (d *Detector) processRequiredRules(fragment Fragment, currentRaw string, r config.Rule, encodedSegments []*codec.EncodedSegment, primaryFindings []report.Finding, logger zerolog.Logger, paths fragmentPaths) []report.Finding {
	if len(primaryFindings) == 0 {
		logger.Debug().Msg("no primary findings to process for required rules")
		return primaryFindings
	}

	// Pre-collect all required rule findings once
	allRequiredFindings := make(map[string][]report.Finding)

	for _, requiredRule := range r.RequiredRules {
		rule, ok := d.Config.Rules[requiredRule.RuleID]
		if !ok {
			logger.Error().Str("rule-id", requiredRule.RuleID).Msg("required rule not found in config")
			continue
		}

		// Mark fragment as inherited to prevent infinite recursion
		inheritedFragment := fragment
		inheritedFragment.InheritedFromFinding = true

		// Call detectRule once for each required rule
		requiredFindings := d.detectRuleWithPaths(inheritedFragment, currentRaw, rule, encodedSegments, paths)
		allRequiredFindings[requiredRule.RuleID] = requiredFindings

		logger.Debug().
			Str("rule-id", requiredRule.RuleID).
			Int("findings", len(requiredFindings)).
			Msg("collected required rule findings")
	}

	var finalFindings []report.Finding

	// Now process each primary finding against the pre-collected required findings
	for _, primaryFinding := range primaryFindings {
		var requiredFindings []*report.RequiredFinding

		for _, requiredRule := range r.RequiredRules {
			foundRequiredFindings, exists := allRequiredFindings[requiredRule.RuleID]
			if !exists {
				continue // Rule wasn't found earlier, skip
			}

			// Filter findings that are within proximity of the primary finding
			for _, requiredFinding := range foundRequiredFindings {
				if d.withinProximity(primaryFinding, requiredFinding, requiredRule) {
					req := &report.RequiredFinding{
						RuleID:      requiredFinding.RuleID,
						StartLine:   requiredFinding.StartLine,
						EndLine:     requiredFinding.EndLine,
						StartColumn: requiredFinding.StartColumn,
						EndColumn:   requiredFinding.EndColumn,
						Line:        requiredFinding.Line,
						Match:       requiredFinding.Match,
						Secret:      requiredFinding.Secret,
					}
					requiredFindings = append(requiredFindings, req)
				}
			}
		}

		// Check if we have at least one auxiliary finding for each required rule
		if len(requiredFindings) > 0 && d.hasAllRequiredRules(requiredFindings, r.RequiredRules) {
			// Create a finding with auxiliary findings
			newFinding := primaryFinding // Copy the primary finding
			newFinding.AddRequiredFindings(requiredFindings)
			finalFindings = append(finalFindings, newFinding)

			logger.Debug().
				Str("primary-rule", r.RuleID).
				Int("primary-line", primaryFinding.StartLine).
				Int("auxiliary-count", len(requiredFindings)).
				Msg("multi-part rule satisfied")
		}
	}

	return finalFindings
}

// hasAllRequiredRules checks if we have at least one auxiliary finding for each required rule
func (d *Detector) hasAllRequiredRules(auxiliaryFindings []*report.RequiredFinding, requiredRules []*config.Required) bool {
	foundRules := make(map[string]bool)
	// AuxiliaryFinding
	for _, aux := range auxiliaryFindings {
		foundRules[aux.RuleID] = true
	}

	for _, required := range requiredRules {
		if !foundRules[required.RuleID] {
			return false
		}
	}

	return true
}

func (d *Detector) withinProximity(primary, required report.Finding, requiredRule *config.Required) bool {
	// fmt.Println(requiredRule.WithinLines)
	// If neither within_lines nor within_columns is set, findings just need to be in the same fragment
	if requiredRule.WithinLines == nil && requiredRule.WithinColumns == nil {
		return true
	}

	// Check line proximity (vertical distance)
	if requiredRule.WithinLines != nil {
		lineDiff := abs(primary.StartLine - required.StartLine)
		if lineDiff > *requiredRule.WithinLines {
			return false
		}
	}

	// Check column proximity (horizontal distance)
	if requiredRule.WithinColumns != nil {
		// Use the start column of each finding for proximity calculation
		colDiff := abs(primary.StartColumn - required.StartColumn)
		if colDiff > *requiredRule.WithinColumns {
			return false
		}
	}

	return true
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (d *Detector) markGitleaksIgnoreUsed(entry string) {
	d.ignoreMutex.Lock()
	d.gitleaksIgnoreUsed[entry] = struct{}{}
	d.ignoreMutex.Unlock()
}

// UnusedGitleaksIgnore returns loaded ignore entries that did not suppress a
// finding in this completed run. "Not triggered" is not proof that an entry is
// obsolete: the relevant file, commit or revision may never have been examined.
func (d *Detector) UnusedGitleaksIgnore() []string {
	d.ignoreMutex.Lock()
	defer d.ignoreMutex.Unlock()
	unused := make([]string, 0)
	for entry := range d.gitleaksIgnore {
		if _, ok := d.gitleaksIgnoreUsed[entry]; !ok {
			unused = append(unused, entry)
		}
	}
	sort.Strings(unused)
	return unused
}

// AddFinding synchronously adds a finding to the findings slice
func (d *Detector) AddFinding(finding report.Finding) {
	globalFingerprint := fmt.Sprintf("%s:%s:%d", finding.File, finding.RuleID, finding.StartLine)
	if finding.Commit != "" {
		finding.Fingerprint = fmt.Sprintf("%s:%s:%s:%d", finding.Commit, finding.File, finding.RuleID, finding.StartLine)
	} else {
		finding.Fingerprint = globalFingerprint
	}

	// check if we should ignore this finding
	logger := logging.With().Str("path", finding.File).Str("rule-id", finding.RuleID).Logger()
	if _, ok := d.gitleaksIgnore[globalFingerprint]; ok {
		d.markGitleaksIgnoreUsed(globalFingerprint)
		logger.Debug().
			Str("fingerprint", globalFingerprint).
			Msg("skipping finding: global fingerprint")
		return
	} else if finding.Commit != "" {
		// Awkward nested if because I'm not sure how to chain these two conditions.
		if _, ok := d.gitleaksIgnore[finding.Fingerprint]; ok {
			d.markGitleaksIgnoreUsed(finding.Fingerprint)
			logger.Debug().
				Str("fingerprint", finding.Fingerprint).
				Msgf("skipping finding: fingerprint")
			return
		}
	}

	if d.baseline != nil && !IsNew(finding, d.Redact, d.baseline) {
		logger.Debug().
			Str("fingerprint", finding.Fingerprint).
			Msgf("skipping finding: baseline")
		return
	}

	d.findingMutex.Lock()
	d.findings = append(d.findings, finding)
	if d.Verbose {
		printFinding(finding, d.NoColor, d.FileLinks)
	}
	d.findingMutex.Unlock()
}

// Findings returns the findings added to the detector
func (d *Detector) Findings() []report.Finding {
	return d.findings
}

// AddCommit synchronously adds a commit to the commit slice
func (d *Detector) addCommit(commit string) {
	d.commitMutex.Lock()
	d.commitMap[commit] = true
	d.commitMutex.Unlock()
}

// commitAllowed reports whether any allowlist contains a matching commit.
// Informative reporting uses it only to label an already-made exclusion
// decision accurately.
func commitAllowed(commit string, allowlists []*config.Allowlist) bool {
	if commit == "" {
		return false
	}
	for _, a := range allowlists {
		if ok, _ := a.CommitAllowed(commit); ok {
			return true
		}
	}
	return false
}

// checkCommitOrPathAllowed evaluates |fragment| against all provided |allowlists|.
//
// If the match condition is "OR", only commit and path are checked.
// Otherwise, if regexes or stopwords are defined this will fail.
func checkCommitOrPathAllowed(
	logger zerolog.Logger,
	fragment Fragment,
	paths fragmentPaths,
	allowlists []*config.Allowlist,
) (bool, *zerolog.Event) {
	if fragment.FilePath == "" && fragment.CommitSHA == "" {
		return false, nil
	}

	for _, a := range allowlists {
		var (
			isAllowed        bool
			allowlistChecks  []bool
			commitAllowed, _ = a.CommitAllowed(fragment.CommitSHA)
			pathAllowed      = paths.allowed(a)
		)
		// If the condition is "AND" we need to check all conditions.
		if a.MatchCondition == config.AllowlistMatchAnd {
			if len(a.Commits) > 0 {
				allowlistChecks = append(allowlistChecks, commitAllowed)
			}
			if len(a.Paths) > 0 {
				allowlistChecks = append(allowlistChecks, pathAllowed)
			}
			// These will be checked later.
			if len(a.Regexes) > 0 {
				continue
			}
			if len(a.StopWords) > 0 {
				continue
			}

			isAllowed = allTrue(allowlistChecks)
		} else {
			isAllowed = commitAllowed || pathAllowed
		}
		if isAllowed {
			event := logger.Trace().Str("condition", a.MatchCondition.String())
			if commitAllowed {
				event.Bool("allowed-commit", commitAllowed)
			}
			if pathAllowed {
				event.Bool("allowed-path", pathAllowed)
			}
			return true, event
		}
	}
	return false, nil
}

// checkFindingAllowed evaluates |finding| against all provided |allowlists|.
//
// If the match condition is "OR", only regex and stopwords are run. (Commit and path should be handled separately).
// Otherwise, all conditions are checked.
//
// TODO: The method signature is awkward. I can't think of a better way to log helpful info.
func checkFindingAllowed(
	logger zerolog.Logger,
	finding report.Finding,
	fragment Fragment,
	paths fragmentPaths,
	currentLine string,
	allowlists []*config.Allowlist,
) (bool, *zerolog.Event) {
	for _, a := range allowlists {
		allowlistTarget := finding.Secret
		switch a.RegexTarget {
		case "match":
			allowlistTarget = finding.Match
		case "line":
			allowlistTarget = currentLine
		}

		var (
			checks              []bool
			isAllowed           bool
			commitAllowed       bool
			commit              string
			pathAllowed         bool
			regexAllowed        = a.RegexAllowed(allowlistTarget)
			containsStopword, _ = a.ContainsStopWord(finding.Secret)
		)
		// If the condition is "AND" we need to check all conditions.
		if a.MatchCondition == config.AllowlistMatchAnd {
			// Determine applicable checks.
			if len(a.Commits) > 0 {
				commitAllowed, commit = a.CommitAllowed(fragment.CommitSHA)
				checks = append(checks, commitAllowed)
			}
			if len(a.Paths) > 0 {
				pathAllowed = paths.allowed(a)
				checks = append(checks, pathAllowed)
			}
			if len(a.Regexes) > 0 {
				checks = append(checks, regexAllowed)
			}
			if len(a.StopWords) > 0 {
				checks = append(checks, containsStopword)
			}

			isAllowed = allTrue(checks)
		} else {
			isAllowed = regexAllowed || containsStopword
		}

		if isAllowed {
			event := logger.Trace().
				Str("condition", a.MatchCondition.String())
			if commitAllowed {
				event.Str("allowed-commit", commit)
			}
			if pathAllowed {
				event.Bool("allowed-path", pathAllowed)
			}
			if regexAllowed {
				event.Bool("allowed-regex", regexAllowed)
			}
			if containsStopword {
				event.Bool("allowed-stopword", true)
			}
			return true, event
		}
	}
	return false, nil
}

func allTrue(bools []bool) bool {
	for _, check := range bools {
		if !check {
			return false
		}
	}
	return true
}
