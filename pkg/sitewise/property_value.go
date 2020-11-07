package sitewise

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueInput {

	var (
		propertyId *string
		assetId    *string
	)

	assetId = getAssetId(query.BaseQuery)
	propertyId = getPropertyId(query.BaseQuery)

	return &iotsitewise.GetAssetPropertyValueInput{
		AssetId:    assetId,
		PropertyId: propertyId,
	}

}

func GetAssetPropertyValue(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyValue, error) {

	awsReq := valueQueryToInput(query)

	resp, err := client.GetAssetPropertyValueWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyValue{PropertyValue: resp.PropertyValue}, nil
}
