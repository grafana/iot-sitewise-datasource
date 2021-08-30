package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueInput {

	return &iotsitewise.GetAssetPropertyValueInput{
		AssetId:       getAssetId(query.BaseQuery),
		PropertyId:    getPropertyId(query.BaseQuery),
		PropertyAlias: getPropertyAlias(query.BaseQuery),
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
