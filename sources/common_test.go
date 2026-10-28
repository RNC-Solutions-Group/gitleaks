package sources

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func Test_readUntilSafeBoundary(t *testing.T) {
	// Arrange
	cases := []struct {
		name     string
		r        io.Reader
		expected string
	}{
		// Current split is fine, exit early.
		{
			name:     "safe original split - LF",
			r:        strings.NewReader("abc\n\ndefghijklmnop\n\nqrstuvwxyz"),
			expected: "abc\n\n",
		},
		{
			name:     "safe original split - CRLF",
			r:        strings.NewReader("a\r\n\r\nbcdefghijklmnop\n"),
			expected: "a\r\n\r\n",
		},
		// Current split is bad, look for a better one.
		{
			name:     "safe split - LF",
			r:        strings.NewReader("abcdefg\nhijklmnop\n\nqrstuvwxyz"),
			expected: "abcdefg\nhijklmnop\n\n",
		},
		{
			name:     "safe split - CRLF",
			r:        strings.NewReader("abcdefg\r\nhijklmnop\r\n\r\nqrstuvwxyz"),
			expected: "abcdefg\r\nhijklmnop\r\n\r\n",
		},
		{
			name:     "safe split - blank line",
			r:        strings.NewReader("abcdefg\nhijklmnop\n\t  \t\nqrstuvwxyz"),
			expected: "abcdefg\nhijklmnop\n\t  \t\n",
		},
		// Current split is bad and there is no blank line, finish the current line.
		{
			name:     "no blank line - finish the line",
			r:        strings.NewReader("abcdefg\nhijklmnopqrstuvwxyz\n0123"),
			expected: "abcdefg\nhijklmnopqrstuvwxyz\n",
		},
		{
			name:     "no blank line - finish the line at EOF",
			r:        strings.NewReader("abcdefg\nhijklmnopqrstuvwxyz"),
			expected: "abcdefg\nhijklmnopqrstuvwxyz",
		},
		// Current split is bad, exhaust options.
		{
			name:     "no safe split",
			r:        strings.NewReader("abcdefg\nhijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			expected: "abcdefg\nhijklmnopqrstuvwxyz0123456789ABCDEFGH",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := make([]byte, 5)
			n, err := c.r.Read(buf)
			require.NoError(t, err)

			// Act
			reader := bufio.NewReader(c.r)
			peekBuf := bytes.NewBuffer(buf[:n])
			err = readUntilSafeBoundary(reader, n, 20, peekBuf)
			require.NoError(t, err)

			// Assert
			t.Log(peekBuf.String())
			require.Equal(t, c.expected, peekBuf.String())
		})
	}
}

func TestShouldSkipPath(t *testing.T) {
	mustCompile := func(pattern string) []*regexp.Regexp {
		return []*regexp.Regexp{regexp.MustCompile(pattern)}
	}

	tests := map[string]struct {
		allowlist *config.Allowlist
		path      string
		want      bool
	}{
		"or path matches": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchOr,
				Paths:          mustCompile(`\.md$`),
			},
			path: "docs/README.md",
			want: true,
		},
		"or path does not match": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchOr,
				Paths:          mustCompile(`\.md$`),
			},
			path: "main.go",
			want: false,
		},
		"or path alone is sufficient with regexes": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchOr,
				Paths:          mustCompile(`\.md$`),
				Regexes:        mustCompile(`PLACEHOLDER`),
			},
			path: "docs/README.md",
			want: true,
		},
		"and path with nonmatching regex must be evaluated": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchAnd,
				Paths:          mustCompile(`\.md$`),
				Regexes:        mustCompile(`PLACEHOLDER`),
			},
			path: "docs/README.md",
			want: false,
		},
		"and path with stopword must be evaluated": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchAnd,
				Paths:          mustCompile(`\.md$`),
				StopWords:      []string{"example"},
			},
			path: "docs/README.md",
			want: false,
		},
		"and path with commit must be evaluated": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchAnd,
				Paths:          mustCompile(`\.md$`),
				Commits:        []string{"0000000000000000000000000000000000000000"},
			},
			path: "docs/README.md",
			want: false,
		},
		"and path only is sufficient": {
			allowlist: &config.Allowlist{
				MatchCondition: config.AllowlistMatchAnd,
				Paths:          mustCompile(`\.md$`),
			},
			path: "docs/README.md",
			want: true,
		},
		"nil allowlist is safe": {
			allowlist: nil,
			path:      "docs/README.md",
			want:      false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.allowlist != nil {
				require.NoError(t, test.allowlist.Validate())
			}
			cfg := &config.Config{Allowlists: []*config.Allowlist{test.allowlist}}
			assert.Equal(t, test.want, shouldSkipPath(cfg, test.path))
		})
	}

	t.Run("nil config", func(t *testing.T) {
		assert.False(t, shouldSkipPath(nil, "any/path.md"))
	})
}
