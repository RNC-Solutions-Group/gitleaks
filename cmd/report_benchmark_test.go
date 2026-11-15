package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkReportDestinationValidation(b *testing.B) {
	dir := b.TempDir()
	primary := filepath.Join(dir, "report.json")
	companion := informativeCompanionPath(primary)
	grouped := filepath.Join(dir, "grouped.json")

	cases := []struct {
		name      string
		companion string
		grouped   string
	}{
		{name: "primary_only"},
		{name: "primary_grouped", grouped: grouped},
		{name: "primary_companion_grouped", companion: companion, grouped: grouped},
	}
	for _, test := range cases {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := validateReportDestinationSet(primary, test.companion, test.grouped, reportAggregate); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkAggregateAllocation(b *testing.B) {
	for _, occupied := range []int{0, 10, 100, 1_000, 10_000} {
		b.Run(fmt.Sprintf("occupied_%d", occupied), func(b *testing.B) {
			dir := b.TempDir()
			dest := filepath.Join(dir, "report.json")
			seed := filepath.Join(dir, "occupied.seed")
			if occupied > 0 {
				if err := os.WriteFile(seed, []byte("OCCUPIED"), 0o600); err != nil {
					b.Fatal(err)
				}
				candidates := []string{dest}
				for number := 1; number <= occupied; number++ {
					candidates = append(candidates, numberedReportPath(dest, number))
				}
				for _, candidate := range candidates {
					if err := os.Link(seed, candidate); err != nil {
						if err := os.WriteFile(candidate, []byte("OCCUPIED"), 0o600); err != nil {
							b.Fatal(err)
						}
					}
				}
			}

			out, err := newReportOutput(dest, reportAggregate, "", nil, "", ".", "")
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				path, err := out.write(&stringReporter{payload: "bench"}, nil)
				if err != nil {
					b.Fatal(err)
				}
				b.StopTimer()
				if err := os.Remove(path); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
			}
		})
	}
}
