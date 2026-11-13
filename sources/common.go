package sources

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"runtime"

	"github.com/mholt/archives"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/logging"
)

const maxPeekSize = 25 * 1_000 // 25kb

// maxLinePeekSize bounds the extra read used to finish a long line at a
// blank-line boundary. It keeps the normal 25 KB whitespace lookahead while
// allowing lines beyond 125 KB to be scanned without splitting a secret across
// fragments. Memory remains bounded per file; a line that cannot be finished
// within this cap is an operational error rather than a silent split.
const maxLinePeekSize = 1 * 1024 * 1024 // 1 MiB

var isWhitespace [256]bool
var isWindows = runtime.GOOS == "windows"

func init() {
	// define whitespace characters
	isWhitespace[' '] = true
	isWhitespace['\t'] = true
	isWhitespace['\n'] = true
	isWhitespace['\r'] = true
}

// isArchive does a light check to see if the provided path is an archive or
// compressed file. The File source already does this, so this exists mainly
// to avoid expensive calls before sending things to the File source
func isArchive(ctx context.Context, path string) bool {
	format, _, err := archives.Identify(ctx, path, nil)
	return err == nil && format != nil
}

// shouldSkipPath checks a path against all the allowlists to see if it can
// be skipped
func shouldSkipPath(cfg *config.Config, path string) bool {
	if cfg == nil {
		logging.Trace().Str("path", path).Msg("not skipping path because config is nil")
		return false
	}

	for _, a := range cfg.Allowlists {
		if !pathIsSufficientToSkip(a) {
			continue
		}
		if a.PathAllowed(path) ||
			// TODO: Remove this in v9.
			// This is an awkward hack to mitigate https://github.com/gitleaks/gitleaks/issues/1641.
			(isWindows && a.PathAllowed(filepath.ToSlash(path))) {
			return true
		}
	}

	return false
}

// shouldSkipPathAtRoot also evaluates allowlist paths against path relative to
// the directory scan root. A directory scan reports path the way its root was
// supplied, so an anchored pattern written against the root (for example
// `^src/vendor/`) must skip the same subtree for `.`, a relative root, and an
// absolute root. The native-separator form keeps existing Windows allowlist
// patterns working. This only reaches the early-skip decision for allowlists
// whose path alone is sufficient; AND allowlists with content or commit checks
// remain deferred to the detector.
func shouldSkipPathAtRoot(cfg *config.Config, root, path string) bool {
	if cfg == nil || len(cfg.Allowlists) == 0 {
		return false
	}
	if shouldSkipPath(cfg, path) {
		return true
	}
	rel := RootRelativePath(root, path)
	if rel == "" || rel == path {
		return false
	}
	if shouldSkipPath(cfg, rel) {
		return true
	}
	native := filepath.FromSlash(rel)
	return native != rel && shouldSkipPath(cfg, native)
}

// pathIsSufficientToSkip reports whether a path match alone is enough to
// exclude a file without reading its contents.
//
// With the default OR condition, any one matching criterion suppresses the
// finding, so a path match is sufficient. With AND, content and commit checks
// must also be evaluated by the detector. Commits cannot be evaluated during
// enumeration, so a path plus commit allowlist is deferred as well.
func pathIsSufficientToSkip(a *config.Allowlist) bool {
	if a == nil {
		return false
	}
	if a.MatchCondition != config.AllowlistMatchAnd {
		return true
	}
	return len(a.Regexes) == 0 && len(a.StopWords) == 0 && len(a.Commits) == 0
}

// readUntilSafeBoundary consumes |f| until it finds two consecutive `\n` characters, up to |maxPeekSize|.
// This hopefully avoids splitting. (https://github.com/gitleaks/gitleaks/issues/1651)
func readUntilSafeBoundary(r *bufio.Reader, n int, maxPeekSize int, peekBuf *bytes.Buffer) error {
	if peekBuf.Len() == 0 {
		return nil
	}

	// Does the buffer end in consecutive newlines?
	var (
		data         = peekBuf.Bytes()
		lastChar     = data[len(data)-1]
		newlineCount = 0 // Tracks consecutive newlines
	)

	if isWhitespace[lastChar] {
		for i := len(data) - 1; i >= 0; i-- {
			lastChar = data[i]
			if lastChar == '\n' {
				newlineCount++

				// Stop if two consecutive newlines are found
				if newlineCount >= 2 {
					return nil
				}
			} else if isWhitespace[lastChar] {
				// The presence of other whitespace characters (`\r`, ` `, `\t`) shouldn't reset the count.
				// (Intentionally do nothing.)
			} else {
				break
			}
		}
	}

	// If not, read ahead until we (hopefully) find some.
	newlineCount = 0
	for {
		data = peekBuf.Bytes()
		// Check if the last character is a newline.
		lastChar = data[len(data)-1]
		if lastChar == '\n' {
			newlineCount++

			// Stop if two consecutive newlines are found
			if newlineCount >= 2 {
				break
			}
		} else if isWhitespace[lastChar] {
			// The presence of other whitespace characters (`\r`, ` `, `\t`) shouldn't reset the count.
			// (Intentionally do nothing.)
		} else {
			newlineCount = 0 // Reset if a non-newline character is found
		}

		// Stop looking for a blank line at the normal lookahead limit and
		// finish the current line instead, up to the bounded line cap.
		if (peekBuf.Len() - n) >= maxPeekSize {
			return readUntilNewline(r, maxLinePeekSize, peekBuf)
		}

		// Read additional data into a temporary buffer
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		peekBuf.WriteByte(b)
	}
	return nil
}

// readUntilNewline consumes r until the end of the current line or until the
// bounded extra-read budget is exhausted. It is used when the blank-line
// lookahead ends inside a long line; a later chunk must not start mid-secret.
func readUntilNewline(r *bufio.Reader, maxBytes int, peekBuf *bytes.Buffer) error {
	if data := peekBuf.Bytes(); len(data) > 0 && data[len(data)-1] == '\n' {
		return nil
	}
	for i := 0; i < maxBytes; i++ {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		peekBuf.WriteByte(b)
		if b == '\n' {
			return nil
		}
	}
	return fmt.Errorf("line exceeds the %d-byte safe boundary; refusing to report a complete scan", maxBytes)
}
