package sitewise

import (
	"errors"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/sqlds/v4"
)

func TestMacros(t *testing.T) {
	tests := []struct {
		name        string
		macro       string
		query       *sqlutil.Query
		args        []string
		expected    string
		expectedErr error
	}{
		// SelectAll
		{
			name:  "SelectAll valid case",
			macro: "selectAll",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				RawSQL: "SELECT $__selectAll FROM asset",
			},
			args:        []string{},
			expected:    "asset_id, asset_name, asset_description, asset_model_id",
			expectedErr: nil,
		},
		{
			name:  "SelectAll unknown table",
			macro: "selectAll",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				RawSQL: "SELECT $__selectAll FROM azzet",
			},
			args:        []string{},
			expected:    "selectAll",
			expectedErr: TableColumnsNotFoundError,
		},
		{
			name:  "SelectAll incomplete query",
			macro: "selectAll",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				RawSQL: "SELECT $__selectAll FRO",
			},
			args:        []string{},
			expected:    "selectAll",
			expectedErr: TableColumnsNotFoundError,
		},
		// RawTimeFrom
		{
			name:  "RawTimeFrom",
			macro: "rawTimeFrom",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "1672531200",
			expectedErr: nil,
		},
		// RawTimeFrom
		{
			name:  "RawTimeTo default format",
			macro: "rawTimeTo",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "1672617600",
			expectedErr: nil,
		},
		// UnixEpochFilter
		{
			name:  "UnixEpochFilter valid case",
			macro: "unixEpochFilter",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{"timestamp"},
			expected:    "timestamp >= 1672531200 and timestamp <= 1672617600",
			expectedErr: nil,
		},
		{
			name:  "UnixEpochFilter invalid argument count",
			macro: "unixEpochFilter",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "",
			expectedErr: sqlds.ErrorBadArgumentCount,
		},
		// resolution
		{
			name:  "resolution less than 1m",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 30,
			},
			args:        []string{},
			expected:    "1m",
			expectedErr: nil,
		},
		{
			name:  "resolution 1m",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 60,
			},
			args:        []string{},
			expected:    "1m",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 15m",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 90,
			},
			args:        []string{},
			expected:    "15m",
			expectedErr: nil,
		},
		{
			name:  "resolution 15m",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 900,
			},
			args:        []string{},
			expected:    "15m",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 1h",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 1000,
			},
			args:        []string{},
			expected:    "1h",
			expectedErr: nil,
		},
		{
			name:  "resolution 1h",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 3600,
			},
			args:        []string{},
			expected:    "1h",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 1d",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 4000,
			},
			args:        []string{},
			expected:    "1d",
			expectedErr: nil,
		},
		{
			name:  "resolution 1d",
			macro: "resolution",
			query: &sqlutil.Query{
				Interval: 86400,
			},
			args:        []string{},
			expected:    "1d",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := macros[tt.macro](tt.query, tt.args)
			if (err != nil || tt.expectedErr != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error %v, expecting %v", err, tt.expectedErr)
			}
			if res != tt.expected {
				t.Errorf("unexpected result %v, expecting %v", res, tt.expected)
			}
		})
	}
}
