package sources

import (
	"path/filepath"
	"strings"

	"github.com/zricethezav/gitleaks/v8/internal/informative"
)

// Fragment represents a fragment of a source with its meta data
type Fragment struct {
	// Raw is the raw content of the fragment
	Raw string

	Bytes []byte

	// FilePath is the path to the file if applicable.
	// The path separator MUST be normalized to `/`.
	FilePath    string
	SymlinkFile string
	// WindowsFilePath is the path with the original separator.
	// This provides a backwards-compatible solution to https://github.com/gitleaks/gitleaks/issues/1565.
	WindowsFilePath string `json:"-"` // TODO: remove this in v9.

	// Root is the directory a directory scan was rooted at, spelled the way the
	// caller supplied it. FilePath keeps that same spelling for reporting, so
	// path-based matching also evaluates the path relative to Root. It is empty
	// for Git scans, stdin, and single-file targets.
	//
	// Root-relative directory matching was proposed in gitleaks/gitleaks#2248
	// (@hamodywe); this implementation keeps the reported path and fingerprint
	// contract unchanged and handles archive members lexically.
	Root string `json:"-"`

	// CommitSHA is the SHA of the commit if applicable
	CommitSHA string // TODO: remove this in v9 and use CommitInfo instead

	// StartLine is the line number this fragment starts on
	StartLine int

	// CommitInfo captures additional information about the git commit if applicable
	CommitInfo *CommitInfo

	InheritedFromFinding bool // Indicates if this fragment is inherited from a finding

	// InformativeTarget is the optional coverage-record handle for this
	// fragment's target. It carries reporting metadata only and is never part
	// of finding output.
	InformativeTarget *informative.Target
}

// RootRelativePath returns filePath relative to root with `/` separators, or
// "" when filePath is outside root (or root is empty). Archive members are
// handled naturally: both the container and the `!`-separated member path are
// part of the lexical path under root.
func RootRelativePath(root, filePath string) string {
	if root == "" || filePath == "" {
		return ""
	}
	rel, err := filepath.Rel(root, filepath.FromSlash(filePath))
	if err != nil {
		return ""
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return ""
	}
	return filepath.ToSlash(rel)
}
