package util

import (
	"testing"
	"time"
)

func TestGetFormattedTimeRange(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC time",
			input:    time.Date(2023, 1, 1, 10, 20, 30, 0, time.UTC),
			expected: "2023-01-01 10:20:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetFormattedTimeRange(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}
