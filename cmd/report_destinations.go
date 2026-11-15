package cmd

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// validateReportDestinationSet checks all intended output paths as one set.
// Exact aliases are compared directly. In aggregate mode, predictable naming
// families are compared structurally, so validation does not enumerate every
// possible numbered filename. The checks are preflight-only and do not claim
// crash-atomic publication across files.
func validateReportDestinationSet(primary, companion, grouped string, mode reportMode) error {
	paths := []struct {
		label string
		path  string
	}{
		{"primary report", usableDestination(primary)},
		{"informative companion", companion},
		{"grouped report", grouped},
	}
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if paths[i].path == "" || paths[j].path == "" {
				continue
			}
			if samePath(paths[i].path, paths[j].path) {
				return fmt.Errorf("%s path %q collides with %s path %q", paths[i].label, paths[i].path, paths[j].label, paths[j].path)
			}
		}
	}

	if mode != reportAggregate {
		return nil
	}

	primary = usableDestination(primary)
	if primary == "" {
		return nil
	}

	if grouped != "" {
		if aggregateFamiliesOverlap(primary, grouped) {
			return fmt.Errorf("primary report path %q can collide with grouped report path %q through aggregate numbering", primary, grouped)
		}
		if companion != "" && companionFamilyContains(primary, grouped) {
			return fmt.Errorf("grouped report path %q can collide with an aggregate informative companion for primary report %q", grouped, primary)
		}
	}
	if companion != "" && grouped != "" && aggregateFamilyContains(grouped, companion) {
		return fmt.Errorf("grouped report path %q can collide with informative companion path %q through aggregate numbering", grouped, companion)
	}
	return nil
}

func usableDestination(path string) string {
	if path == "" || path == "-" {
		return ""
	}
	return path
}

// aggregateFamiliesOverlap reports whether the base name or any numbered name
// generated for one destination can equal the other destination's base name or
// generated name. With aggregate numbering this is a static naming property:
// report.json overlaps report_001.json, but not report_001_001.json.
func aggregateFamiliesOverlap(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	if !sameReportDirectory(a, b) {
		return false
	}
	return aggregateFamilyContains(a, b) || aggregateFamilyContains(b, a)
}

// aggregateFamilyContains reports whether member is the aggregate family base
// or a numbered name generated from family.
func aggregateFamilyContains(family, member string) bool {
	if family == "" || member == "" {
		return false
	}
	if samePath(family, member) {
		return true
	}
	if !sameReportDirectory(family, member) {
		return false
	}
	return familyNameMember(filepath.Base(family), filepath.Base(member))
}

// companionFamilyContains reports whether member is the informative companion
// for the aggregate primary base or for one of its numbered candidates. The
// companion name is derived from the selected primary name, not from an
// independent numbering family.
func companionFamilyContains(primary, member string) bool {
	if primary == "" || member == "" {
		return false
	}
	if samePath(informativeCompanionPath(primary), member) {
		return true
	}
	if !sameReportDirectory(primary, member) {
		return false
	}

	primaryBase := filepath.Base(primary)
	primaryStem := strings.TrimSuffix(primaryBase, filepath.Ext(primaryBase))
	if primaryStem == "" {
		primaryStem = primaryBase
	}
	prefix := primaryStem + "_"
	const suffix = "-informative.json"
	memberBase := filepath.Base(member)
	if !hasReportNameSuffix(memberBase, suffix) || !hasReportNamePrefix(memberBase, prefix) {
		return false
	}
	digits := memberBase[len(prefix) : len(memberBase)-len(suffix)]
	return isCanonicalReportNumber(digits)
}

func sameReportDirectory(a, b string) bool {
	return samePath(filepath.Dir(a), filepath.Dir(b))
}

func familyNameMember(familyBase, memberBase string) bool {
	ext := filepath.Ext(familyBase)
	stem := strings.TrimSuffix(familyBase, ext)
	if stem == "" {
		stem = familyBase
		ext = ""
	}
	prefix := stem + "_"
	if !hasReportNameSuffix(memberBase, ext) || !hasReportNamePrefix(memberBase, prefix) {
		return false
	}
	digits := memberBase[len(prefix) : len(memberBase)-len(ext)]
	return isCanonicalReportNumber(digits)
}

// isCanonicalReportNumber accepts exactly the decimal spelling produced by
// numberedReportPath (%03d for the first 999 candidates). This keeps e.g.
// report_0001.json outside report.json's generated family.
func isCanonicalReportNumber(digits string) bool {
	if digits == "" {
		return false
	}
	number, err := strconv.Atoi(digits)
	if err != nil || number < 1 {
		return false
	}
	return fmt.Sprintf("%03d", number) == digits
}

func hasReportNamePrefix(name, prefix string) bool {
	if len(name) < len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.EqualFold(name[:len(prefix)], prefix)
	}
	return name[:len(prefix)] == prefix
}

func hasReportNameSuffix(name, suffix string) bool {
	if len(name) < len(suffix) {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.EqualFold(name[len(name)-len(suffix):], suffix)
	}
	return name[len(name)-len(suffix):] == suffix
}
