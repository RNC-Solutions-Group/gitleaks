package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateFamiliesOverlap(t *testing.T) {
	dir := t.TempDir()
	other := t.TempDir()
	path := func(root, name string) string { return filepath.Join(root, name) }

	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{name: "same base", a: path(dir, "report.json"), b: path(dir, "report.json"), want: true},
		{name: "first candidate", a: path(dir, "report.json"), b: path(dir, "report_001.json"), want: true},
		{name: "v14a boundary", a: path(dir, "report.json"), b: path(dir, "report_1000.json"), want: true},
		{name: "beyond v14a boundary", a: path(dir, "report.json"), b: path(dir, "report_1001.json"), want: true},
		{name: "high candidate", a: path(dir, "report.json"), b: path(dir, "report_10001.json"), want: true},
		{name: "noncanonical padding", a: path(dir, "report.json"), b: path(dir, "report_0001.json"), want: false},
		{name: "noncanonical width", a: path(dir, "report.json"), b: path(dir, "report_1.json"), want: false},
		{name: "nested family member", a: path(dir, "report_001.json"), b: path(dir, "report_001_001.json"), want: true},
		{name: "reverse nested family", a: path(dir, "report_001_001.json"), b: path(dir, "report_001.json"), want: true},
		{name: "informative name", a: path(dir, "report.json"), b: path(dir, "report-informative.json"), want: false},
		{name: "double numbered suffix", a: path(dir, "report.json"), b: path(dir, "report_001_001.json"), want: false},
		{name: "different directory", a: path(dir, "report.json"), b: path(other, "report_001.json"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, aggregateFamiliesOverlap(test.a, test.b))
		})
	}
}

func TestCompanionFamilyContains(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "report.json")

	tests := []struct {
		name   string
		member string
		want   bool
	}{
		{name: "base companion", member: filepath.Join(dir, "report-informative.json"), want: true},
		{name: "candidate one", member: filepath.Join(dir, "report_001-informative.json"), want: true},
		{name: "v14a boundary", member: filepath.Join(dir, "report_1000-informative.json"), want: true},
		{name: "beyond v14a boundary", member: filepath.Join(dir, "report_1001-informative.json"), want: true},
		{name: "noncanonical padding", member: filepath.Join(dir, "report_0001-informative.json"), want: false},
		{name: "noncanonical width", member: filepath.Join(dir, "report_1-informative.json"), want: false},
		{name: "grouped aggregate member", member: filepath.Join(dir, "report_001-informative_001.json"), want: false},
		{name: "ordinary primary candidate", member: filepath.Join(dir, "report_001.json"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, companionFamilyContains(primary, test.member))
		})
	}
}

func TestReportDestinationSetUsesStructuralAggregateChecks(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "report.json")
	companion := informativeCompanionPath(primary)

	require.NoError(t, validateReportDestinationSet(primary, companion, filepath.Join(dir, "grouped.json"), reportAggregate))
	require.NoError(t, validateReportDestinationSet(primary, companion, filepath.Join(dir, "report_001-informative_001.json"), reportAggregate))

	for _, grouped := range []string{
		filepath.Join(dir, "report_1001.json"),
		filepath.Join(dir, "report_1001-informative.json"),
	} {
		require.Error(t, validateReportDestinationSet(primary, companion, grouped, reportAggregate), grouped)
	}
}
