package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameReportFileWithRetriesTransientFailure(t *testing.T) {
	transient := errors.New("transient rename failure")
	var (
		attempts int
		sleeps   []time.Duration
		oldPath  string
		newPath  string
	)
	err := renameReportFileWith(
		func(old, new string) error {
			attempts++
			oldPath, newPath = old, new
			if attempts < 3 {
				return transient
			}
			return nil
		},
		func(delay time.Duration) { sleeps = append(sleeps, delay) },
		func(err error) bool { return errors.Is(err, transient) },
		"report.tmp",
		"report.json",
	)
	require.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, []time.Duration{time.Millisecond, 2 * time.Millisecond}, sleeps)
	assert.Equal(t, "report.tmp", oldPath)
	assert.Equal(t, "report.json", newPath)
}

func TestRenameReportFileWithStopsOnPermanentFailure(t *testing.T) {
	permanent := errors.New("permanent rename failure")
	attempts := 0
	err := renameReportFileWith(
		func(string, string) error {
			attempts++
			return permanent
		},
		func(time.Duration) { t.Fatal("must not sleep for a permanent failure") },
		func(error) bool { return false },
		"report.tmp",
		"report.json",
	)
	require.ErrorIs(t, err, permanent)
	assert.Equal(t, 1, attempts)
}

func TestRenameReportFileWithExhaustsBoundedAttempts(t *testing.T) {
	transient := errors.New("always transient")
	attempts := 0
	sleeps := 0
	err := renameReportFileWith(
		func(string, string) error {
			attempts++
			return transient
		},
		func(time.Duration) { sleeps++ },
		func(error) bool { return true },
		"report.tmp",
		"report.json",
	)
	require.ErrorIs(t, err, transient)
	assert.Equal(t, reportRenameAttempts, attempts)
	assert.Equal(t, reportRenameAttempts-1, sleeps)
}
