package fields

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func FieldTypeForPropertyValue(property *iotsitewise.DescribeAssetPropertyOutput) data.FieldType {
	switch *property.AssetProperty.DataType {
	case "BOOLEAN":
		return data.FieldTypeBool
	case "INTEGER":
		return data.FieldTypeInt64
	case "STRING":
		return data.FieldTypeString
	default:
		return data.FieldTypeFloat64
	}
}

// Map values from ???:
//   https://docs.microsoft.com/en-us/rest/api/monitor/metrics/list#unit
// to
//   https://github.com/grafana/grafana/blob/master/packages/grafana-data/src/valueFormats/categories.ts#L24
func ToGrafanaUnit(unit *string) string {
	if unit == nil {
		return ""
	}

	switch *unit {
	case "BitsPerSecond":
		return "bps"
	case "Bytes":
		return "decbytes" // or ICE
	case "BytesPerSecond":
		return "Bps"
	case "Count":
		return "short" // this is used for integers
	case "CountPerSecond":
		return "cps"
	case "Percent":
		return "percent"
	case "Milliseconds":
		return "ms"
	case "Seconds":
		return "s"
	}
	return *unit // this will become a suffix in the display
}
