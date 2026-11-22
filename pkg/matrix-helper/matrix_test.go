package matrixhelper

import (
	"testing"
	"time"
)

func Test_NeedsRetentionUpdate(t *testing.T) {
	sevenDays := 7 * 24 * time.Hour
	sixDays := 6 * 24 * time.Hour
	eightDays := 8 * 24 * time.Hour
	oneHour := time.Hour
	zero := 0 * time.Hour
	tests := []struct {
		name    string
		current *time.Duration
		desired time.Duration
		want    bool
	}{
		{"absent_policy_applies", nil, sevenDays, true},
		{"absent_policy_applies_even_for_zero_desired", nil, zero, true},
		{"equal_policy_skips", &sevenDays, sevenDays, false},
		{"equal_zero_policy_skips", &zero, zero, false},
		{"shorter_policy_applies", &sixDays, sevenDays, true},
		{"longer_policy_applies", &eightDays, sevenDays, true},
		{"manual_drift_is_corrected", &oneHour, sevenDays, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NeedsRetentionUpdate(tt.current, tt.desired); got != tt.want {
				t.Fatalf("expected update=%v, got %v", tt.want, got)
			}
		})
	}
}
