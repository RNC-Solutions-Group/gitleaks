package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zricethezav/gitleaks/v8/report"
)

const maxPriorReportSize = 10 * 1_000 * 1_000

// validateReportDestination checks a report path without creating, truncating,
// or removing it. It rejects destinations that are symlinks, directories, or
// special files, and it rejects existing files inside the scan target unless
// they are recognizable prior reports for the selected format.
//
// The report writer performs the actual open/write/replace only after a
// successful scan. This function deliberately has no side effects on the
// destination.
func validateReportDestination(reportPath, reportFormat, source, configPath, baselinePath, gitleaksIgnorePath string) error {
	if reportPath == "" || reportPath == report.StdoutReportPath {
		return nil
	}

	info, err := os.Lstat(reportPath)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("report path %q is a symbolic link; refusing to replace its target", reportPath)
		}
		if info.IsDir() {
			return fmt.Errorf("report path %q is a directory", reportPath)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("report path %q is not a regular file", reportPath)
		}
	case os.IsNotExist(err):
		info = nil
	default:
		return fmt.Errorf("could not inspect report path %q: %w", reportPath, err)
	}

	parent := filepath.Dir(reportPath)
	parentInfo, err := os.Stat(parent)
	if err != nil {
		return fmt.Errorf("report directory %q is not accessible: %w", parent, err)
	}
	if !parentInfo.IsDir() {
		return fmt.Errorf("report directory %q is not a directory", parent)
	}

	if source != "" {
		sourceInfo, statErr := os.Stat(source)
		if statErr == nil {
			if !sourceInfo.IsDir() {
				if samePath(reportPath, source) {
					return fmt.Errorf("report path %q collides with the scan target", reportPath)
				}
			} else if info != nil && pathWithin(reportPath, source) {
				if reportFormat == "template" {
					return fmt.Errorf("report path %q is inside the scan target and cannot be identified as a prior report", reportPath)
				}
				ok, priorErr := looksLikePriorReport(reportPath, reportFormat)
				if priorErr != nil {
					return fmt.Errorf("could not inspect existing report %q: %w", reportPath, priorErr)
				}
				if !ok {
					return fmt.Errorf("report path %q is inside the scan target and does not look like a prior %s report", reportPath, reportFormat)
				}
			}
		}
	}

	for _, input := range []string{
		configPath,
		baselinePath,
		gitleaksIgnorePath,
		filepath.Join(gitleaksIgnorePath, ".gitleaksignore"),
		filepath.Join(source, ".gitleaksignore"),
	} {
		if input == "" {
			continue
		}
		if samePath(reportPath, input) {
			return fmt.Errorf("report path %q collides with scan input %q", reportPath, input)
		}
	}

	return nil
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
// considered a child; a destination equal to the root is handled separately.
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

// looksLikePriorReport recognizes a prior output for the selected format. It is
// intentionally conservative: a file inside the scan target that cannot be
// recognized is treated as a possible input and rejected.
func looksLikePriorReport(path, format string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxPriorReportSize+1))
	if err != nil {
		return false, err
	}
	if len(data) > maxPriorReportSize {
		return false, fmt.Errorf("existing file is larger than %d bytes", maxPriorReportSize)
	}

	switch format {
	case "json":
		trimmed := bytes.TrimSpace(data)
		if len(trimmed) == 0 || trimmed[0] != '[' {
			return false, nil
		}
		var findings []report.Finding
		if err := json.Unmarshal(trimmed, &findings); err != nil {
			return false, nil
		}
		return true, nil

	case "csv":
		if len(bytes.TrimSpace(data)) == 0 {
			// The CSV reporter intentionally writes an empty file when there
			// are no findings, so a prior empty CSV is recognizable output.
			return true, nil
		}
		reader := csv.NewReader(bytes.NewReader(data))
		header, err := reader.Read()
		if err != nil {
			return false, nil
		}
		fields := make(map[string]bool, len(header))
		for _, field := range header {
			fields[field] = true
		}
		return fields["RuleID"] && fields["File"] && fields["Secret"], nil

	case "junit":
		var suites struct {
			XMLName xml.Name `xml:"testsuites"`
			Suites  []struct {
				Name string `xml:"name,attr"`
			} `xml:"testsuite"`
		}
		if err := xml.Unmarshal(data, &suites); err != nil {
			return false, nil
		}
		return len(suites.Suites) == 1 && suites.Suites[0].Name == "gitleaks", nil

	case "sarif":
		var sarif report.Sarif
		if err := json.Unmarshal(data, &sarif); err != nil {
			return false, nil
		}
		return sarif.Version == "2.1.0" && sarif.Runs != nil, nil

	default:
		return false, nil
	}
}

// writeReportAtomic writes the report to a uniquely created temporary file in
// the destination directory, syncs and closes it, and then replaces the
// destination. Any failure before the rename preserves the prior destination.
func writeReportAtomic(reportPath string, reporter report.Reporter, findings []report.Finding) error {
	dir := filepath.Dir(reportPath)
	base := filepath.Base(reportPath)
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
	if err := os.Rename(tempName, reportPath); err != nil {
		return fmt.Errorf("could not replace report destination: %w", err)
	}
	renamed = true
	return nil
}
