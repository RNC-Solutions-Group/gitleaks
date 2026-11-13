package detect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetLocation tests the getLocation function.
func TestGetLocation(t *testing.T) {
	tests := []struct {
		linePairs    [][]int
		start        int
		end          int
		wantLocation Location
	}{
		{
			linePairs: [][]int{
				{0, 39},
				{40, 55},
				{56, 57},
			},
			start: 35,
			end:   38,
			wantLocation: Location{
				startLine:            1,
				startColumn:          36,
				endLine:              1,
				endColumn:            38,
				startColumnCodepoint: 36,
				endColumnCodepoint:   38,
				startLineIndex:       1,
				endLineIndex:         40,
			},
		},
		{
			linePairs: [][]int{
				{0, 39},
				{40, 55},
				{56, 57},
			},
			start: 40,
			end:   44,
			wantLocation: Location{
				startLine:            2,
				startColumn:          1,
				endLine:              2,
				endColumn:            4,
				startColumnCodepoint: 1,
				endColumnCodepoint:   4,
				startLineIndex:       41,
				endLineIndex:         56,
			},
		},
	}

	for _, test := range tests {
		loc := location(test.linePairs, "", []int{test.start, test.end})
		assert.Equal(t, test.wantLocation, loc)
	}
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
