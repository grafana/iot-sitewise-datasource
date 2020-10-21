package framer

import (
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func fieldTypeForPropertyValue(property *iotsitewise.DescribeAssetPropertyOutput) data.FieldType {
	switch *property.AssetProperty.DataType {
	case "BOOLEAN":
		return data.FieldTypeNullableBool
	case "INTEGER":
		return data.FieldTypeNullableInt64
	case "STRING":
		return data.FieldTypeNullableString
	default:
		return data.FieldTypeNullableFloat64
	}
}

func getTime(ts *iotsitewise.TimeInNanos) time.Time {
	sec := *ts.TimeInSeconds

	if nanos := ts.OffsetInNanos; nanos != nil {
		return time.Unix(sec, *nanos)
	}

	return time.Unix(sec, 0)
}

func getPropertyVariantValue(variant *iotsitewise.Variant) interface{} {

	if val := variant.BooleanValue; val != nil {
		return val
	}

	if val := variant.DoubleValue; val != nil {
		return val
	}

	if val := variant.IntegerValue; val != nil {
		return val
	}

	if val := variant.StringValue; val != nil {
		return val
	}

	return nil
}

func newPropertyValueField(property *iotsitewise.DescribeAssetPropertyOutput, length int) *data.Field {
	valueField := newFieldWithName(*property.AssetProperty.Name, fieldTypeForPropertyValue(property), length)
	valueField.Config = &data.FieldConfig{
		Unit: toGrafanaUnit(property.AssetProperty.Unit),
	}
	return valueField
}

func newFieldWithName(name string, fieldType data.FieldType, length int) *data.Field {
	field := data.NewFieldFromFieldType(fieldType, length)
	field.Name = name
	return field
}

// Map values from ???:
//   https://docs.microsoft.com/en-us/rest/api/monitor/metrics/list#unit
// to
//   https://github.com/grafana/grafana/blob/master/packages/grafana-data/src/valueFormats/categories.ts#L24
func toGrafanaUnit(unit *string) string {
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
