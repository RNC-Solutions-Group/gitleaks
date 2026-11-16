package detect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetLocation tests the getLocation function.
func TestGetLocation(t *testing.T) {
	tests := []struct {
		name         string
		raw          string
		start        int
		end          int
		wantLocation Location
	}{
		{
			name:  "match on first line",
			raw:   strings.Repeat("x", 35) + "abc" + strings.Repeat("y", 1) + "\ntail\n",
			start: 35,
			end:   38,
			wantLocation: Location{
				startLine:            0,
				startColumn:          36,
				endLine:              0,
				endColumn:            38,
				startColumnCodepoint: 36,
				endColumnCodepoint:   38,
				startLineIndex:       0,
				endLineIndex:         39,
			},
		},
		{
			name:  "match on later line",
			raw:   "first line\n....abc\n",
			start: 15,
			end:   18,
			wantLocation: Location{
				startLine:            1,
				startColumn:          5,
				endLine:              1,
				endColumn:            7,
				startColumnCodepoint: 5,
				endColumnCodepoint:   7,
				startLineIndex:       11,
				endLineIndex:         18,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loc := location(newLineRegexp.FindAllStringIndex(test.raw, -1), test.raw, []int{test.start, test.end})
			assert.Equal(t, test.wantLocation, loc)
		})
	}
}

func TestGetLocationColumnsAreRelativeToLineStart(t *testing.T) {
	raw := "first line\n..f..\nlast line\n"
	start := strings.Index(raw, "\n..f..") + 3
	loc := location(newLineRegexp.FindAllStringIndex(raw, -1), raw, []int{start, start + 1})

	// #1424: a later-line match three bytes into its line is column 3, not 4.
	assert.Equal(t, 1, loc.startLine)
	assert.Equal(t, 3, loc.startColumn)
	assert.Equal(t, 3, loc.endColumn)

	first := location(newLineRegexp.FindAllStringIndex(raw, -1), raw, []int{0, 1})
	assert.Equal(t, 0, first.startLine)
	assert.Equal(t, 1, first.startColumn)
}

// TestGetLocationCountsCodepoints checks that columns are counted in Unicode
// codepoints rather than bytes, which is what the SARIF specification requires.
// https://github.com/gitleaks/gitleaks/issues/1962
func TestGetLocationCountsCodepoints(t *testing.T) {
	const secret = "token=abc123"

	tests := []struct {
		name   string
		prefix string
		// Column the secret starts at, counting characters from 1.
		wantStartColumn int
	}{
		{name: "ascii", prefix: "key: ", wantStartColumn: 6},
		{name: "arabic", prefix: "مفتاح سري: ", wantStartColumn: 12},
		{name: "cjk", prefix: "密钥: ", wantStartColumn: 5},
		{name: "emoji", prefix: "🔑🔒 ", wantStartColumn: 4},
		{name: "leading combining mark", prefix: "e\u0301: ", wantStartColumn: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw := test.prefix + secret + "\n"
			newlineIndices := newLineRegexp.FindAllStringIndex(raw, -1)
			start := len(test.prefix) // byte offset, as the regex engine reports it
			end := start + len(secret)

			loc := location(newlineIndices, raw, []int{start, end})

			assert.Equal(t, test.wantStartColumn, loc.startColumnCodepoint)
			assert.Equal(t, test.wantStartColumn+len(secret)-1, loc.endColumnCodepoint,
				"the secret is ASCII, so it spans len(secret) columns")
		})
	}
}

func TestGetLocationCountsCodepointsOnLaterCRLFLine(t *testing.T) {
	raw := "αβγ\r\nkey=🔑abc\r\n"
	start := strings.Index(raw, "key=")
	end := start + len("key=🔑abc")
	loc := location(newLineRegexp.FindAllStringIndex(raw, -1), raw, []int{start, end})

	assert.Equal(t, 1, loc.startLine)
	assert.Equal(t, 1, loc.startColumnCodepoint)
	assert.Equal(t, 1, loc.endLine)
	assert.Equal(t, 8, loc.endColumnCodepoint)
	assert.Equal(t, "key=🔑abc", strings.TrimSuffix(raw[loc.startLineIndex:loc.endLineIndex], "\r"))
}
