package framer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrames(t *testing.T) {
	tests := []struct {
		name     string
		results  QueryResults
		expected int
	}{
		{
			name: "Empty results",
			results: QueryResults{
				Rows: []iotsitewisetypes.Row{},
			},
			expected: 0,
		},
		{
			name: "Single row",
			results: QueryResults{
				Rows: []iotsitewisetypes.Row{
					{
						Data: []iotsitewisetypes.Datum{
							{ScalarValue: aws.String("true")},
						},
					},
				},
				Columns: []iotsitewisetypes.ColumnInfo{
					{
						Name: aws.String("Test Field"),
						Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
					},
				},
			},
			expected: 1,
		},
		{
			name: "Multiple rows",
			results: QueryResults{
				Rows: []iotsitewisetypes.Row{
					{
						Data: []iotsitewisetypes.Datum{
							{ScalarValue: aws.String("true")},
						},
					},
					{
						Data: []iotsitewisetypes.Datum{
							{ScalarValue: aws.String("false")},
						},
					},
				},
				Columns: []iotsitewisetypes.ColumnInfo{
					{
						Name: aws.String("Test Field"),
						Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
					},
				},
			},
			expected: 2,
		},
		{
			name: "Null values",
			results: QueryResults{
				Rows: []iotsitewisetypes.Row{
					{
						Data: []iotsitewisetypes.Datum{
							{ScalarValue: nil},
						},
					},
					{
						Data: []iotsitewisetypes.Datum{
							{ScalarValue: aws.String("null")},
						},
					},
				},
				Columns: []iotsitewisetypes.ColumnInfo{
					{
						Name: aws.String("Test Field"),
						Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
					},
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frames, err := tt.results.Frames(context.Background(), nil)
			require.NoError(t, err)
			require.Len(t, frames, 1)
			if tt.expected > 0 {
				assert.Equal(t, tt.expected, frames[0].Fields[0].Len())
			} else {
				assert.Equal(t, 0, len(frames[0].Fields))
				backend.Logger.Debug("weirdness", "frames", frames)
			}
		})
	}
}

func TestSetValue(t *testing.T) {
	tests := []struct {
		name        string
		col         iotsitewisetypes.ColumnInfo
		scalarValue []string
		expected    interface{}
		expectError bool
	}{
		{
			name: "BOOLEAN true",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
			},
			scalarValue: []string{"true"},
			expected:    true,
			expectError: false,
		},
		{
			name: "BOOLEAN false",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
			},
			scalarValue: []string{"false"},
			expected:    false,
			expectError: false,
		},
		{
			name: "INT",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeInt},
			},
			scalarValue: []string{"123"},
			expected:    int64(123),
			expectError: false,
		},
		{
			name: "INTEGER",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeInt},
			},
			scalarValue: []string{"123"},
			expected:    int64(123),
			expectError: false,
		},
		{
			name: "STRING",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeString},
			},
			scalarValue: []string{"test"},
			expected:    "test",
			expectError: false,
		},
		{
			name: "DOUBLE",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeDouble},
			},
			scalarValue: []string{"123.456"},
			expected:    123.456,
			expectError: false,
		},
		//{
		//	// FIXME: this test fails due either to TZ locale or daylight savings time
		//	name: "TIMESTAMP",
		//	col: iotsitewisetypes.ColumnInfo{
		//		Name: aws.String("Test Field"),
		//		Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeTimestamp},
		//	},
		//	scalarValue: []string{"1736116323"},
		//	expected:    time.Date(2025, time.January, 5, 22, 32, 03, 0, time.Local),
		//	expectError: false,
		//},
		{
			name: "Unsupported type",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarType("UNSUPPORTED")},
			},
			scalarValue: []string{"test"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid BOOLEAN",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeBoolean},
			},
			scalarValue: []string{"notabool"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid INT",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarType("INT")},
			},
			scalarValue: []string{"notanint"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid INTEGER",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeInt},
			},
			scalarValue: []string{"notanint"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid DOUBLE",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeDouble},
			},
			scalarValue: []string{"notadouble"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid TIMESTAMP",
			col: iotsitewisetypes.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeTimestamp},
			},
			scalarValue: []string{"notatimestamp"},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := fields.DatumField(1, tt.col)
			err := SetValue(tt.col, tt.scalarValue[0], field, 0)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, field.At(0))
			}
		})
	}
}
