package detect

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func IsNew(finding report.Finding, redact uint, baseline []report.Finding) bool {
	// Explicitly testing each property as it gives significantly better performance in comparison to cmp.Equal(). Drawback is that
	// the code requires maintenance if/when the Finding struct changes
	for _, b := range baseline {
		if finding.RuleID == b.RuleID &&
			finding.Description == b.Description &&
			finding.StartLine == b.StartLine &&
			finding.EndLine == b.EndLine &&
			finding.StartColumn == b.StartColumn &&
			finding.EndColumn == b.EndColumn &&
			(redact > 0 || (finding.Match == b.Match && finding.Secret == b.Secret)) &&
			finding.File == b.File &&
			finding.Commit == b.Commit &&
			finding.Author == b.Author &&
			finding.Email == b.Email &&
			finding.Date == b.Date &&
			finding.Message == b.Message &&
			// Omit checking finding.Fingerprint - if the format of the fingerprint changes, the users will see unexpected behaviour
			finding.Entropy == b.Entropy {
			return false
		}
	}
	return true
}

func LoadBaseline(baselinePath string) ([]report.Finding, error) {
	bytes, err := os.ReadFile(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("could not open %s", baselinePath)
	}

	var previousFindings []report.Finding
	err = json.Unmarshal(bytes, &previousFindings)
	if err != nil {
		return nil, fmt.Errorf("the format of the file %s is not supported", baselinePath)
	}

	return previousFindings, nil
}

func (d *Detector) AddBaseline(baselinePath string, source string) error {
	if baselinePath != "" {
		absoluteSource, err := filepath.Abs(source)
		if err != nil {
			return err
		}

		absoluteBaseline, err := filepath.Abs(baselinePath)
		if err != nil {
			return err
		}

		relativeBaseline, err := filepath.Rel(absoluteSource, absoluteBaseline)
		if err != nil {
			return err
		}

		baseline, err := LoadBaseline(baselinePath)
		if err != nil {
			return err
		}

		// A directory scan reports File the way its root was supplied, and
		// IsNew compares File. Bring entries that are demonstrably under the
		// source to their root-relative spelling in memory, without modifying
		// the baseline file. Entries that cannot be placed under the source
		// stay exactly as written. This is the baseline part of the design
		// proposed in gitleaks/gitleaks#2248 (@hamodywe).
		for i := range baseline {
			if rel, ok := baselineFileRelative(absoluteSource, baseline[i].File); ok {
				baseline[i].File = rel
			}
		}

		d.baseline = baseline
		baselinePath = relativeBaseline

	}

	d.baselinePath = baselinePath
	return nil
}

// baselineFileRelative maps a baseline File value onto the current source's
// root-relative spelling when that can be established. Absolute entries under
// the source and ordinary root-relative entries are unambiguous. A relative
// root spelling (for example `repo/src/x.py` produced by `gitleaks dir repo`)
// is also accepted when the leading source basename can be removed and the
// resulting path exists under the source. The existence check prevents
// stripping a legitimate nested path that merely begins with the same name.
func baselineFileRelative(source, file string) (string, bool) {
	if file == "" {
		return "", false
	}
	if filepath.IsAbs(file) {
		rel := sources.RootRelativePath(source, file)
		return rel, rel != ""
	}

	clean := filepath.Clean(filepath.FromSlash(file))
	if filepath.IsLocal(clean) {
		base := filepath.Base(source)
		if clean != base && strings.HasPrefix(clean, base+string(filepath.Separator)) {
			// A relative-root baseline spells paths as `repo/src/x.py`.
			// Do not strip that prefix if it is also a real nested path
			// under the current root.
			if baselineSourcePathExists(source, clean) {
				return filepath.ToSlash(clean), true
			}
			stripped := strings.TrimPrefix(clean, base+string(filepath.Separator))
			if baselineSourcePathExists(source, stripped) {
				return filepath.ToSlash(stripped), true
			}
		}
		return filepath.ToSlash(clean), true
	}

	// The path may have been recorded relative to a different working
	// directory. Resolve it against the current one as a final fallback.
	abs, err := filepath.Abs(file)
	if err != nil {
		return "", false
	}
	rel := sources.RootRelativePath(source, abs)
	return rel, rel != ""
}

// baselineSourcePathExists reports whether the outer file of a baseline path
// exists beneath source. Archive members keep their `!`-separated inner path;
// only the container can be stat'ed.
func baselineSourcePathExists(source, file string) bool {
	candidate := filepath.Join(source, filepath.FromSlash(file))
	if sources.RootRelativePath(source, candidate) != "" {
		if _, err := os.Stat(candidate); err == nil {
			return true
		}
	}

	outer, _, nested := strings.Cut(filepath.ToSlash(file), sources.InnerPathSeparator)
	if !nested {
		return false
	}
	candidate = filepath.Join(source, filepath.FromSlash(outer))
	if sources.RootRelativePath(source, candidate) == "" {
		return false
	}
	_, err := os.Stat(candidate)
	return err == nil
}
