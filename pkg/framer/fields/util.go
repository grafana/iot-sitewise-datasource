package fields

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func FieldTypeForPropertyValue(property *iotsitewise.DescribeAssetPropertyOutput) data.FieldType {
	dataType := util.GetPropertyDataType(property)

	switch dataType {
	case "BOOLEAN":
		return data.FieldTypeBool
	case "INTEGER":
		return data.FieldTypeInt64
	case "INT":
		return data.FieldTypeInt64
	case "STRING":
		return data.FieldTypeString
	case "STRUCT":
		return data.FieldTypeString
	case "DOUBLE":
		return data.FieldTypeFloat64
	case "TIMESTAMP":
		return data.FieldTypeTime
	default:
		return data.FieldTypeString
	}
}

func FieldTypeForQueryResult(column iotsitewise.ColumnInfo) data.FieldType {
	// Override the type for event_timestamp
	if *column.Name == "event_timestamp" {
		return data.FieldTypeTime
	}

	switch *column.Type.ScalarType {
	case "BOOLEAN":
		return data.FieldTypeBool
	case "INTEGER":
		return data.FieldTypeInt64
	case "INT":
		return data.FieldTypeInt64
	case "STRING":
		return data.FieldTypeString
	case "DOUBLE":
		return data.FieldTypeFloat64
	case "TIMESTAMP":
		return data.FieldTypeTime
	default:
		backend.Logger.Debug("Unknown scalar type", "type", *column.Type.ScalarType)
		return data.FieldTypeString
	}
}

// Map values from ???:
//
//	https://docs.microsoft.com/en-us/rest/api/monitor/metrics/list#unit
//
// to
//
//	https://github.com/grafana/grafana/blob/master/packages/grafana-data/src/valueFormats/categories.ts#L24
func ToGrafanaUnit(unit *string) string {
	if unit == nil {
		return ""
	}

	switch *unit {
	case "Watts":
		return "watt"
	case "Kilowatts":
		return "kwatt"
	case "Count":
		return "short" // this is used for integers
	case "Percent":
		return "percent"
	case "Milliseconds":
		return "ms"
	case "Seconds":
		return "s"
	}
	return *unit // this will become a suffix in the display
}
