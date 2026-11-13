package sources

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const archiveTestSecret = "aws_access_key_id = \"AKIAIOSFODNN7EXAMPLE\""

func collectAllFragments(t *testing.T, f *File) (string, error) {
	t.Helper()
	var out strings.Builder
	err := f.Fragments(context.Background(), func(fragment Fragment, err error) error {
		if err != nil {
			return err
		}
		out.WriteString(fragment.Raw)
		return nil
	})
	return out.String(), err
}

func TestFileFragmentsPlaintextWithArchiveLikeName(t *testing.T) {
	names := []string{
		"notes.7z.txt",
		"data.zipx",
		"x.brotli",
		"secret.xzabc",
		"notes.gz.txt",
		"backup.tar.gz.txt",
	}
	content := "header\n" + archiveTestSecret + "\ntrailer\n"
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			f := &File{Content: strings.NewReader(content), Path: name, MaxArchiveDepth: 2}
			raw, err := collectAllFragments(t, f)
			require.NoError(t, err)
			assert.Contains(t, raw, archiveTestSecret)
			assert.Equal(t, content, raw)
		})
	}
}

func TestFileFragmentsShortPlaintextWithArchiveName(t *testing.T) {
	f := &File{Content: strings.NewReader("short"), Path: "short.zip", MaxArchiveDepth: 2}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Equal(t, "short", raw)
}
func TestFileFragmentsNonSeekablePrefixIsPreserved(t *testing.T) {
	content := "prefix-" + archiveTestSecret + "-suffix"
	f := &File{
		Content:         struct{ io.Reader }{Reader: strings.NewReader(content)},
		Path:            "notes.7z.txt",
		MaxArchiveDepth: 2,
	}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Equal(t, content, raw)
}

func TestFileFragmentsSeekFailureFallsBackToPrefixReplay(t *testing.T) {
	content := "prefix-" + archiveTestSecret
	f := &File{
		Content:         &seekErrorReader{Reader: strings.NewReader(content)},
		Path:            "notes.zip",
		MaxArchiveDepth: 2,
	}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Equal(t, content, raw)
}

func TestFileFragmentsArchiveVerificationReadErrorIsReturned(t *testing.T) {
	f := &File{
		Content: &afterDataErrorReader{
			data: []byte("prefix-" + archiveTestSecret),
			err:  errors.New("synthetic verification read failure"),
		},
		Path:            "notes.zip",
		MaxArchiveDepth: 2,
	}
	_, err := collectAllFragments(t, f)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "archive verification")
	assert.Contains(t, err.Error(), "synthetic verification read failure")
}

func TestFileFragmentsActualArchivesStillWork(t *testing.T) {
	t.Run("zip", func(t *testing.T) {
		data := makeZip(t, "secret.txt", archiveTestSecret)
		f := &File{Content: bytes.NewReader(data), Path: "real.zip", MaxArchiveDepth: 2}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Contains(t, raw, archiveTestSecret)
	})

	t.Run("non-seekable zip", func(t *testing.T) {
		data := makeZip(t, "secret.txt", archiveTestSecret)
		f := &File{
			Content:         struct{ io.Reader }{Reader: bytes.NewReader(data)},
			Path:            "real.zip",
			MaxArchiveDepth: 2,
		}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Contains(t, raw, archiveTestSecret)
	})

	t.Run("gzip", func(t *testing.T) {
		data := makeGzip(t, []byte(archiveTestSecret))
		f := &File{Content: bytes.NewReader(data), Path: "real.gz", MaxArchiveDepth: 2}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Contains(t, raw, archiveTestSecret)
	})

	t.Run("tar", func(t *testing.T) {
		data := makeTar(t, "secret.txt", archiveTestSecret)
		f := &File{Content: bytes.NewReader(data), Path: "real.tar", MaxArchiveDepth: 2}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Contains(t, raw, archiveTestSecret)
	})

	t.Run("tar.gz", func(t *testing.T) {
		data := makeTarGz(t, "secret.txt", archiveTestSecret)
		f := &File{Content: bytes.NewReader(data), Path: "real.tar.gz", MaxArchiveDepth: 2}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Contains(t, raw, archiveTestSecret)
	})
}

func TestFileFragmentsCorruptArchiveIsOperationalError(t *testing.T) {
	// A valid zip signature with a truncated body must route to the extractor
	// and report the extraction failure, not silently fall back to plaintext.
	f := &File{
		Content:         bytes.NewReader([]byte("PK\x03\x04truncated")),
		Path:            "broken.zip",
		MaxArchiveDepth: 2,
	}
	_, err := collectAllFragments(t, f)
	require.Error(t, err)
}

func TestFileFragmentsArchiveDepthDisabledStillSkips(t *testing.T) {
	f := &File{
		Content:         bytes.NewReader(makeZip(t, "secret.txt", archiveTestSecret)),
		Path:            "real.zip",
		MaxArchiveDepth: 0,
	}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Empty(t, raw)
}

func TestFileSizeLimitPlainKnownAndUnknown(t *testing.T) {
	t.Run("known over limit", func(t *testing.T) {
		f := &File{Content: strings.NewReader("x"), Path: "big.txt", Size: 101, MaxFileSize: 100}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Empty(t, raw)
	})
	t.Run("known at limit", func(t *testing.T) {
		content := strings.Repeat("x", 100)
		f := &File{Content: strings.NewReader(content), Path: "limit.txt", Size: 100, MaxFileSize: 100}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Equal(t, content, raw)
	})
	t.Run("unknown over limit", func(t *testing.T) {
		f := &File{Content: strings.NewReader(strings.Repeat("x", 101)), Path: "unknown.txt", MaxFileSize: 100}
		raw, err := collectAllFragments(t, f)
		require.NoError(t, err)
		assert.Empty(t, raw)
	})
}

func TestArchiveMemberSizeLimit(t *testing.T) {
	zipData := makeZip(t, "big.txt", strings.Repeat("A", 200))
	f := &File{Content: bytes.NewReader(zipData), Path: "members.zip", MaxArchiveDepth: 2, MaxFileSize: 100}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Empty(t, raw)

	zipSmall := makeZip(t, "small.txt", archiveTestSecret)
	f = &File{Content: bytes.NewReader(zipSmall), Path: "small.zip", MaxArchiveDepth: 2, MaxFileSize: 1_000}
	raw, err = collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Contains(t, raw, archiveTestSecret)
}

func TestCompressedStreamSizeLimit(t *testing.T) {
	over := makeGzip(t, []byte(strings.Repeat("A", 200)))
	f := &File{Content: bytes.NewReader(over), Path: "over.gz", MaxArchiveDepth: 2, MaxFileSize: 100}
	raw, err := collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Empty(t, raw)

	under := makeGzip(t, []byte(archiveTestSecret))
	f = &File{Content: bytes.NewReader(under), Path: "under.gz", MaxArchiveDepth: 2, MaxFileSize: 1_000}
	raw, err = collectAllFragments(t, f)
	require.NoError(t, err)
	assert.Contains(t, raw, archiveTestSecret)
}

func TestFilesChecksArchiveMembersNotOuterSize(t *testing.T) {
	dir := t.TempDir()
	archivePath := filepath.Join(dir, "members.zip")
	require.NoError(t, os.WriteFile(archivePath, makeZip(t, "secret.txt", archiveTestSecret), 0o600))
	info, err := os.Stat(archivePath)
	require.NoError(t, err)
	require.Greater(t, info.Size(), int64(100), "fixture archive must exceed the outer size limit")

	source := Files{
		Path:            archivePath,
		Sema:            semgroup.NewGroup(context.Background(), 1),
		MaxArchiveDepth: 2,
		MaxFileSize:     100,
	}
	var raw strings.Builder
	err = source.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		raw.WriteString(fragment.Raw)
		return nil
	})
	require.NoError(t, err)
	assert.Contains(t, raw.String(), archiveTestSecret)
}

type seekErrorReader struct {
	*strings.Reader
}

func (r *seekErrorReader) Seek(int64, int) (int64, error) {
	return 0, errors.New("synthetic seek failure")
}

type afterDataErrorReader struct {
	data []byte
	err  error
}

func (r *afterDataErrorReader) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		return n, nil
	}
	return 0, r.err
}

func makeZip(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	require.NoError(t, err)
	_, err = io.WriteString(w, content)
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func makeGzip(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(content)
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func makeTar(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Size: int64(len(content)), Mode: 0o600}))
	_, err := io.WriteString(tw, content)
	require.NoError(t, err)
	require.NoError(t, tw.Close())
	return buf.Bytes()
}

func makeTarGz(t *testing.T, name, content string) []byte {
	t.Helper()
	var compressed bytes.Buffer
	gw := gzip.NewWriter(&compressed)
	_, err := gw.Write(makeTar(t, name, content))
	require.NoError(t, err)
	require.NoError(t, gw.Close())
	return compressed.Bytes()
}

func TestFileSizeLimitReadErrorIsReturned(t *testing.T) {
	f := &File{
		Content:         &afterDataErrorReader{data: []byte(strings.Repeat("A", 20)), err: errors.New("synthetic spool read failure")},
		Path:            "stream.txt",
		MaxFileSize:     100,
		MaxArchiveDepth: 0,
	}
	_, err := collectAllFragments(t, f)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "synthetic spool read failure")
}
