package cmd

import "time"

// reportRenameAttempts bounds retries of a transient filesystem rename failure.
// The observed Windows failure was an Access is denied error while replacing a
// freshly reserved companion; an immediate retry succeeded. The retry is
// deliberately small and is applied only to errors classified by the
// platform-specific helper, not to every rename error.
const reportRenameAttempts = 6

type renameFileFunc func(oldPath, newPath string) error
type renameSleepFunc func(time.Duration)

// renameReportFileWith retries only errors selected by retryable. Backoff is
// bounded; callers still receive the final error unchanged.
func renameReportFileWith(rename renameFileFunc, sleep renameSleepFunc, retryable func(error) bool, oldPath, newPath string) error {
	var (
		err   error
		delay = time.Millisecond
	)
	for attempt := 0; attempt < reportRenameAttempts; attempt++ {
		err = rename(oldPath, newPath)
		if err == nil {
			return nil
		}
		if !retryable(err) || attempt == reportRenameAttempts-1 {
			return err
		}
		sleep(delay)
		if delay < 8*time.Millisecond {
			delay *= 2
		}
	}
	return err
}
