package sources

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/require"
)

// fragmentsOf returns the raw content of each fragment yielded for content.
func fragmentsOf(t *testing.T, content string) []string {
	t.Helper()

	var raws []string
	file := &File{Path: "file.txt"}
	err := file.fileFragments(
		context.Background(),
		bufio.NewReader(strings.NewReader(content)),
		func(fragment Fragment, err error) error {
			require.NoError(t, err)
			raws = append(raws, fragment.Raw)
			return nil
		},
	)
	require.NoError(t, err)

	return raws
}

func Test_fileFragments_filetype(t *testing.T) {
	// The first line spans more than one chunk, so no newline has been read
	// yet when the second chunk starts.
	firstLine := strings.Repeat("a", 2*defaultBufferSize)
	plain := fragmentsOf(t, firstLine+"\nafter\n")
	require.Greater(t, len(plain), 1, "the first line must span more than one chunk")
	boundary := len(plain[0])

	t.Run("binary signature at the start of a later chunk", func(t *testing.T) {
		content := firstLine[:boundary] + "MZ" + firstLine[boundary+2:] + "\nafter\n"
		got := strings.Join(fragmentsOf(t, content), "")

		// Compare lengths first to keep the failure message short.
		require.Equal(t, len(content), len(got), "the fragments must cover the whole file")
		require.True(t, content == got, "the fragments must match the file content")
	})

	t.Run("binary signature at the start of the file", func(t *testing.T) {
		content := "MZ" + firstLine[2:] + "\nafter\n"

		require.Empty(t, fragmentsOf(t, content))
	})
}

func TestArchiveMemberHonorsMaxFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "large.zip")
	require.NoError(t, writeZipFile(archivePath, "large.txt", bytes.Repeat([]byte("A"), 2_000)))

	source := Files{
		Path:            archivePath,
		Sema:            semgroup.NewGroup(context.Background(), 1),
		MaxArchiveDepth: 1,
		MaxFileSize:     1_000,
	}

	var fragments []Fragment
	err := source.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	require.Empty(t, fragments)
}

func TestCompressedFileHonorsMaxFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	gzipPath := filepath.Join(tmpDir, "large.txt.gz")
	require.NoError(t, writeGzipFile(gzipPath, bytes.Repeat([]byte("A"), 2_000)))

	source := Files{
		Path:            gzipPath,
		Sema:            semgroup.NewGroup(context.Background(), 1),
		MaxArchiveDepth: 1,
		MaxFileSize:     1_000,
	}

	var fragments []Fragment
	err := source.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	require.Empty(t, fragments)
}

func writeZipFile(path, name string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	w, err := zw.Create(name)
	if err != nil {
		_ = zw.Close()
		return err
	}
	if _, err := w.Write(content); err != nil {
		_ = zw.Close()
		return err
	}
	return zw.Close()
}

func writeGzipFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	if _, err := gw.Write(content); err != nil {
		_ = gw.Close()
		return err
	}
	return gw.Close()
}
