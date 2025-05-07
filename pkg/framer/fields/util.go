package fields

import (
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func FieldTypeForPropertyValue(property *iotsitewise.DescribeAssetPropertyOutput) data.FieldType {
	dataType := util.GetPropertyDataType(property)

	switch dataType {
	case types.PropertyDataTypeBoolean:
		return data.FieldTypeBool
	case types.PropertyDataTypeInteger:
		return data.FieldTypeInt64
	case types.PropertyDataTypeString:
		return data.FieldTypeString
	case types.PropertyDataTypeStruct:
		return data.FieldTypeString
	default:
		return data.FieldTypeFloat64
	}
}

func FieldTypeForQueryResult(column types.ColumnInfo) data.FieldType {
	// Override the type for event_timestamp
	if *column.Name == "event_timestamp" {
		return data.FieldTypeTime
	}

	switch column.Type.ScalarType {
	case types.ScalarTypeBoolean:
		return data.FieldTypeBool
	case types.ScalarTypeInt:
		return data.FieldTypeInt64
	case types.ScalarTypeString:
		return data.FieldTypeString
	case types.ScalarTypeDouble:
		return data.FieldTypeFloat64
	case types.ScalarTypeTimestamp:
		return data.FieldTypeTime
	default:
		backend.Logger.Debug("Unknown scalar type", "type", column.Type.ScalarType)
		return data.FieldTypeString
	}
}

// ToGrafanaUnit maps values from ???:
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
