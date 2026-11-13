package detect

import (
	// "encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"

	"charm.land/lipgloss/v2"
)

// samePath compares two paths lexically, including their absolute forms. It
// performs no filesystem access.
func samePath(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}

	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	return errA == nil && errB == nil && filepath.Clean(absA) == filepath.Clean(absB)
}

// sameFileOrPath compares two paths lexically and, when both exist, by file
// identity. It is used to avoid scanning the report destination as an input.
func sameFileOrPath(a, b string) bool {
	if samePath(a, b) {
		return true
	}
	if a == "" || b == "" {
		return false
	}

	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(infoA, infoB)
}

var linkCleaner = strings.NewReplacer(
	" ", "%20",
	"%", "%25",
)

func createScmLink(remote *sources.RemoteInfo, finding report.Finding) string {
	if remote.Platform == scm.UnknownPlatform ||
		remote.Platform == scm.NoPlatform ||
		finding.Commit == "" {
		return ""
	}

	// Clean the path.
	filePath, _, hasInnerPath := strings.Cut(finding.File, sources.InnerPathSeparator)
	filePath = linkCleaner.Replace(filePath)

	switch remote.Platform {
	case scm.GitHubPlatform:
		link := fmt.Sprintf("%s/blob/%s/%s", remote.Url, finding.Commit, filePath)
		if hasInnerPath {
			return link
		}
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext == ".ipynb" || ext == ".md" {
			link += "?plain=1"
		}
		if finding.StartLine != 0 {
			link += fmt.Sprintf("#L%d", finding.StartLine)
		}
		if finding.EndLine != finding.StartLine {
			link += fmt.Sprintf("-L%d", finding.EndLine)
		}
		return link
	case scm.GitLabPlatform:
		link := fmt.Sprintf("%s/blob/%s/%s", remote.Url, finding.Commit, filePath)
		if hasInnerPath {
			return link
		}
		if finding.StartLine != 0 {
			link += fmt.Sprintf("#L%d", finding.StartLine)
		}
		if finding.EndLine != finding.StartLine {
			link += fmt.Sprintf("-%d", finding.EndLine)
		}
		return link
	case scm.AzureDevOpsPlatform:
		link := fmt.Sprintf("%s/commit/%s?path=/%s", remote.Url, finding.Commit, filePath)
		// Add line information if applicable
		if hasInnerPath {
			return link
		}
		if finding.StartLine != 0 {
			link += fmt.Sprintf("&line=%d", finding.StartLine)
		}
		if finding.EndLine != finding.StartLine {
			link += fmt.Sprintf("&lineEnd=%d", finding.EndLine)
		}
		// This is a bit dirty, but Azure DevOps does not highlight the line when the lineStartColumn and lineEndColumn are not provided
		link += "&lineStartColumn=1&lineEndColumn=10000000&type=2&lineStyle=plain&_a=files"
		return link
	case scm.GiteaPlatform:
		link := fmt.Sprintf("%s/src/commit/%s/%s", remote.Url, finding.Commit, filePath)
		if hasInnerPath {
			return link
		}
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext == ".ipynb" || ext == ".md" {
			link += "?display=source"
		}
		if finding.StartLine != 0 {
			link += fmt.Sprintf("#L%d", finding.StartLine)
		}
		if finding.EndLine != finding.StartLine {
			link += fmt.Sprintf("-L%d", finding.EndLine)
		}
		return link
	case scm.BitbucketPlatform:
		link := fmt.Sprintf("%s/src/%s/%s", remote.Url, finding.Commit, filePath)
		if hasInnerPath {
			return link
		}
		if finding.StartLine != 0 {
			link += fmt.Sprintf("#lines-%d", finding.StartLine)
		}
		if finding.EndLine != finding.StartLine {
			link += fmt.Sprintf(":%d", finding.EndLine)
		}
		return link
	default:
		// This should never happen.
		return ""
	}
}

// shannonEntropy calculates the entropy of data using the formula defined here:
// https://en.wiktionary.org/wiki/Shannon_entropy
// Another way to think about what this is doing is calculating the number of bits
// needed to on average encode the data. So, the higher the entropy, the more random the data, the
// more bits needed to encode that data.
func shannonEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}

	charCounts := make(map[rune]int)
	for _, char := range data {
		charCounts[char]++
	}

	invLength := 1.0 / float64(len(data))
	for _, count := range charCounts {
		freq := float64(count) * invLength
		entropy -= freq * math.Log2(freq)
	}

	return entropy
}

// filter will dedupe and redact findings
func filter(findings []report.Finding, redact uint) []report.Finding {
	var retFindings []report.Finding
	for _, f := range findings {
		include := true
		if strings.Contains(strings.ToLower(f.RuleID), "generic") {
			for _, fPrime := range findings {
				if f.StartLine == fPrime.StartLine &&
					f.Commit == fPrime.Commit &&
					f.RuleID != fPrime.RuleID &&
					strings.Contains(fPrime.Secret, f.Secret) &&
					!strings.Contains(strings.ToLower(fPrime.RuleID), "generic") {

					logging.Trace().
						Str("generic-rule", f.RuleID).
						Str("preferred-rule", fPrime.RuleID).
						Int("line", f.StartLine).
						Msg("skipping generic finding")
					include = false
					break
				}
			}
		}
		if include {
			retFindings = append(retFindings, f)
		}
	}

	if redact == 0 {
		return retFindings
	}

	// A finding's Line can contain other reported credentials on a different
	// line range of the same file/commit (for example a multiline match and a
	// same-line match). Group known exposures by source identity rather than by
	// exact line range so a selected finding cannot render a neighboring
	// credential. Empty-capture findings contribute their Match when it is not
	// a filename-only marker; otherwise their surrounding context would remain
	// unredacted.
	secretsBySource := make(map[string][]string)
	for _, f := range findings {
		key := findingSourceKey(f)
		if f.Secret != "" {
			secretsBySource[key] = append(secretsBySource[key], f.Secret)
		} else if f.Match != "" && !strings.HasPrefix(f.Match, "file detected:") {
			secretsBySource[key] = append(secretsBySource[key], f.Match)
		}
	}
	for i := range retFindings {
		retFindings[i].RedactWith(redact, secretsBySource[findingSourceKey(retFindings[i])])
	}
	return retFindings
}

func findingSourceKey(f report.Finding) string {
	return f.File + "\x00" + f.Commit
}

func printFinding(f report.Finding, noColor bool) {
	// trim all whitespace and tabs
	f.Line = strings.TrimSpace(f.Line)
	f.Secret = strings.TrimSpace(f.Secret)
	f.Match = strings.TrimSpace(f.Match)

	isFileMatch := strings.HasPrefix(f.Match, "file detected:")
	skipColor := noColor
	finding := ""
	var secret lipgloss.Style
	// Styled output is written through a down-sampling writer so that redirected
	// or piped output stays free of ANSI escape sequences, as it was with
	// Lip Gloss v1.
	out := report.ColorWriter()

	// Matches from filenames do not have a |line| or |secret|
	if !isFileMatch {
		matchInLineIDX := strings.Index(f.Line, f.Match)
		secretInMatchIdx := strings.Index(f.Match, f.Secret)

		skipColor = false

		if matchInLineIDX == -1 || secretInMatchIdx == -1 || noColor {
			skipColor = true
			matchInLineIDX = 0
			secretInMatchIdx = 0
		}

		start := f.Line[0:matchInLineIDX]
		startRunes := []rune(start)
		if len(startRunes) > 20 {
			start = "..." + string(startRunes[len(startRunes)-20:])
		}

		matchBeginning := lipgloss.NewStyle().SetString(f.Match[0:secretInMatchIdx]).Foreground(lipgloss.Color("#f5d445"))
		secret = lipgloss.NewStyle().SetString(f.Secret).
			Bold(true).
			Italic(true).
			Foreground(lipgloss.Color("#f05c07"))
		matchEndIdx := secretInMatchIdx + len(f.Secret)
		if matchEndIdx > len(f.Match) {
			matchEndIdx = len(f.Match)
		}
		matchEnd := lipgloss.NewStyle().SetString(f.Match[matchEndIdx:]).Foreground(lipgloss.Color("#f5d445"))

		lineEndIdx := matchInLineIDX + len(f.Match)
		if len(f.Line)-1 <= lineEndIdx {
			lineEndIdx = len(f.Line)
		}

		lineEnd := f.Line[lineEndIdx:]

		if len([]rune(f.Secret)) > 100 {
			secret = lipgloss.NewStyle().SetString(truncateRunes(f.Secret, 100)).
				Bold(true).
				Italic(true).
				Foreground(lipgloss.Color("#f05c07"))
		}
		if len([]rune(lineEnd)) > 20 {
			lineEnd = string([]rune(lineEnd)[:20]) + "..."
		}

		finding = fmt.Sprintf("%s%s%s%s%s\n", strings.TrimPrefix(strings.TrimLeft(start, " "), "\n"), matchBeginning, secret, matchEnd, lineEnd)
	}

	if skipColor || isFileMatch {
		fmt.Fprintf(out, "%-12s %s\n", "Finding:", f.Match)
		fmt.Fprintf(out, "%-12s %s\n", "Secret:", f.Secret)
	} else {
		fmt.Fprintf(out, "%-12s %s", "Finding:", finding)
		fmt.Fprintf(out, "%-12s %s\n", "Secret:", secret)
	}

	fmt.Fprintf(out, "%-12s %s\n", "RuleID:", f.RuleID)
	fmt.Fprintf(out, "%-12s %f\n", "Entropy:", f.Entropy)

	if f.File == "" {
		f.PrintRequiredFindings()
		fmt.Fprintln(out, "")
		return
	}
	if len(f.Tags) > 0 {
		fmt.Fprintf(out, "%-12s %s\n", "Tags:", f.Tags)
	}
	fmt.Fprintf(out, "%-12s %s\n", "File:", f.File)
	fmt.Fprintf(out, "%-12s %d\n", "Line:", f.StartLine)
	if f.Commit == "" {
		fmt.Fprintf(out, "%-12s %s\n", "Fingerprint:", f.Fingerprint)
		f.PrintRequiredFindings()
		fmt.Fprintln(out, "")
		return
	}
	fmt.Fprintf(out, "%-12s %s\n", "Commit:", f.Commit)
	fmt.Fprintf(out, "%-12s %s\n", "Author:", f.Author)
	fmt.Fprintf(out, "%-12s %s\n", "Email:", f.Email)
	fmt.Fprintf(out, "%-12s %s\n", "Date:", f.Date)
	fmt.Fprintf(out, "%-12s %s\n", "Fingerprint:", f.Fingerprint)
	if f.Link != "" {
		fmt.Fprintf(out, "%-12s %s\n", "Link:", f.Link)
	}

	f.PrintRequiredFindings()
	fmt.Fprintln(out, "")
}

func truncateRunes(value string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max-3]) + "..."
}
