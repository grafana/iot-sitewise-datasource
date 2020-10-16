package fdata

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

type AssetPropertyValue iotsitewise.GetAssetPropertyValueOutput

func (a AssetPropertyValue) Rows() [][]interface{} {
	rows := [][]interface{}{
		{getTimeValue(a.PropertyValue.Timestamp), getPropertyVariantValue(a.PropertyValue.Value)},
	}

	fmt.Println(rows)
	return rows
}

func getTimeValue(ts *iotsitewise.TimeInNanos) time.Time {
	var sec int64 = 0
	var nsec int64 = 0
	if ts.TimeInSeconds != nil {
		sec = *ts.TimeInSeconds
	}
	if ts.OffsetInNanos != nil {
		nsec = *ts.OffsetInNanos
	}
	return time.Unix(sec, nsec)
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
