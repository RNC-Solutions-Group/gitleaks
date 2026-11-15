package cmd

import "fmt"

const maxAggregateCollisionCandidates = 1000

// validateReportDestinationSet checks all intended output paths as one set.
// Direct aliases and aggregate-numbered candidates are compared. The checks
// are preflight-only and do not claim crash-atomic publication across files.
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
			if mode == reportAggregate && aggregateDestinationsOverlap(paths[i].path, paths[j].path) {
				return fmt.Errorf("%s path %q can collide with %s path %q through aggregate numbering", paths[i].label, paths[i].path, paths[j].label, paths[j].path)
			}
		}
	}
	if mode == reportAggregate && companion != "" && grouped != "" && usableDestination(primary) != "" {
		for number := 1; number <= maxAggregateCollisionCandidates; number++ {
			candidate := informativeCompanionPath(numberedReportPath(primary, number))
			if samePath(candidate, grouped) || aggregateDestinationsOverlap(candidate, grouped) {
				return fmt.Errorf("grouped report path %q can collide with aggregate informative companion %q", grouped, candidate)
			}
		}
	}
	return nil
}

func usableDestination(path string) string {
	if path == "" || path == "-" {
		return ""
	}
	return path
}

func aggregateDestinationsOverlap(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	if samePath(a, b) {
		return true
	}
	for number := 1; number <= maxAggregateCollisionCandidates; number++ {
		if samePath(a, numberedReportPath(b, number)) || samePath(numberedReportPath(a, number), b) {
			return true
		}
	}
	return false
}
