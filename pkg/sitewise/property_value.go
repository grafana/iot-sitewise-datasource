package sitewise

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fdata"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

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

func GetAssetPropertyValue(ctx context.Context, client client.Client, query models.AssetPropertyValueQuery) (*fdata.AssetPropertyValue, error) {

	awsReq := valueQueryToInput(query)

	resp, err := client.GetAssetPropertyValueWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &fdata.AssetPropertyValue{PropertyValue: resp.PropertyValue}, nil
}
