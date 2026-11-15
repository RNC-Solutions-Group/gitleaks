package codec

import (
	"unicode"
	"unicode/utf8"
)

var printableASCII [256]bool

func init() {
	for b := 0; b < len(printableASCII); b++ {
		if '\x08' < b && b < '\x7f' {
			printableASCII[b] = true
		}
	}
}

// isPrintable reports whether b is valid UTF-8 made up entirely of printable
// characters. ASCII uses the historical table, preserving its allowances for
// tab/newline/carriage return. Multibyte runes must be printable Unicode;
// control characters, format characters and invalid UTF-8 are rejected.
//
// #2000: rejecting every non-ASCII byte caused valid decoded text such as
// "café" or "日本語" to be silently skipped.
func isPrintable(b []byte) bool {
	for i := 0; i < len(b); {
		c := b[i]
		if c < utf8.RuneSelf {
			if !printableASCII[c] {
				return false
			}
			i++
			continue
		}

		r, size := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError && size == 1 {
			return false
		}
		if !unicode.IsPrint(r) {
			return false
		}
		i += size
	}
	return true
}

// hasByte can be used to check if a string has at least one of the provided
// bytes. Note: make sure byteset is long enough to handle the largest byte in
// the string.
func hasByte(data string, byteset []bool) bool {
	for i := 0; i < len(data); i++ {
		if byteset[data[i]] {
			return true
		}
	}
	return false
}
