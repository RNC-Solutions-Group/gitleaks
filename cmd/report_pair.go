package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/report"
)

// enableInformative configures the optional informative JSON companion. The
// destination is validated before scanning; publication is paired with the
// primary report but is not claimed to be atomic across two filesystem paths.
func (o *reportOutput) enableInformative(render func(primaryPath string) ([]byte, error)) error {
	if o == nil {
		return errors.New("report output is not configured")
	}
	if o.path == report.StdoutReportPath {
		return errors.New("stdout reporting cannot have an informative companion file")
	}
	if render == nil {
		return errors.New("informative report renderer is not configured")
	}
	if o.companion != nil {
		return errors.New("informative report companion is already configured")
	}

	companionPath := informativeCompanionPath(o.path)
	if samePath(companionPath, o.path) {
		return fmt.Errorf("informative report path %q collides with the primary report path", companionPath)
	}
	if err := o.validateAdditionalPath(companionPath, o.mode == reportReject); err != nil {
		return err
	}
	o.companion = &companionReport{
		path:   companionPath,
		render: render,
	}
	return nil
}

// validateAdditionalPath applies destination safety and collision rules to a
// second report file. rejectExists implements the reject-existing policy; the
// primary report's existing policy remains independent for overwrite and
// aggregate modes.
func (o *reportOutput) validateAdditionalPath(path string, rejectExists bool) error {
	exists, err := inspectReportPath(path)
	if err != nil {
		return err
	}
	if err := validateReportParent(path); err != nil {
		return err
	}
	if err := o.checkInputCollision(path); err != nil {
		return err
	}
	if rejectExists && exists {
		return fmt.Errorf(
			"informative report path %q already exists and report output policy is reject; "+
				"enable auto-overwrite or aggregation, or pass --overwrite or --aggregate",
			path,
		)
	}
	return nil
}

// informativeCompanionPath derives the companion filename from the actual
// primary filename. The informative document is always JSON, regardless of the
// primary report format.
func informativeCompanionPath(primaryPath string) string {
	dir, base := filepath.Split(primaryPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if stem == "" {
		// Preserve dotfiles as a stem rather than producing a leading hyphen.
		stem = base
	}
	return filepath.Join(dir, stem+"-informative.json")
}

// writePaired renders both outputs before publication, then publishes the
// informative companion first and the primary report second. Render or
// validation failures leave both existing destinations untouched. A late
// primary publication failure can leave the new companion in place for
// overwrite mode because two independent paths cannot be replaced atomically.
func (o *reportOutput) writePaired(reporter report.Reporter, findings []report.Finding) (string, error) {
	if o.companion == nil {
		return "", errors.New("informative report companion is not configured")
	}

	// Repeat the non-destructive checks immediately before rendering. For
	// reject mode, no reservation exists yet, so this also detects a destination
	// that appeared during the scan.
	if err := o.validateForPublication(); err != nil {
		return "", err
	}
	if err := o.validateCompanionRequestedPath(); err != nil {
		return "", err
	}

	primaryTemp, err := writeReportTemp(o.path, "report", func(out io.WriteCloser) error {
		return reporter.Write(out, findings)
	})
	if err != nil {
		return "", err
	}
	published := false
	defer func() {
		if !published {
			_ = os.Remove(primaryTemp)
		}
	}()

	allocation, err := o.allocatePair()
	if err != nil {
		return "", err
	}
	defer allocation.release()

	informativeBytes, err := o.companion.render(allocation.primary)
	if err != nil {
		return "", fmt.Errorf("could not render informative report: %w", err)
	}
	informativeTemp, err := writeReportTemp(o.companion.path, "informative report", func(out io.WriteCloser) error {
		_, writeErr := out.Write(informativeBytes)
		return writeErr
	})
	if err != nil {
		return "", err
	}
	defer func() {
		if !published {
			_ = os.Remove(informativeTemp)
		}
	}()

	// Revalidate the actual destinations, including an aggregate candidate that
	// differs from the requested filename, immediately before publication.
	if err := o.validatePairSafety(allocation); err != nil {
		return "", err
	}

	// In overwrite mode an existing companion is preserved before it is
	// replaced. If the primary rename then fails, the old companion can be
	// restored so it is not presented as belonging to the still-old primary.
	// This is not crash-atomic; it closes the handled-failure window.
	companionBackup, hadCompanion, err := backupCompanion(allocation.companion, o.mode == reportOverwrite)
	if err != nil {
		return "", err
	}
	keepBackup := false
	defer func() {
		if companionBackup != "" && !keepBackup {
			_ = os.Remove(companionBackup)
		}
	}()

	// Publish the companion first so a primary publication failure leaves the
	// primary output untouched. The reverse order would risk losing the primary
	// report on a late companion failure.
	if err := o.rename(informativeTemp, allocation.companion); err != nil {
		return "", fmt.Errorf("could not publish informative report %q: %w", allocation.companion, err)
	}
	companionPublished, _ := os.Stat(allocation.companion)
	if err := o.rename(primaryTemp, allocation.primary); err != nil {
		// Capture the identity of the newly published companion after the
		// rename. Windows may replace the reservation's identity, so a current
		// stat is used rather than the pre-publication reservation.
		if hadCompanion {
			current, statErr := os.Lstat(allocation.companion)
			switch {
			case errors.Is(statErr, fs.ErrNotExist):
				if restoreErr := o.rename(companionBackup, allocation.companion); restoreErr != nil {
					keepBackup = true
					return "", fmt.Errorf(
						"could not publish report %q after publishing informative report %q: %w (companion restore failed: %v; preserved backup at %q)",
						allocation.primary,
						allocation.companion,
						err,
						restoreErr,
						companionBackup,
					)
				}
			case statErr != nil:
				keepBackup = true
				return "", fmt.Errorf(
					"could not publish report %q after publishing informative report %q: %w (could not inspect companion for rollback: %v; preserved backup at %q)",
					allocation.primary,
					allocation.companion,
					err,
					statErr,
					companionBackup,
				)
			case sameReportSnapshot(current, companionPublished):
				if removeErr := os.Remove(allocation.companion); removeErr != nil {
					keepBackup = true
					return "", fmt.Errorf(
						"could not publish report %q after publishing informative report %q: %w (could not remove the newly published companion: %v; preserved backup at %q)",
						allocation.primary,
						allocation.companion,
						err,
						removeErr,
						companionBackup,
					)
				}
				if restoreErr := o.rename(companionBackup, allocation.companion); restoreErr != nil {
					keepBackup = true
					return "", fmt.Errorf(
						"could not publish report %q after publishing informative report %q: %w (companion restore failed: %v; preserved backup at %q)",
						allocation.primary,
						allocation.companion,
						err,
						restoreErr,
						companionBackup,
					)
				}
			default:
				// A concurrent writer replaced our companion after publication.
				// Never overwrite or delete that output; retain the prior
				// companion for operator recovery instead of the deferred cleanup.
				keepBackup = true
				return "", fmt.Errorf(
					"could not publish report %q after publishing informative report %q: %w (companion was replaced concurrently; concurrent output preserved at %q; prior companion preserved at %q)",
					allocation.primary,
					allocation.companion,
					err,
					allocation.companion,
					companionBackup,
				)
			}
		} else {
			removeIfSnapshot(allocation.companion, companionPublished)
		}
		return "", fmt.Errorf(
			"could not publish report %q after publishing informative report %q: %w",
			allocation.primary,
			allocation.companion,
			err,
		)
	}

	published = true
	logging.Info().Str("path", allocation.companion).Msg("informative report written")
	return allocation.primary, nil
}

func (o *reportOutput) validateCompanionRequestedPath() error {
	if o.companion == nil {
		return errors.New("informative report companion is not configured")
	}
	return o.validateAdditionalPath(o.companion.path, o.mode == reportReject)
}

func (o *reportOutput) validatePairSafety(allocation *pairAllocation) error {
	if allocation == nil {
		return errors.New("report pair allocation is missing")
	}
	for _, path := range []string{allocation.primary, allocation.companion} {
		if _, err := inspectReportPath(path); err != nil {
			return err
		}
	}
	return o.validatePairStatic(allocation.primary, allocation.companion)
}

// validatePairStatic checks conditions that do not depend on whether a
// destination currently exists. Aggregate allocation uses it so an occupied
// numbered candidate is skipped rather than treated as an invalid destination.
func (o *reportOutput) validatePairStatic(primary, companion string) error {
	for _, path := range []string{primary, companion} {
		if err := validateReportParent(path); err != nil {
			return err
		}
		if err := o.checkInputCollision(path); err != nil {
			return err
		}
	}
	return nil
}

type pairAllocation struct {
	primary   string
	companion string

	primaryReserved   os.FileInfo
	companionReserved os.FileInfo
}

// release removes only empty reservations created by this allocation. Published
// files have a different file identity and are never removed by cleanup.
func (a *pairAllocation) release() {
	if a == nil {
		return
	}
	releaseReservedPath(a.primary, a.primaryReserved)
	releaseReservedPath(a.companion, a.companionReserved)
}

func releaseReservedPath(path string, reserved os.FileInfo) {
	if reserved == nil {
		return
	}
	current, err := os.Lstat(path)
	if err == nil && os.SameFile(current, reserved) {
		_ = os.Remove(path)
	}
}

func removeIfSame(path string, expected os.FileInfo) {
	if expected == nil {
		return
	}
	current, err := os.Lstat(path)
	if err == nil && os.SameFile(current, expected) {
		_ = os.Remove(path)
	}
}

func removeIfSnapshot(path string, expected os.FileInfo) {
	if expected == nil {
		return
	}
	current, err := os.Lstat(path)
	if err == nil && sameReportSnapshot(current, expected) {
		_ = os.Remove(path)
	}
}

func sameReportSnapshot(current, expected os.FileInfo) bool {
	if current == nil || expected == nil || !os.SameFile(current, expected) {
		return false
	}
	return current.Size() == expected.Size() && current.ModTime().Equal(expected.ModTime())
}

// backupCompanion preserves an existing overwrite-mode companion before
// publication. A same-directory hardlink is preferred; a copy is the portable
// fallback for filesystems without hardlink support.
func backupCompanion(path string, enabled bool) (string, bool, error) {
	if !enabled {
		return "", false, nil
	}
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("could not inspect existing informative report %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", false, fmt.Errorf("existing informative report %q is not a regular file", path)
	}

	backup, err := os.CreateTemp(filepath.Dir(path), ".gitleaks-companion-backup-*")
	if err != nil {
		return "", false, fmt.Errorf("could not reserve informative report backup: %w", err)
	}
	backupPath := backup.Name()
	if closeErr := backup.Close(); closeErr != nil {
		_ = os.Remove(backupPath)
		return "", false, fmt.Errorf("could not close informative report backup: %w", closeErr)
	}
	_ = os.Remove(backupPath)

	if linkErr := os.Link(path, backupPath); linkErr != nil {
		if copyErr := copyFile(path, backupPath, info.Mode().Perm()); copyErr != nil {
			_ = os.Remove(backupPath)
			return "", false, fmt.Errorf("could not back up informative report %q: %w", path, copyErr)
		}
	}
	return backupPath, true, nil
}

func copyFile(source, destination string, mode os.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

// allocatePair selects the logical pair. Overwrite uses the requested names.
// Aggregate retries complete pairs so a numbered primary is never separated
// from independently numbered informative output. Reject requires both names
// to be absent and reserves them together.
func (o *reportOutput) allocatePair() (*pairAllocation, error) {
	switch o.mode {
	case reportOverwrite:
		return &pairAllocation{
			primary:   o.path,
			companion: o.companion.path,
		}, nil
	case reportAggregate:
		for number := 0; ; number++ {
			primary := o.path
			if number > 0 {
				primary = numberedReportPath(o.path, number)
			}
			companion := informativeCompanionPath(primary)
			if err := o.validatePairStatic(primary, companion); err != nil {
				return nil, err
			}
			allocation, err := o.reservePair(primary, companion)
			if err == nil {
				return allocation, nil
			}
			if !errors.Is(err, fs.ErrExist) {
				return nil, err
			}
		}
	case reportReject:
		allocation, err := o.reservePair(o.path, o.companion.path)
		if errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf(
				"report pair %q and %q must both be absent when report output policy is reject: %w",
				o.path,
				o.companion.path,
				err,
			)
		}
		return allocation, err
	default:
		return nil, fmt.Errorf("unsupported report mode %s", o.mode)
	}
}

// reservePair claims both filenames with exclusive creates. If the second
// reservation fails, the first empty reservation is removed before the error is
// returned so aggregate allocation can try the next complete pair.
func (o *reportOutput) reservePair(primary, companion string) (*pairAllocation, error) {
	primaryPath, primaryInfo, err := reserveReportPath(primary)
	if err != nil {
		return nil, err
	}
	companionPath, companionInfo, err := reserveReportPath(companion)
	if err != nil {
		releaseReservedPath(primaryPath, primaryInfo)
		return nil, err
	}
	return &pairAllocation{
		primary:           primaryPath,
		companion:         companionPath,
		primaryReserved:   primaryInfo,
		companionReserved: companionInfo,
	}, nil
}

// writeReportTemp renders a report to a restrictive temporary file in the
// destination directory. The caller owns the returned filename and must rename
// it or remove it.
func writeReportTemp(reportPath, label string, write func(io.WriteCloser) error) (string, error) {
	temp, err := createReportTemp(reportPath)
	if err != nil {
		return "", fmt.Errorf("could not create temporary %s: %w", label, err)
	}
	tempName := temp.Name()
	cleanup := func() {
		_ = temp.Close()
		_ = os.Remove(tempName)
	}

	if runtime.GOOS != "windows" {
		if err := temp.Chmod(0o600); err != nil {
			cleanup()
			return "", fmt.Errorf("could not set %s permissions: %w", label, err)
		}
	}
	if err := write(temp); err != nil {
		cleanup()
		return "", fmt.Errorf("could not write %s: %w", label, err)
	}
	if err := temp.Sync(); err != nil {
		cleanup()
		return "", fmt.Errorf("could not flush %s: %w", label, err)
	}
	if err := temp.Close(); err != nil {
		_ = os.Remove(tempName)
		return "", fmt.Errorf("could not close %s: %w", label, err)
	}
	return tempName, nil
}
