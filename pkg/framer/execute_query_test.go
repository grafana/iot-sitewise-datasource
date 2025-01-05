package framer

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetValue(t *testing.T) {
	tests := []struct {
		name        string
		col         *iotsitewise.ColumnInfo
		scalarValue []string
		expected    interface{}
		expectError bool
	}{
		{
			name: "BOOLEAN true",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("BOOLEAN")},
			},
			scalarValue: []string{"true"},
			expected:    true,
			expectError: false,
		},
		{
			name: "BOOLEAN false",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("BOOLEAN")},
			},
			scalarValue: []string{"false"},
			expected:    false,
			expectError: false,
		},
		{
			name: "INT",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("INT")},
			},
			scalarValue: []string{"123"},
			expected:    int64(123),
			expectError: false,
		},
		{
			name: "INTEGER",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("INTEGER")},
			},
			scalarValue: []string{"123"},
			expected:    int64(123),
			expectError: false,
		},
		{
			name: "STRING",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("STRING")},
			},
			scalarValue: []string{"test"},
			expected:    "test",
			expectError: false,
		},
		{
			name: "DOUBLE",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("DOUBLE")},
			},
			scalarValue: []string{"123.456"},
			expected:    123.456,
			expectError: false,
		},
		{
			name: "TIMESTAMP",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("TIMESTAMP")},
			},
			scalarValue: []string{"1736116323"},
			expected:    time.Date(2025, time.January, 5, 22, 32, 03, 0, time.Local),
			expectError: false,
		},
		{
			name: "Unsupported type",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("UNSUPPORTED")},
			},
			scalarValue: []string{"test"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid BOOLEAN",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("BOOLEAN")},
			},
			scalarValue: []string{"notabool"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid INT",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("INT")},
			},
			scalarValue: []string{"notanint"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid INTEGER",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("INTEGER")},
			},
			scalarValue: []string{"notanint"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid DOUBLE",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("DOUBLE")},
			},
			scalarValue: []string{"notadouble"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid TIMESTAMP",
			col: &iotsitewise.ColumnInfo{
				Name: aws.String("Test Field"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("TIMESTAMP")},
			},
			scalarValue: []string{"notatimestamp"},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := fields.DatumField(*tt.col, 1)
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
