package sources

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/mholt/archives"
	"github.com/rs/zerolog"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/logging"
)

const defaultBufferSize = 100 * 1_000 // 100kb
const InnerPathSeparator = "!"

// archivePeekSize bounds how much content is inspected when verifying a
// filename-derived archive match. It covers the magic/header checks used by the
// archive package, including brotli's bounded trial decode.
const archivePeekSize = 64 * 1_000

type seekReaderAt interface {
	io.ReaderAt
	io.Seeker
}

// File is a source for yielding fragments from a file or other reader
type File struct {
	// Content provides a reader to the file's content
	Content io.Reader
	// Path is the resolved real path of the file
	Path string
	// Symlink represents a symlink to the file if that's how it was discovered
	Symlink string
	// Buffer is used for reading the content in chunks
	Buffer []byte
	// Config is the gitleaks config used for shouldSkipPath. If not set, then
	// shouldSkipPath is ignored
	Config *config.Config
	// MaxFileSize is the largest uncompressed file or member size, in bytes,
	// that will be scanned. Zero means no explicit limit. Propagation follows
	// gitleaks/gitleaks#2099 (@alephnull-sh); unknown-size streams are bounded
	// with a temporary spool instead of buffering whole contents in memory.
	MaxFileSize int
	// Size is the known size in bytes. For archive members this is the
	// uncompressed member size. Zero means unknown or empty; a positive
	// MaxFileSize is then enforced by bounded spooling.
	Size int64
	// outerPaths is the list of container paths (e.g. archives) that lead to
	// this file
	outerPaths []string
	// MaxArchiveDepth limits how deep the sources will explore nested archives
	MaxArchiveDepth int
	// archiveDepth is the current archive nesting depth
	archiveDepth int
	// Informative is the optional coverage collector. Nil disables collection.
	Informative *informative.Collector
	// Target is the optional coverage-record handle. Files creates it during
	// enumeration; direct File users may leave it nil.
	Target *informative.Target
	// TargetKind overrides the inferred target kind for direct File users
	// (for example, stdin versus an unnamed stream).
	TargetKind informative.Kind
	// Commit is the Git commit associated with this content, when applicable.
	// It is reporting metadata only.
	Commit string
}

// ensureTarget creates a target for direct File users that did not enumerate
// through Files. Enumeration paths pass an already-observed target.
func (s *File) ensureTarget() *informative.Target {
	if s.Target != nil || s.Informative == nil {
		return s.Target
	}
	kind := s.TargetKind
	if kind == "" {
		kind = informative.KindFile
		switch {
		case s.archiveDepth > 0:
			kind = informative.KindArchiveMember
		case s.Path == "":
			kind = informative.KindStream
		}
	}
	s.Target = s.Informative.Track(kind, s.FullPath(), s.Commit, s.Size, s.Size > 0)
	return s.Target
}

// Fragments yields fragments for the this source
func (s *File) Fragments(ctx context.Context, yield FragmentsFunc) error {
	s.ensureTarget()
	format, _, err := archives.Identify(ctx, s.Path, nil)
	if err != nil && !errors.Is(err, archives.NoMatch) {
		return fmt.Errorf("could not identify file format: %w", err)
	}
	// Identify was called without content, so a filename alone can select a
	// format. Verify the selected format against a bounded content prefix
	// before routing the file into archive/decompression handling.
	if err == nil && format != nil {
		matches, matchErr := s.contentMatchesArchive(ctx, format)
		if matchErr != nil {
			return matchErr
		}
		if !matches {
			format = nil
		}
	}
	if err == nil && format != nil {
		if s.archiveDepth+1 > s.MaxArchiveDepth {
			var event *zerolog.Event

			// Warn if the feature is enabled; else emit a trace log.
			if s.MaxArchiveDepth != 0 {
				event = logging.Warn()
			} else {
				event = logging.Trace()
			}

			event.Str(
				"path", s.FullPath(),
			).Int(
				"max_archive_depth", s.MaxArchiveDepth,
			).Msg("skipping archive: exceeds max archive depth")

			s.Target.Exclude(informative.ReasonArchiveDepth, informative.ExclusionOptions{
				Detail: fmt.Sprintf("maximum archive depth %d", s.MaxArchiveDepth),
			})
			return nil
		}
		if extractor, ok := format.(archives.Extractor); ok {
			return s.extractorFragments(ctx, extractor, s.Content, yield)
		}
		if decompressor, ok := format.(archives.Decompressor); ok {
			return s.decompressorFragments(ctx, decompressor, s.Content, yield)
		}
		logging.Warn().Str("path", s.FullPath()).Msg("skipping unknown archive type")
	}

	cleanup, skipped, err := s.applyPlainSizeLimit(ctx)
	if err != nil {
		return err
	}
	if skipped {
		return nil
	}
	if cleanup != nil {
		defer cleanup()
	}
	return s.fileFragments(ctx, bufio.NewReader(s.Content), yield)
}

// contentMatchesArchive reports whether the selected format's stream signature
// matches a bounded prefix. Consumed bytes are restored before returning, even
// on a read error, so the caller can report the error without losing content.
// The content-verification approach is adapted from gitleaks/gitleaks#2317
// (@Varshith-Kali), with explicit read-error propagation and non-seekable
// prefix replay added here.
func (s *File) contentMatchesArchive(ctx context.Context, format archives.Format) (bool, error) {
	peeked := make([]byte, archivePeekSize)
	n, readErr := io.ReadFull(s.Content, peeked)
	peeked = peeked[:n]
	s.restorePrefix(peeked)

	if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
		return false, fmt.Errorf("could not read file for archive verification: %w", readErr)
	}

	match, err := verifyArchiveStream(ctx, format, s.Path, peeked)
	if err != nil {
		return false, fmt.Errorf("could not verify archive content: %w", err)
	}
	return match.ByStream, nil
}

func (s *File) restorePrefix(prefix []byte) {
	if len(prefix) == 0 {
		return
	}
	if seeker, ok := s.Content.(io.Seeker); ok {
		if _, err := seeker.Seek(-int64(len(prefix)), io.SeekCurrent); err == nil {
			return
		}
	}
	s.Content = io.MultiReader(bytes.NewReader(prefix), s.Content)
}

// verifyArchiveStream checks the format's stream signature against a bounded
// prefix. Combined formats are checked at the outer compression layer;
// verifying the inner archive would require unbounded decompression.
func verifyArchiveStream(ctx context.Context, format archives.Format, name string, peeked []byte) (archives.MatchResult, error) {
	if ca, ok := format.(archives.CompressedArchive); ok && ca.Compression != nil {
		return ca.Compression.Match(ctx, name, bytes.NewReader(peeked))
	}
	return format.Match(ctx, name, bytes.NewReader(peeked))
}

// applyPlainSizeLimit enforces MaxFileSize for non-archive content. Known
// sizes are checked without reading. Unknown or zero sizes are bounded through
// a temporary spool so an oversized stream is excluded before any fragment is
// yielded.
func (s *File) applyPlainSizeLimit(ctx context.Context) (cleanup func(), skipped bool, err error) {
	if s.MaxFileSize <= 0 {
		return nil, false, nil
	}

	if s.Size > int64(s.MaxFileSize) {
		s.excludeSizeLimit(s.Size, int64(s.MaxFileSize), true)
		logging.Warn().
			Str("path", s.FullPath()).
			Msgf("skipping file: too large max_size=%dMB, size=%dMB", s.MaxFileSize/1_000_000, s.Size/1_000_000)
		return nil, true, nil
	}
	if s.Size > 0 {
		return nil, false, nil
	}

	content, size, cleanup, skipped, err := spoolLimitedReader(ctx, s.Content, s.MaxFileSize)
	if err != nil {
		return nil, false, err
	}
	if skipped {
		cleanup()
		s.excludeSizeLimit(size, int64(s.MaxFileSize), false)
		logging.Warn().
			Str("path", s.FullPath()).
			Msgf("skipping file: too large max_size=%dMB", s.MaxFileSize/1_000_000)
		return nil, true, nil
	}
	s.Content = content
	s.Size = size
	s.Target.SetSize(size, true)
	return cleanup, false, nil
}

func (s *File) excludeSizeLimit(size, limit int64, exactSize bool) {
	options := informative.ExclusionOptions{
		LimitBytes: &limit,
		Detail:     fmt.Sprintf("maximum file size %d bytes", limit),
	}
	if exactSize {
		options.SizeBytes = &size
	} else {
		options.OmitSize = true
		options.Detail = fmt.Sprintf("maximum file size %d bytes; stopped after %d bytes", limit, size)
	}
	s.Target.Exclude(informative.ReasonSizeLimit, options)
}

func int64Ptr(value int64) *int64 {
	return &value
}

func spoolLimitedReader(ctx context.Context, reader io.Reader, max int) (io.Reader, int64, func(), bool, error) {
	if max <= 0 {
		return reader, 0, nil, false, nil
	}
	tmp, err := os.CreateTemp("", "gitleaks-content-")
	if err != nil {
		return nil, 0, nil, false, fmt.Errorf("could not create temporary content file: %w", err)
	}
	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}

	n, err := io.Copy(tmp, io.LimitReader(reader, int64(max)+1))
	if err != nil {
		cleanup()
		return nil, 0, nil, false, fmt.Errorf("could not read file: %w", err)
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		cleanup()
		return nil, 0, nil, false, ctxErr
	}
	if n > int64(max) {
		cleanup()
		return nil, n, func() {}, true, nil
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		cleanup()
		return nil, 0, nil, false, fmt.Errorf("could not rewind temporary content file: %w", err)
	}
	return tmp, n, cleanup, false, nil
}

// extractorFragments recursively crawls archives and yields fragments
func (s *File) extractorFragments(ctx context.Context, extractor archives.Extractor, reader io.Reader, yield FragmentsFunc) error {
	if _, isSeekReaderAt := reader.(seekReaderAt); !isSeekReaderAt {
		switch extractor.(type) {
		case archives.SevenZip, archives.Zip:
			tmpfile, err := os.CreateTemp("", "gitleaks-archive-")
			if err != nil {
				return fmt.Errorf("could not create temporary archive file: %w", err)
			}
			cleanup := func() {
				_ = tmpfile.Close()
				_ = os.Remove(tmpfile.Name())
			}
			defer cleanup()

			var (
				copied  int64
				copyErr error
			)
			if s.MaxFileSize > 0 {
				copied, copyErr = io.Copy(tmpfile, io.LimitReader(reader, int64(s.MaxFileSize)+1))
			} else {
				copied, copyErr = io.Copy(tmpfile, reader)
			}
			if copyErr != nil {
				return fmt.Errorf("could not copy archive file: %w", copyErr)
			}
			if s.MaxFileSize > 0 && copied > int64(s.MaxFileSize) {
				s.excludeSizeLimit(copied, int64(s.MaxFileSize), false)
				logging.Warn().
					Str("path", s.FullPath()).
					Msgf("skipping archive: too large max_size=%dMB", s.MaxFileSize/1_000_000)
				return nil
			}
			if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("could not rewind archive file: %w", err)
			}
			reader = tmpfile
		}
	}

	return extractor.Extract(ctx, reader, func(_ context.Context, d archives.FileInfo) error {
		if d.IsDir() {
			return nil
		}

		path := filepath.Clean(d.NameInArchive)
		memberPath := s.FullPath() + InnerPathSeparator + filepath.ToSlash(path)
		memberTarget := s.Informative.Track(
			informative.KindArchiveMember,
			memberPath,
			s.Commit,
			d.Size(),
			true,
		)
		if s.Config != nil && shouldSkipPath(s.Config, path) {
			memberTarget.Exclude(informative.ReasonPathAllowlist, informative.ExclusionOptions{
				Detail: "global allowlist",
			})
			logging.Debug().Str("path", s.FullPath()).Msg("skipping file: global allowlist")
			return nil
		}

		if s.MaxFileSize > 0 && d.Size() > int64(s.MaxFileSize) {
			memberTarget.Exclude(informative.ReasonSizeLimit, informative.ExclusionOptions{
				LimitBytes: int64Ptr(int64(s.MaxFileSize)),
				Detail:     fmt.Sprintf("maximum archive member size %d bytes", s.MaxFileSize),
			})
			logging.Warn().
				Str("path", filepath.ToSlash(filepath.Join(s.FullPath(), path))).
				Msgf("skipping file: too large max_size=%dMB, size=%dMB", s.MaxFileSize/1_000_000, d.Size()/1_000_000)
			return nil
		}
		if d.Size() == 0 {
			// Keep archive opening behavior unchanged. An empty member yields
			// no fragments, but any archive-level open error still propagates.
			memberTarget.Exclude(informative.ReasonEmptyFile, informative.ExclusionOptions{
				Detail: "empty archive member",
			})
		}

		innerReader, err := d.Open()
		if err != nil {
			return fmt.Errorf("could not open archive entry %q: %w", d.NameInArchive, err)
		}
		file := &File{
			Content:         innerReader,
			Path:            path,
			Symlink:         s.Symlink,
			outerPaths:      append(s.outerPaths, filepath.ToSlash(s.Path)),
			MaxArchiveDepth: s.MaxArchiveDepth,
			MaxFileSize:     s.MaxFileSize,
			Size:            d.Size(),
			archiveDepth:    s.archiveDepth + 1,
			Informative:     s.Informative,
			Target:          memberTarget,
			Commit:          s.Commit,
		}

		if err := file.Fragments(ctx, yield); err != nil {
			_ = innerReader.Close()
			return err
		}
		if err := innerReader.Close(); err != nil {
			return fmt.Errorf("could not close archive entry %q: %w", d.NameInArchive, err)
		}
		return nil
	})
}

// decompressorFragments recursively crawls archives and yields fragments
func (s *File) decompressorFragments(ctx context.Context, decompressor archives.Decompressor, reader io.Reader, yield FragmentsFunc) error {
	innerReader, err := decompressor.OpenReader(reader)
	if err != nil {
		return fmt.Errorf("could not open compressed file: %w", err)
	}

	if s.MaxFileSize > 0 {
		content, size, cleanup, skipped, err := spoolLimitedReader(ctx, innerReader, s.MaxFileSize)
		if closeErr := innerReader.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("could not close compressed file: %w", closeErr)
		}
		if err != nil {
			return err
		}
		if skipped {
			cleanup()
			s.excludeSizeLimit(size, int64(s.MaxFileSize), false)
			logging.Warn().
				Str("path", s.FullPath()).
				Msgf("skipping file: too large max_size=%dMB", s.MaxFileSize/1_000_000)
			return nil
		}
		defer cleanup()
		s.Content = content
		s.Size = size
		s.Target.SetSize(size, true)
		return s.fileFragments(ctx, bufio.NewReader(content), yield)
	}

	if err := s.fileFragments(ctx, bufio.NewReader(innerReader), yield); err != nil {
		_ = innerReader.Close()
		return err
	}
	if err := innerReader.Close(); err != nil {
		return fmt.Errorf("could not close compressed file: %w", err)
	}
	return nil
}

// fileFragments reads the file into fragments to yield
func (s *File) fileFragments(ctx context.Context, reader *bufio.Reader, yield FragmentsFunc) error {
	// Create a buffer if the caller hasn't provided one
	if s.Buffer == nil {
		s.Buffer = make([]byte, defaultBufferSize)
	}

	totalLines := 0
	isFirstChunk := true
	var readErr error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fragment := Fragment{
				FilePath:          s.FullPath(),
				InformativeTarget: s.Target,
			}

			n, err := reader.Read(s.Buffer)
			if n == 0 {
				if err != nil && err != io.EOF {
					if readErr == nil {
						readErr = fmt.Errorf("could not read file: %w", err)
					}
					return readErr
				}
				return readErr
			}

			// Only check the filetype on the first chunk. totalLines stays zero
			// while a long first line spans chunks, so it cannot identify the
			// start of the file.
			if isFirstChunk {
				isFirstChunk = false
				// TODO: could other optimizations be introduced here?
				if mimetype, err := filetype.Match(s.Buffer[:n]); err != nil {
					return yield(
						fragment,
						fmt.Errorf("could not read file: could not determine type: %w", err),
					)
				} else if mimetype.MIME.Type == "application" {
					s.Target.Exclude(informative.ReasonBinaryContent, informative.ExclusionOptions{
						Detail: "MIME type " + mimetype.MIME.Value,
					})
					logging.Debug().
						Str("mime_type", mimetype.MIME.Value).
						Str("path", s.FullPath()).
						Msgf("skipping binary file")

					return nil
				}
			}

			// Try to split chunks across large areas of whitespace, if possible.
			peekBuf := bytes.NewBuffer(s.Buffer[:n])
			if boundaryErr := readUntilSafeBoundary(reader, n, maxPeekSize, peekBuf); boundaryErr != nil {
				if readErr == nil {
					readErr = fmt.Errorf("could not read file: could not read until safe boundary: %w", boundaryErr)
				}
				logging.Warn().Err(boundaryErr).Msgf("issue reading file")
			}

			fragment.Raw = peekBuf.String()
			fragment.Bytes = peekBuf.Bytes()
			fragment.StartLine = totalLines + 1

			// Count the number of newlines in this chunk
			totalLines += strings.Count(fragment.Raw, "\n")

			if len(s.Symlink) > 0 {
				fragment.SymlinkFile = s.Symlink
			}

			if isWindows {
				fragment.FilePath = filepath.ToSlash(fragment.FilePath)
				fragment.SymlinkFile = filepath.ToSlash(s.Symlink)
				fragment.WindowsFilePath = s.FullPath()
			}

			if err != nil && err != io.EOF {
				if readErr == nil {
					readErr = fmt.Errorf("could not read file: %w", err)
				}
				logging.Warn().Err(err).Msgf("issue reading file")
			}

			// Scan the content that was successfully read even if the reader
			// reported an error alongside it; the error is returned below.
			if yieldErr := yield(fragment, nil); yieldErr != nil {
				return yieldErr
			}
			if err == io.EOF {
				return readErr
			}
		}
	}
}

// FullPath returns the File.Path with any preceding outer paths
func (s *File) FullPath() string {
	if len(s.outerPaths) > 0 {
		return strings.Join(
			// outerPaths have already been normalized to slash
			append(s.outerPaths, s.Path),
			InnerPathSeparator,
		)
	}

	return s.Path
}
