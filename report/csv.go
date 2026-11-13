package report

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type CsvReporter struct {
}

var _ Reporter = (*CsvReporter)(nil)

// sanitizeCSVCell neutralizes spreadsheet formula prefixes. Findings can carry
// attacker-controlled paths, author metadata, matches and messages, and a CSV
// opened in Excel/LibreOffice/Sheets may otherwise evaluate those cells. The
// leading apostrophe is the conventional text marker; the original value is
// otherwise preserved. Leading whitespace, controls and a UTF-8 BOM are
// skipped when deciding whether the cell begins like a formula.
func sanitizeCSVCell(value string) string {
	if value == "" {
		return value
	}
	candidate := strings.TrimLeftFunc(value, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == '\ufeff'
	})
	if candidate == "" {
		return value
	}
	first, _ := utf8.DecodeRuneInString(candidate)
	if isCSVFormulaRune(first) {
		return "'" + value
	}
	return value
}

func isCSVFormulaRune(r rune) bool {
	switch r {
	case '=', '+', '-', '@':
		// ASCII spreadsheet formula prefixes.
		return true
	case '＝', '＋', '－', '＠':
		// Fullwidth variants accepted by some spreadsheet input paths.
		return true
	default:
		return false
	}
}

func sanitizeCSVRow(row []string) []string {
	out := make([]string, len(row))
	for i, value := range row {
		out[i] = sanitizeCSVCell(value)
	}
	return out
}

func (r *CsvReporter) Write(w io.WriteCloser, findings []Finding) error {
	if len(findings) == 0 {
		return nil
	}

	var (
		cw  = csv.NewWriter(w)
		err error
	)
	columns := []string{"RuleID",
		"Commit",
		"File",
		"SymlinkFile",
		"Secret",
		"Match",
		"StartLine",
		"EndLine",
		"StartColumn",
		"EndColumn",
		"Author",
		"Message",
		"Date",
		"Email",
		"Fingerprint",
		"Tags",
	}
	// A miserable attempt at "omitempty" so tests don't yell at me.
	if findings[0].Link != "" {
		columns = append(columns, "Link")
	}

	if err = cw.Write(columns); err != nil {
		return err
	}
	for _, f := range findings {
		row := []string{f.RuleID,
			f.Commit,
			f.File,
			f.SymlinkFile,
			f.Secret,
			f.Match,
			strconv.Itoa(f.StartLine),
			strconv.Itoa(f.EndLine),
			strconv.Itoa(f.StartColumn),
			strconv.Itoa(f.EndColumn),
			f.Author,
			f.Message,
			f.Date,
			f.Email,
			f.Fingerprint,
			strings.Join(f.Tags, " "),
		}
		if findings[0].Link != "" {
			row = append(row, f.Link)
		}

		if err = cw.Write(sanitizeCSVRow(row)); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}
