package sitewise

import (
	"errors"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
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
			name:  "TimeFrom",
			macro: "timeFrom",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "TIMESTAMP '2023-01-01 00:00:00'",
			expectedErr: nil,
		},
		// RawTimeFrom
		{
			name:  "TimeTo default format",
			macro: "timeTo",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "TIMESTAMP '2023-01-02 00:00:00'",
			expectedErr: nil,
		},
		// UnixEpochFilter
		{
			name:  "TimeFilter valid case",
			macro: "timeFilter",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{"timestamp"},
			expected:    "timestamp >= TIMESTAMP '2023-01-01 00:00:00' and timestamp <= TIMESTAMP '2023-01-02 00:00:00'",
			expectedErr: nil,
		},
		{
			name:  "TimeFilter invalid argument count",
			macro: "timeFilter",
			query: &sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			args:        []string{},
			expected:    "",
			expectedErr: ErrorBadArgumentCount,
		},
		// resolution
		{
			name:  "resolution less than 1m",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(30 * time.Second),
			},
			args:        []string{},
			expected:    "1m",
			expectedErr: nil,
		},
		{
			name:  "resolution 1m",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(60 * time.Second),
			},
			args:        []string{},
			expected:    "1m",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 15m",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(90 * time.Second),
			},
			args:        []string{},
			expected:    "15m",
			expectedErr: nil,
		},
		{
			name:  "resolution 15m",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(900 * time.Second),
			},
			args:        []string{},
			expected:    "15m",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 1h",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(1000 * time.Second),
			},
			args:        []string{},
			expected:    "1h",
			expectedErr: nil,
		},
		{
			name:  "resolution 1h",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(3600 * time.Second),
			},
			args:        []string{},
			expected:    "1h",
			expectedErr: nil,
		},
		{
			name:  "resolution less than 1d",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(4000 * time.Second),
			},
			args:        []string{},
			expected:    "1d",
			expectedErr: nil,
		},
		{
			name:  "resolution 1d",
			macro: "autoResolution",
			query: &sqlutil.Query{
				Interval: time.Duration(86400 * time.Second),
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
