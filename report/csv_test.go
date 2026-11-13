package report

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeCSVCellUnicodeAndIdempotence(t *testing.T) {
	tests := map[string]string{
		"\u200b=SUM(1,1)":    "'\u200b=SUM(1,1)",
		"\u00a0+cmd":         "'\u00a0+cmd",
		"＝1+1":               "'＝1+1",
		"－1+1":               "'－1+1",
		"＠user":              "'＠user",
		"  normal text":      "  normal text",
		"normal":             "normal",
		"'already-quoted":    "'already-quoted",
		"\ufeff=HYPERLINK()": "'\ufeff=HYPERLINK()",
	}
	for input, want := range tests {
		got := sanitizeCSVCell(input)
		assert.Equal(t, want, got)
		assert.Equal(t, got, sanitizeCSVCell(got), "sanitization must be idempotent for %q", input)
	}
}

func TestCSVFormulaInjectionIsNeutralized(t *testing.T) {
	path := filepath.Join(t.TempDir(), "formula.csv")
	file, err := os.Create(path)
	require.NoError(t, err)

	finding := Finding{
		RuleID:      "=2+2",
		Commit:      "+cmd",
		File:        "  =path",
		SymlinkFile: "@link",
		Secret:      "-1+1",
		Match:       "\t=HYPERLINK(\"https://example.invalid\")",
		Author:      "\ufeff=author",
		Message:     "\n=SUM(1,1)",
		Date:        "2026-01-01T00:00:00Z",
		Email:       "@example.invalid",
		Fingerprint: "=fingerprint",
		Tags:        []string{"=tag1", "tag2"},
	}
	require.NoError(t, (&CsvReporter{}).Write(file, []Finding{finding}))
	require.NoError(t, file.Close())

	handle, err := os.Open(path)
	require.NoError(t, err)
	defer handle.Close()
	rows, err := csv.NewReader(handle).ReadAll()
	require.NoError(t, err)
	require.Len(t, rows, 2)
	row := rows[1]

	assert.Equal(t, "'=2+2", row[0])
	assert.Equal(t, "'+cmd", row[1])
	assert.Equal(t, "'  =path", row[2])
	assert.Equal(t, "'@link", row[3])
	assert.Equal(t, "'-1+1", row[4])
	assert.Equal(t, "'\t=HYPERLINK(\"https://example.invalid\")", row[5])
	assert.Equal(t, "'\ufeff=author", row[10])
	assert.Equal(t, "'\n=SUM(1,1)", row[11])
	assert.Equal(t, "2026-01-01T00:00:00Z", row[12], "normal dates remain unchanged")
	assert.Equal(t, "'@example.invalid", row[13])
	assert.Equal(t, "'=fingerprint", row[14])
	assert.Equal(t, "'=tag1 tag2", row[15])
}

func TestWriteCSV(t *testing.T) {
	tests := []struct {
		findings       []Finding
		testReportName string
		expected       string
		wantEmpty      bool
	}{
		{
			testReportName: "simple",
			expected:       filepath.Join(expectPath, "report", "csv_simple.csv"),
			findings: []Finding{
				{
					RuleID:      "test-rule",
					Match:       "line containing secret",
					Secret:      "a secret",
					StartLine:   1,
					EndLine:     2,
					StartColumn: 1,
					EndColumn:   2,
					Message:     "opps",
					File:        "auth.py",
					SymlinkFile: "",
					Commit:      "0000000000000000",
					Author:      "John Doe",
					Email:       "johndoe@gmail.com",
					Date:        "10-19-2003",
					Fingerprint: "fingerprint",
					Tags:        []string{"tag1", "tag2", "tag3"},
				},
			}},
		{

			wantEmpty:      true,
			testReportName: "empty",
			expected:       filepath.Join(expectPath, "report", "this_should_not_exist.csv"),
			findings:       []Finding{},
		},
	}

	reporter := CsvReporter{}
	for _, test := range tests {
		t.Run(test.testReportName, func(t *testing.T) {
			tmpfile, err := os.Create(filepath.Join(t.TempDir(), test.testReportName+".csv"))
			require.NoError(t, err)
			defer tmpfile.Close()

			err = reporter.Write(tmpfile, test.findings)
			require.NoError(t, err)
			assert.FileExists(t, tmpfile.Name())

			got, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			if test.wantEmpty {
				assert.Empty(t, got)
				return
			}

			want, err := os.ReadFile(test.expected)
			require.NoError(t, err)

			wantStr := lineEndingReplacer.Replace(string(want))
			gotStr := lineEndingReplacer.Replace(string(got))
			assert.Equal(t, wantStr, gotStr)
		})
	}
}
