package sitewise

import (
	"context"
	"fmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
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

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueInput {

	var (
		propertyId *string
		assetId    *string
	)

	if query.AssetId != "" {
		assetId = aws.String(query.AssetId)
	}

	if query.PropertyId != "" {
		propertyId = aws.String(query.PropertyId)
	}

	return &iotsitewise.GetAssetPropertyValueInput{
		AssetId:    assetId,
		PropertyId: propertyId,
	}

}

func GetAssetPropertyValue(ctx context.Context, client client.Client, query models.AssetPropertyValueQuery) (*AssetPropertyValue, error) {

	awsReq := valueQueryToInput(query)

	resp, err := client.GetAssetPropertyValueWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &AssetPropertyValue{PropertyValue: resp.PropertyValue}, nil
}
