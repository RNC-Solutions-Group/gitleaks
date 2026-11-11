package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/report"
)

// reportMode is the resolved file-output behavior for a run.
type reportMode int

const (
	// reportOverwrite replaces the requested destination after a successful
	// scan. This preserves the historical Gitleaks behavior.
	reportOverwrite reportMode = iota
	// reportAggregate writes a numbered report and never replaces an existing
	// report file.
	reportAggregate
	// reportReject refuses to run when the requested destination exists.
	reportReject
)

func (m reportMode) String() string {
	switch m {
	case reportOverwrite:
		return "overwrite"
	case reportAggregate:
		return "aggregate"
	case reportReject:
		return "reject"
	default:
		return "unknown"
	}
}

// resolveReportMode applies CLI precedence over the configuration policy.
// Aggregation takes precedence when both configuration settings are enabled.
func resolveReportMode(cfg config.Config, cmd *cobra.Command) (reportMode, error) {
	overwriteSet := cmd.Flags().Changed("overwrite")
	aggregateSet := cmd.Flags().Changed("aggregate")
	if overwriteSet && aggregateSet {
		return reportReject, errors.New("--overwrite and --aggregate cannot be used together")
	}

	if aggregateSet {
		if aggregate, err := cmd.Flags().GetBool("aggregate"); err != nil {
			return reportReject, err
		} else if aggregate {
			return reportAggregate, nil
		}
	}
	if overwriteSet {
		if overwrite, err := cmd.Flags().GetBool("overwrite"); err != nil {
			return reportReject, err
		} else if overwrite {
			return reportOverwrite, nil
		}
	}

	if cfg.Aggregation {
		return reportAggregate, nil
	}
	if cfg.AutoOverwrite {
		return reportOverwrite, nil
	}
	return reportReject, nil
}

// reportOutput owns the output policy for one run. The requested destination is
// validated without modification before scanning. The actual report is written
// to a temporary file and published only after the scan and writer succeed.
type reportOutput struct {
	path            string
	mode            reportMode
	source          string
	sourceIsFile    bool
	protectedInputs []string
	rename          func(string, string) error
}

// newReportOutput validates the requested destination before scanning. It does
// not create, truncate, rename, or remove the destination.
func newReportOutput(
	reportPath string,
	mode reportMode,
	source string,
	configPath string,
	baselinePath string,
	gitleaksIgnorePath string,
	reportTemplatePath string,
) (*reportOutput, error) {
	out := &reportOutput{
		path:   reportPath,
		mode:   mode,
		source: source,
		rename: os.Rename,
	}
	out.protectedInputs = reportInputPaths(configPath, baselinePath, gitleaksIgnorePath, reportTemplatePath, source)

	if sourceInfo, err := os.Stat(source); err == nil && !sourceInfo.IsDir() {
		out.sourceIsFile = true
	}

	if _, err := out.validateRequestedPath(); err != nil {
		return nil, err
	}
	return out, nil
}

func reportInputPaths(configPath, baselinePath, gitleaksIgnorePath, reportTemplatePath, source string) []string {
	inputs := []string{configPath, baselinePath, reportTemplatePath}
	if gitleaksIgnorePath != "" {
		inputs = append(inputs, gitleaksIgnorePath, filepath.Join(gitleaksIgnorePath, ".gitleaksignore"))
	}
	if source != "" {
		inputs = append(inputs, filepath.Join(source, ".gitleaksignore"))
	}
	return inputs
}

// validateRequestedPath checks the requested destination and policy. A
// destination that already exists is rejected only in reject mode; aggregation
// will move to a numbered filename instead.
func (o *reportOutput) validateRequestedPath() (bool, error) {
	exists, err := inspectReportPath(o.path)
	if err != nil {
		return false, err
	}
	if err := validateReportParent(o.path); err != nil {
		return false, err
	}
	if err := o.checkInputCollision(o.path); err != nil {
		return false, err
	}
	if o.mode == reportReject && exists {
		return true, fmt.Errorf(
			"report path %q already exists and both auto-overwrite and aggregation are disabled; "+
				"enable auto-overwrite or aggregation in the configuration, or pass --overwrite or --aggregate",
			o.path,
		)
	}
	return exists, nil
}

// validateForPublication repeats the non-destructive checks immediately before
// report generation and again before publication.
func (o *reportOutput) validateForPublication() error {
	_, err := o.validateRequestedPath()
	return err
}

func inspectReportPath(path string) (bool, error) {
	info, err := os.Lstat(path)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink != 0 {
			return true, fmt.Errorf("report path %q is a symbolic link; refusing to replace its target", path)
		}
		if info.IsDir() {
			return true, fmt.Errorf("report path %q is a directory", path)
		}
		if !info.Mode().IsRegular() {
			return true, fmt.Errorf("report path %q is not a regular file", path)
		}
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, fmt.Errorf("could not inspect report path %q: %w", path, err)
	}
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

func (o *reportOutput) checkInputCollision(path string) error {
	if o.sourceIsFile && samePath(path, o.source) {
		return fmt.Errorf("report path %q collides with the scan target", path)
	}
	for _, input := range o.protectedInputs {
		if input == "" {
			continue
		}
		if samePath(path, input) {
			return fmt.Errorf("report path %q collides with scan input %q", path, input)
		}
	}
	return nil
}

// write produces the report after a completed scan and returns the actual
// output filename. A failed writer or publication leaves the previous report
// destination untouched.
func (o *reportOutput) write(reporter report.Reporter, findings []report.Finding) (string, error) {
	if err := o.validateForPublication(); err != nil {
		return "", err
	}

	temp, err := createReportTemp(o.path)
	if err != nil {
		return "", err
	}
	tempName := temp.Name()
	published := false
	defer func() {
		if !published {
			_ = temp.Close()
			_ = os.Remove(tempName)
		}
	}()

	if runtime.GOOS != "windows" {
		if err := temp.Chmod(0o600); err != nil {
			return "", fmt.Errorf("could not set report permissions: %w", err)
		}
	}
	if err := reporter.Write(temp, findings); err != nil {
		return "", fmt.Errorf("could not write report: %w", err)
	}
	if err := temp.Sync(); err != nil {
		return "", fmt.Errorf("could not flush report: %w", err)
	}
	if err := temp.Close(); err != nil {
		return "", fmt.Errorf("could not close report: %w", err)
	}

	// Revalidate immediately before publication. In particular, a destination
	// that became a symlink, directory, or special file is never followed.
	if err := o.validateForPublication(); err != nil {
		return "", err
	}

	var outputPath string
	switch o.mode {
	case reportOverwrite:
		outputPath = o.path
		if err := o.rename(tempName, outputPath); err != nil {
			return "", fmt.Errorf("could not replace report destination: %w", err)
		}
	case reportAggregate:
		outputPath, err = o.publishAggregate(tempName)
		if err != nil {
			return "", err
		}
	case reportReject:
		outputPath, err = o.publishReject(tempName)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported report mode %s", o.mode)
	}

	published = true
	return outputPath, nil
}

func createReportTemp(reportPath string) (*os.File, error) {
	dir := filepath.Dir(reportPath)
	base := filepath.Base(reportPath)
	temp, err := os.CreateTemp(dir, "."+base+".tmp-*")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary report: %w", err)
	}
	return temp, nil
}

// publishAggregate reserves the first available filename and publishes the
// completed temporary report there. The exclusive create is the allocation
// decision, so concurrent processes cannot claim the same numbered filename.
func (o *reportOutput) publishAggregate(tempName string) (string, error) {
	if path, info, err := reserveReportPath(o.path); err == nil {
		return path, o.publishReservedReport(tempName, path, info)
	} else if !errors.Is(err, fs.ErrExist) {
		return "", err
	}

	for number := 1; ; number++ {
		candidate := numberedReportPath(o.path, number)
		path, info, err := reserveReportPath(candidate)
		if err == nil {
			return path, o.publishReservedReport(tempName, path, info)
		}
		if !errors.Is(err, fs.ErrExist) {
			return "", err
		}
	}
}

// publishReject publishes only when the destination is still absent. The
// exclusive reservation closes the check-then-rename race.
func (o *reportOutput) publishReject(tempName string) (string, error) {
	path, info, err := reserveReportPath(o.path)
	if errors.Is(err, fs.ErrExist) {
		return "", fmt.Errorf("report path %q appeared before publication; refusing to overwrite it", o.path)
	}
	if err != nil {
		return "", err
	}
	return path, o.publishReservedReport(tempName, path, info)
}

func reserveReportPath(path string) (string, os.FileInfo, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", nil, fmt.Errorf("could not reserve report path %q: %w", path, err)
	}
	info, statErr := file.Stat()
	closeErr := file.Close()
	if statErr != nil {
		_ = os.Remove(path)
		return "", nil, fmt.Errorf("could not inspect reserved report path %q: %w", path, statErr)
	}
	if closeErr != nil {
		_ = os.Remove(path)
		return "", nil, fmt.Errorf("could not close reserved report path %q: %w", path, closeErr)
	}
	return path, info, nil
}

func (o *reportOutput) publishReservedReport(tempName, path string, reserved os.FileInfo) error {
	if err := o.rename(tempName, path); err != nil {
		// Do not leave an empty reservation behind. Never remove a file that
		// has been replaced by another process.
		if current, statErr := os.Lstat(path); statErr == nil && os.SameFile(current, reserved) {
			_ = os.Remove(path)
		}
		return fmt.Errorf("could not publish report %q: %w", path, err)
	}
	return nil
}

func numberedReportPath(path string, number int) string {
	dir, base := filepath.Split(path)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if stem == "" {
		// Treat dotfiles as a stem rather than producing a leading "_001".
		stem = base
		ext = ""
	}
	return filepath.Join(dir, fmt.Sprintf("%s_%03d%s", stem, number, ext))
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
