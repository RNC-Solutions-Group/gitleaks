package detect

import "unicode/utf8"

// Location represents a location in a file
type Location struct {
	startLine   int
	endLine     int
	startColumn int
	endColumn   int
	// Codepoint columns are used by SARIF, whose column fields are defined in
	// codepoints. They are kept separate so existing JSON/CSV/console columns
	// retain their historical byte-based contract.
	startColumnCodepoint int
	endColumnCodepoint   int
	startLineIndex       int
	endLineIndex         int
}

func location(newlineIndices [][]int, raw string, matchIndex []int) Location {
	var (
		prevNewLine int
		location    Location
		lineSet     bool
		_lineNum    int
		// lineStart is the byte index where the current line starts,
		// i.e. one byte after the newline ending the previous line.
		lineStart int
	)

	start := matchIndex[0]
	end := matchIndex[1]

	// default startLineIndex to 0
	location.startLineIndex = 0

	// Fixes: https://github.com/zricethezav/gitleaks/issues/1037
	// When a fragment does NOT have any newlines, a default "newline"
	// will be counted to make the subsequent location calculation logic work
	// for fragments will no newlines.
	if len(newlineIndices) == 0 {
		newlineIndices = [][]int{
			{len(raw), len(raw) + 1},
		}
	}

	for lineNum, pair := range newlineIndices {
		_lineNum = lineNum
		newLineByteIndex := pair[0]
		if prevNewLine <= start && start < newLineByteIndex {
			lineSet = true
			location.startLine = lineNum
			location.endLine = lineNum
			location.startColumn = (start - prevNewLine) + 1 // +1 because counting starts at 1
			location.startColumnCodepoint = location.startColumn
			if raw != "" {
				location.startColumnCodepoint = columnOf(raw, lineStart, start) + 1
			}
			location.startLineIndex = lineStart
			location.endLineIndex = newLineByteIndex
		}
		if prevNewLine < end && end <= newLineByteIndex {
			location.endLine = lineNum
			location.endColumn = (end - prevNewLine)
			location.endColumnCodepoint = location.endColumn
			if raw != "" {
				location.endColumnCodepoint = columnOf(raw, lineStart, end)
			}
			location.endLineIndex = newLineByteIndex
		}

		prevNewLine = pair[0]
		lineStart = pair[0] + 1
	}

	if !lineSet {
		// if lines never get set then that means the secret is most likely
		// on the last line of the diff output and the diff output does not have
		// a newline
		location.startColumn = (start - prevNewLine) + 1 // +1 because counting starts at 1
		location.endColumn = (end - prevNewLine)
		location.startColumnCodepoint = location.startColumn
		location.endColumnCodepoint = location.endColumn
		if raw != "" {
			location.startColumnCodepoint = columnOf(raw, lineStart, start) + 1
			location.endColumnCodepoint = columnOf(raw, lineStart, end)
		}
		location.startLine = _lineNum + 1
		location.endLine = _lineNum + 1
		location.startLineIndex = lineStart

		// search for new line byte index
		i := 0
		for end+i < len(raw) {
			if raw[end+i] == '\n' {
				break
			}
			if raw[end+i] == '\r' {
				break
			}
			i++
		}
		location.endLineIndex = end + i
	}
	if location.startLineIndex > location.endLineIndex {
		// An empty trimmed match on a trailing newline moves lineStart past
		// endLineIndex; clamp so Finding.Line slicing cannot panic.
		location.startLineIndex = location.endLineIndex
	}
	return location
}

// columnOf returns how far into the line an offset is, counted in Unicode
// codepoints rather than bytes. SARIF defines startColumn and endColumn in
// codepoints; byte offsets point past the secret on lines containing Arabic,
// CJK, emoji, or other multibyte text. For ASCII the values are identical.
func columnOf(raw string, lineStart, offset int) int {
	if lineStart < 0 || offset < lineStart || offset > len(raw) {
		// Out of range: fall back to the byte distance rather than panic.
		return offset - lineStart
	}
	return utf8.RuneCountInString(raw[lineStart:offset])
}
