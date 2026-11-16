package sources

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type countingReader struct {
	reader io.Reader
	reads  int
	bytes  int
}

func (r *countingReader) Read(p []byte) (int, error) {
	r.reads++
	n, err := r.reader.Read(p)
	r.bytes += n
	return n, err
}

func TestIdentifyArchiveFormatSkipsPrefixWhenTraversalDisabled(t *testing.T) {
	content := &countingReader{reader: strings.NewReader("ordinary plain text\n")}
	f := &File{Content: content, Path: "notes.txt", MaxArchiveDepth: 0}

	format, err := f.identifyArchiveFormat(context.Background())
	require.NoError(t, err)
	assert.Nil(t, format)
	assert.Zero(t, content.reads, "disabled archive traversal must not read a prefix of an ordinary file")
}

func TestIdentifyArchiveFormatVerifiesNamedArchiveWhenTraversalDisabled(t *testing.T) {
	// A misleading archive name must still fall back to plain scanning, while
	// a real archive is recognized so the depth check can exclude it.
	data := makeZip(t, "secret.txt", archiveTestSecret)
	f := &File{Content: bytes.NewReader(data), Path: "real.zip", MaxArchiveDepth: 0}

	format, err := f.identifyArchiveFormat(context.Background())
	require.NoError(t, err)
	require.NotNil(t, format)
}

func TestIdentifyArchiveFormatContinuesContentProbeWhenEnabled(t *testing.T) {
	// #2010: a valid archive with an unrelated extension is still recognized
	// when traversal is enabled.
	data := makeZip(t, "secret.txt", archiveTestSecret)
	f := &File{Content: bytes.NewReader(data), Path: "compressed.other", MaxArchiveDepth: 2}

	format, err := f.identifyArchiveFormat(context.Background())
	require.NoError(t, err)
	require.NotNil(t, format)
}

func TestFragmentReadBufferSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want int
	}{
		{name: "unknown", size: 0, want: defaultBufferSize},
		{name: "tiny", size: 12, want: 4 * 1024},
		{name: "small", size: 2 * 1024, want: 4 * 1024},
		{name: "known medium", size: 16 * 1024, want: 16*1024 + 1},
		{name: "at default", size: defaultBufferSize, want: defaultBufferSize},
		{name: "larger", size: 10 * defaultBufferSize, want: defaultBufferSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fragmentReadBufferSize(tt.size))
		})
	}
}

func TestLooksLikeArchiveContent(t *testing.T) {
	assert.True(t, looksLikeArchiveContent(makeZip(t, "a.txt", "x")))
	assert.True(t, looksLikeArchiveContent([]byte("not a real archive but has a NUL\x00")))
	assert.False(t, looksLikeArchiveContent([]byte("ordinary source text with no binary bytes")))
	assert.False(t, looksLikeArchiveContent(nil))
}
