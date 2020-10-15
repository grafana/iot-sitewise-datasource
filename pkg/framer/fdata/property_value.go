package fdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

type AssetPropertyValue iotsitewise.GetAssetPropertyValueOutput

func (a AssetPropertyValue) Rows() [][]interface{} {
	rows := [][]interface{}{
		{getTimeInMs(a.PropertyValue.Timestamp), getPropertyVariantValue(a.PropertyValue.Value)},
	}

	fmt.Println(rows)
	return rows
}

func getTimeInMs(ts *iotsitewise.TimeInNanos) int64 {

	secMs := *ts.TimeInSeconds * 1e3

	if nanos := ts.OffsetInNanos; nanos != nil {
		nanosMs := *ts.OffsetInNanos / 1e6
		secMs = secMs + nanosMs
	}
	return secMs
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
