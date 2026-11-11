package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zricethezav/gitleaks/v8/report"
)

// reportDestination captures the destination policy established before the
// scan. Revalidate is called again immediately before replacement so a file
// that appeared in the scan target during the scan is not silently replaced.
type reportDestination struct {
	path                 string
	source               string
	inTree               bool
	existedAtValidation  bool
	allowInTreeOverwrite bool
	sourceIsFile         bool
	soleSourceFile       bool
	protectedInputs      []string
}

// validateReportDestination checks a report path without creating, truncating,
// or removing it. It deliberately does not inspect report-like content: a file
// that resembles a report is not proof that Gitleaks owns it.
func validateReportDestination(
	reportPath string,
	source string,
	configPath string,
	baselinePath string,
	gitleaksIgnorePath string,
	allowInTreeOverwrite bool,
) (reportDestination, error) {
	dest := reportDestination{
		path:                 reportPath,
		source:               source,
		allowInTreeOverwrite: allowInTreeOverwrite,
	}
	if reportPath == "" || reportPath == report.StdoutReportPath {
		return dest, nil
	}

	info, err := os.Lstat(reportPath)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink != 0 {
			return dest, fmt.Errorf("report path %q is a symbolic link; refusing to replace its target", reportPath)
		}
		if info.IsDir() {
			return dest, fmt.Errorf("report path %q is a directory", reportPath)
		}
		if !info.Mode().IsRegular() {
			return dest, fmt.Errorf("report path %q is not a regular file", reportPath)
		}
		dest.existedAtValidation = true
	case os.IsNotExist(err):
		info = nil
	default:
		return dest, fmt.Errorf("could not inspect report path %q: %w", reportPath, err)
	}

	if err := validateReportParent(reportPath); err != nil {
		return dest, err
	}

	if source != "" {
		sourceInfo, statErr := os.Stat(source)
		if statErr == nil {
			if !sourceInfo.IsDir() {
				dest.sourceIsFile = true
				if samePath(reportPath, source) {
					return dest, fmt.Errorf("report path %q collides with the scan target", reportPath)
				}
			} else {
				dest.inTree = pathWithin(reportPath, source)
			}
		}
	}

	dest.protectedInputs = []string{
		configPath,
		baselinePath,
		gitleaksIgnorePath,
		filepath.Join(gitleaksIgnorePath, ".gitleaksignore"),
		filepath.Join(source, ".gitleaksignore"),
	}
	for _, input := range dest.protectedInputs {
		if input == "" {
			continue
		}
		if samePath(reportPath, input) {
			return dest, fmt.Errorf("report path %q collides with scan input %q", reportPath, input)
		}
	}

	if dest.inTree && info != nil {
		if !allowInTreeOverwrite {
			return dest, fmt.Errorf(
				"report path %q already exists inside the scan target; use --allow-in-tree-report-overwrite to replace it",
				reportPath,
			)
		}
		sole, err := soleRegularFile(source, reportPath)
		if err != nil {
			return dest, fmt.Errorf("could not inspect scan target %q: %w", source, err)
		}
		if sole {
			dest.soleSourceFile = true
			return dest, fmt.Errorf("report path %q is the sole file in the scan target and cannot be replaced", reportPath)
		}
	}

	return dest, nil
}

// revalidate re-checks the destination immediately before it is replaced.
func (d reportDestination) revalidate() error {
	if d.path == "" || d.path == report.StdoutReportPath {
		return nil
	}

	info, err := os.Lstat(d.path)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("report path %q became a symbolic link", d.path)
		}
		if info.IsDir() {
			return fmt.Errorf("report path %q is a directory", d.path)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("report path %q is not a regular file", d.path)
		}
		if d.inTree && !d.allowInTreeOverwrite {
			return fmt.Errorf("report path %q exists inside the scan target and overwrite was not authorized", d.path)
		}
		if d.inTree && !d.existedAtValidation {
			return fmt.Errorf("report path %q appeared inside the scan target after validation", d.path)
		}
		if d.inTree && d.soleSourceFile {
			return fmt.Errorf("report path %q is the sole file in the scan target", d.path)
		}
	case os.IsNotExist(err):
		// A destination that was new inside the scan target is still new. A
		// previously existing destination may have been removed; writing a new
		// file at the authorized path is safe.
	default:
		return fmt.Errorf("could not inspect report path %q: %w", d.path, err)
	}

	if err := validateReportParent(d.path); err != nil {
		return err
	}

	for _, input := range d.protectedInputs {
		if input == "" {
			continue
		}
		if samePath(d.path, input) {
			return fmt.Errorf("report path %q collides with scan input %q", d.path, input)
		}
	}

	if d.sourceIsFile && samePath(d.path, d.source) {
		return fmt.Errorf("report path %q collides with the scan target", d.path)
	}

	return nil
}

func validateReportParent(reportPath string) error {
	parent := filepath.Dir(reportPath)
	info, err := os.Stat(parent)
	if err != nil {
		return fmt.Errorf("report directory %q is not accessible: %w", parent, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("report directory %q is not a directory", parent)
	}
	return nil
}

// soleRegularFile reports whether reportPath is the only regular file beneath
// source. It is used to avoid replacing the only file being scanned.
func soleRegularFile(source, reportPath string) (bool, error) {
	var (
		count int
		first string
	)
	stop := errors.New("more than one file")
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		count++
		if count == 1 {
			first = path
		}
		if count > 1 {
			return stop
		}
		return nil
	})
	if errors.Is(err, stop) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count == 1 && samePath(first, reportPath), nil
}

// samePath compares two paths lexically and, when both exist, by file identity.
func samePath(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}

	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA == nil && errB == nil && filepath.Clean(absA) == filepath.Clean(absB) {
		return true
	}

	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(infoA, infoB)
}

// pathWithin reports whether path is a child of root. The root itself is not
// considered a child.
func pathWithin(path, root string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}

	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		absPath = resolved
	}
	if resolved, err := filepath.EvalSymlinks(absRoot); err == nil {
		absRoot = resolved
	}

	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return false
	}
	if rel == "." || rel == ".." || filepath.IsAbs(rel) {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// writeReportAtomic writes the report to a uniquely created temporary file in
// the destination directory, syncs and closes it, revalidates the destination,
// and then replaces it. Any failure before the rename preserves the prior
// destination.
func writeReportAtomic(destination reportDestination, reporter report.Reporter, findings []report.Finding) error {
	if err := destination.revalidate(); err != nil {
		return err
	}

	dir := filepath.Dir(destination.path)
	base := filepath.Base(destination.path)
	temp, err := os.CreateTemp(dir, "."+base+".tmp-*")
	if err != nil {
		return fmt.Errorf("could not create temporary report: %w", err)
	}

	tempName := temp.Name()
	renamed := false
	defer func() {
		if !renamed {
			_ = temp.Close()
			_ = os.Remove(tempName)
		}
	}()

	if runtime.GOOS != "windows" {
		if err := temp.Chmod(0o600); err != nil {
			return fmt.Errorf("could not set report permissions: %w", err)
		}
	}

	if err := reporter.Write(temp, findings); err != nil {
		return fmt.Errorf("could not write report: %w", err)
	}
	if err := temp.Sync(); err != nil {
		return fmt.Errorf("could not flush report: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("could not close report: %w", err)
	}

	// Revalidate immediately before replacement. In particular, a new in-tree
	// destination that appeared during the scan is not overwritten.
	if err := destination.revalidate(); err != nil {
		return err
	}

	if err := os.Rename(tempName, destination.path); err != nil {
		return fmt.Errorf("could not replace report destination: %w", err)
	}
	renamed = true
	return nil
}
