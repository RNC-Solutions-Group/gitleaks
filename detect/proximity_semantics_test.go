package detect

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/report"
)

func intPtr(v int) *int { return &v }

func TestWithinProximityBoundaries(t *testing.T) {
	d := &Detector{}
	primary := report.Finding{StartLine: 10, StartColumn: 5}

	tests := []struct {
		name     string
		required report.Finding
		rule     config.Required
		want     bool
	}{
		{
			name:     "no constraints is fragment-level",
			required: report.Finding{StartLine: 100, StartColumn: 100},
			rule:     config.Required{},
			want:     true,
		},
		{
			name:     "line boundary equal matches",
			required: report.Finding{StartLine: 14, StartColumn: 1},
			rule:     config.Required{WithinLines: intPtr(4)},
			want:     true,
		},
		{
			name:     "line boundary exceeded fails",
			required: report.Finding{StartLine: 15, StartColumn: 1},
			rule:     config.Required{WithinLines: intPtr(4)},
			want:     false,
		},
		{
			name:     "column zero requires the same start column",
			required: report.Finding{StartLine: 10, StartColumn: 5},
			rule:     config.Required{WithinColumns: intPtr(0)},
			want:     true,
		},
		{
			name:     "column zero rejects an adjacent column",
			required: report.Finding{StartLine: 10, StartColumn: 6},
			rule:     config.Required{WithinColumns: intPtr(0)},
			want:     false,
		},
		{
			name:     "column boundary equal matches",
			required: report.Finding{StartLine: 10, StartColumn: 8},
			rule:     config.Required{WithinColumns: intPtr(3)},
			want:     true,
		},
		{
			name:     "column boundary exceeded fails",
			required: report.Finding{StartLine: 10, StartColumn: 9},
			rule:     config.Required{WithinColumns: intPtr(3)},
			want:     false,
		},
		{
			name:     "column constraint cannot cross a line even when columns match",
			required: report.Finding{StartLine: 11, StartColumn: 5},
			rule:     config.Required{WithinColumns: intPtr(10)},
			want:     false,
		},
		{
			name:     "both constraints on the same line use both axes",
			required: report.Finding{StartLine: 11, StartColumn: 8},
			rule:     config.Required{WithinLines: intPtr(4), WithinColumns: intPtr(3)},
			want:     false,
		},
		{
			name:     "both constraints same line within both boundaries",
			required: report.Finding{StartLine: 10, StartColumn: 8},
			rule:     config.Required{WithinLines: intPtr(4), WithinColumns: intPtr(3)},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, d.withinProximity(primary, tt.required, &tt.rule))
		})
	}
}
