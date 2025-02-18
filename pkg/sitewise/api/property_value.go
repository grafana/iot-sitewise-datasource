package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueInput {

	return &iotsitewise.GetAssetPropertyValueInput{
		AssetId:       getAssetId(query.BaseQuery),
		PropertyId:    getPropertyId(query.BaseQuery),
		PropertyAlias: getPropertyAlias(query.BaseQuery),
	}
}

func GetAssetPropertyValue(ctx context.Context, sw client.SitewiseAPIClient, query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValue, error) {
	modifiedQuery, err := getAssetIdAndPropertyId(query, sw, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReq := valueQueryToInput(modifiedQuery)

	resp, err := sw.GetAssetPropertyValue(ctx, awsReq)

	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	return modifiedQuery, &framer.AssetPropertyValue{PropertyValue: resp.PropertyValue}, nil
}
