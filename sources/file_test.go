package sources

import (
	"bufio"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collectFileFragments(t *testing.T, content string) []string {
	t.Helper()
	file := &File{Path: "input.txt"}
	var fragments []string
	err := file.fileFragments(
		context.Background(),
		bufio.NewReader(strings.NewReader(content)),
		func(fragment Fragment, err error) error {
			require.NoError(t, err)
			fragments = append(fragments, fragment.Raw)
			return nil
		},
	)
	require.NoError(t, err)
	return fragments
}

func TestFileFragmentsLongLineKeepsSecretIntact(t *testing.T) {
	const secret = "AKIAQ7Z3R5T7Y9U2I4O6"

	tests := map[string]int{
		"100 KB boundary": 100_000 - 5,
		"125 KB boundary": 125_000 - 5,
		"125 KB plus one": 125_001 - 5,
		"150 KB boundary": 150_000 - 5,
		"250 KB boundary": 250_000 - 5,
	}
	for name, prefixLen := range tests {
		t.Run(name, func(t *testing.T) {
			content := strings.Repeat("A", prefixLen) + secret + strings.Repeat("A", 1_000) + "\n"
			fragments := collectFileFragments(t, content)
			require.Equal(t, len(content), len(strings.Join(fragments, "")))
			foundIntact := false
			for _, fragment := range fragments {
				if strings.Contains(fragment, secret) {
					foundIntact = true
					break
				}
			}
			require.True(t, foundIntact, "the complete secret must be present in at least one fragment")
		})
	}
}

func TestFileFragmentsOverLongLineIsOperationalError(t *testing.T) {
	content := strings.Repeat("A", maxLinePeekSize+200_000) + "\n"
	file := &File{Path: "input.txt"}
	var fragments []string
	err := file.fileFragments(
		context.Background(),
		bufio.NewReader(strings.NewReader(content)),
		func(fragment Fragment, err error) error {
			if err != nil {
				return err
			}
			fragments = append(fragments, fragment.Raw)
			return nil
		},
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "safe boundary")
	assert.NotEmpty(t, fragments)
}

func TestFileFragmentsMIMEIsFirstChunkOnly(t *testing.T) {
	const secret = "aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\""

	tests := []struct {
		name      string
		prefixLen int
		lineEnd   string
	}{
		{name: "100 KB boundary control", prefixLen: 100_000, lineEnd: "\n"},
		{name: "125 KB second chunk boundary", prefixLen: 125_000, lineEnd: "\n"},
		{name: "125 KB plus one byte", prefixLen: 125_001, lineEnd: "\n"},
		{name: "250 KB third chunk boundary", prefixLen: 250_000, lineEnd: "\n"},
		{name: "300 KB fourth chunk boundary", prefixLen: 300_000, lineEnd: "\n"},
		{name: "125 KB with CRLF", prefixLen: 125_000, lineEnd: "\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := strings.Repeat("A", test.prefixLen) +
				"MZ" +
				strings.Repeat("A", 1_000) +
				test.lineEnd +
				secret +
				test.lineEnd

			fragments := collectFileFragments(t, content)
			joined := strings.Join(fragments, "")
			require.Equal(t, len(content), len(joined), "the entire text file must be yielded")
			require.Contains(t, joined, secret, "a later binary-looking signature must not reclassify text")
		})
	}

	t.Run("binary signature at the start is still skipped", func(t *testing.T) {
		content := "MZ" + strings.Repeat("A", 1_000) + "\n" + secret + "\n"
		fragments := collectFileFragments(t, content)
		require.Empty(t, fragments)
	})
}
