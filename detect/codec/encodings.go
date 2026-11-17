package codec

import (
	"fmt"
	"math"
	"strings"

	"github.com/zricethezav/gitleaks/v8/regexp"
)

var (
	// encodingsRe is a regex built by combining all the encoding patterns
	// into named capture groups so that a single pass can detect multiple
	// encodings
	encodingsRe *regexp.Regexp
	// encodings contains all the encoding configurations for the detector.
	// The precedence is important. You want more specific encodings to
	// have a higher precedence or encodings that partially encode the
	// values (e.g. percent) unlike encodings that fully encode the string
	// (e.g. base64). If two encoding matches overlap the decoder will use
	// this order to determine which encoding should wait till the next pass.
	encodings = []*encoding{
		{
			kind:    percentKind,
			pattern: `%[0-9A-Fa-f]{2}(?:.*%[0-9A-Fa-f]{2})?`,
			decode:  decodePercent,
		},
		{
			kind:    unicodeKind,
			pattern: `(?:(?:U\+[a-fA-F0-9]{4}(?:\s|$))+|(?i)(?:\\{1,2}u[a-fA-F0-9]{4})+)`,
			decode:  decodeUnicode,
		},
		{
			kind:    hexKind,
			pattern: `[0-9A-Fa-f]{32,}`,
			decode:  decodeHex,
		},
		{
			kind:    base64Kind,
			pattern: `[\w\/+-]{16,}={0,2}`,
			decode:  decodeBase64,
		},
	}
)

// encodingNames is used to map the encodingKinds to their name
var encodingNames = []string{
	"percent",
	"unicode",
	"hex",
	"base64",
}

// encodingKind can be or'd together to capture all of the unique encodings
// that were present in a segment
type encodingKind int

var (
	// make sure these go up by powers of 2
	percentKind = encodingKind(1)
	unicodeKind = encodingKind(2)
	hexKind     = encodingKind(4)
	base64Kind  = encodingKind(8)
)

func (e encodingKind) String() string {
	i := int(math.Log2(float64(e)))
	if i >= len(encodingNames) {
		return ""
	}
	return encodingNames[i]
}

// kinds returns a list of encodingKinds combined in this one
func (e encodingKind) kinds() []encodingKind {
	kinds := []encodingKind{}

	for i := 0; i < len(encodingNames); i++ {
		if kind := int(e) & int(math.Pow(2, float64(i))); kind != 0 {
			kinds = append(kinds, encodingKind(kind))
		}
	}

	return kinds
}

// encodingMatch represents a match of an encoding in the text
type encodingMatch struct {
	encoding *encoding
	startEnd
}

// encoding represent a type of coding supported by the decoder.
type encoding struct {
	// the kind of decoding (e.g. base64, etc)
	kind encodingKind
	// the regex pattern that matches the encoding format
	pattern string
	// take the match and return the decoded value
	decode func(string) string
	// determine which encoding should win out when two overlap
	precedence int
}

func init() {
	count := len(encodings)
	namedPatterns := make([]string, count)
	for i, encoding := range encodings {
		encoding.precedence = count - i
		namedPatterns[i] = fmt.Sprintf(
			"(?P<%s>%s)",
			encoding.kind,
			encoding.pattern,
		)
	}
	encodingsRe = regexp.MustCompile(strings.Join(namedPatterns, "|"))
}

// findEncodingMatches finds as many encodings as it can for this pass.
//
// The combined encoding regexp is expensive on ordinary source text. A cheap
// conservative prefilter rejects data that cannot contain any supported
// encoding, avoiding a full RE2 scan for the common case. It is deliberately
// a superset of the regexp: every match it allows still goes through the
// existing decoder and precedence rules unchanged.
func findEncodingMatches(data string) []encodingMatch {
	if !mayContainEncoding(data) {
		return nil
	}

	var all []encodingMatch
	for _, matchIndex := range encodingsRe.FindAllStringSubmatchIndex(data, -1) {
		// Add the encodingMatch with its proper encoding
		for i, j := 2, 0; i < len(matchIndex); i, j = i+2, j+1 {
			if matchIndex[i] > -1 {
				all = append(all, encodingMatch{
					encoding: encodings[j],
					startEnd: startEnd{
						start: matchIndex[i],
						end:   matchIndex[i+1],
					},
				})
			}
		}
	}

	totalMatches := len(all)
	if totalMatches == 1 {
		return all
	}

	// filter out lower precedence ones that overlap their neigbors
	filtered := make([]encodingMatch, 0, len(all))
	for i, m := range all {
		if i > 0 {
			prev := all[i-1]
			if m.overlaps(prev.startEnd) && prev.encoding.precedence > m.encoding.precedence {
				continue // skip this one
			}
		}
		if i+1 < totalMatches {
			next := all[i+1]
			if m.overlaps(next.startEnd) && next.encoding.precedence > m.encoding.precedence {
				continue // skip this one
			}
		}
		filtered = append(filtered, m)
	}

	return filtered
}

// mayContainEncoding reports whether data can possibly contain an encoding
// understood by encodingsRe. It mirrors the necessary character conditions:
//
//   - percent encoding requires a '%'
//   - Unicode escapes require "U+" or a backslash followed by u/U
//   - base64/hex requires a run of at least 16 characters from the base64
//     character class, including a likely-base64 character. A run without
//     one cannot decode (decodeBase64 rejects it), so it cannot produce a
//     segment; hex runs are a subset of this class and are also covered.
//
// The scan exits as soon as a candidate is found so encoded content pays only
// a short prefix check before the real regexp runs.
func mayContainEncoding(data string) bool {
	run := 0
	likely := false
	for i := 0; i < len(data); i++ {
		c := data[i]
		switch c {
		case '%':
			return true
		case '\\':
			if i+1 < len(data) && (data[i+1] == 'u' || data[i+1] == 'U') {
				return true
			}
		case 'U':
			if i+1 < len(data) && data[i+1] == '+' {
				return true
			}
		}

		if isBase64Alphabet(c) {
			run++
			if isLikelyBase64Char(c) {
				likely = true
			}
			if run >= 16 && likely {
				return true
			}
			continue
		}
		run = 0
		likely = false
	}
	return false
}

func isBase64Alphabet(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9' ||
		c == '_' || c == '/' || c == '+' || c == '-'
}

func isLikelyBase64Char(c byte) bool {
	return c >= '0' && c <= '9' || c == '+' || c == '/' || c == '-' || c == '_'
}
