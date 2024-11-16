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
				RawSQL: "SELECT * FROM asset",
			},
			args:        []string{},
			expected:    "asset_id, asset_name, asset_description, asset_model_id, asset_root_id",
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
				RawSQL: "SELECT * FROM azzet",
			},
			args:        []string{},
			expected:    "*",
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
				RawSQL: "SELECT * FRO",
			},
			args:        []string{},
			expected:    "*",
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
