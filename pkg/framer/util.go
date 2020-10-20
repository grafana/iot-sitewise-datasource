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
