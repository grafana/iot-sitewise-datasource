package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

var (
	MaxSitewiseResults = aws.Int64(250)
)

func getNextToken(query models.BaseQuery) *string {
	if query.NextToken == "" {
		return nil
	}
	return aws.String(query.NextToken)
}

func getAssetId(query models.BaseQuery) *string {
	if len(query.AssetIds) == 0 {
		return nil
	}
	return aws.String(query.AssetIds[0])
}

func getPropertyId(query models.BaseQuery) *string {
	if query.PropertyId == "" {
		return nil
	}
	return aws.String(query.PropertyId)
}

func getPropertyAlias(query models.BaseQuery) *string {
	if query.PropertyAlias == "" {
		return nil
	}
	return aws.String(query.PropertyAlias)
}

func getAssetIdAndPropertyId(query models.AssetPropertyValueQuery, client client.SitewiseClient, ctx context.Context) (models.AssetPropertyValueQuery, error) {
	result := query
	if query.PropertyAlias != "" {
		resp, err := client.DescribeTimeSeriesWithContext(ctx, &iotsitewise.DescribeTimeSeriesInput{
			Alias: getPropertyAlias(query.BaseQuery),
		})
		if err != nil {
			return models.AssetPropertyValueQuery{}, err
		}
		result.AssetIds = []string{*resp.AssetId}
		result.PropertyId = *resp.PropertyId
	}
	return result, nil
}
