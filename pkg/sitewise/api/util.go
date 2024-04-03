package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

const (
	l4eAnomalyResultPropertyName = "AWS/L4E_ANOMALY_RESULT"
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

func getAssetIdAndPropertyId(query models.AssetPropertyValueQuery, client client.SitewiseClient, ctx context.Context) (models.AssetPropertyValueQuery, error) {
	result := query
	if query.PropertyAlias != "" {
		resp, err := client.DescribeTimeSeriesWithContext(ctx, &iotsitewise.DescribeTimeSeriesInput{
			Alias: util.GetPropertyAlias(query.BaseQuery),
		})
		if err != nil {
			return models.AssetPropertyValueQuery{}, err
		}
		if resp.AssetId != nil {
			result.AssetIds = []string{*resp.AssetId}
		} else {
			// For disassociated streams with a propertyAlias
			result.AssetIds = []string{}
		}
		if resp.PropertyId != nil {
			result.PropertyId = *resp.PropertyId
		} else {
			// For disassociated streams without a propertyAlias
			result.PropertyId = ""
		}
	}
	return result, nil
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

func filterAnomalyAssetIds(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) ([]string, error) {
	anomalyAssetIds := []string{}

	switch {
	case query.PropertyAlias != "":
		return nil, nil

	default:
		for _, assetId := range query.AssetIds {
			var id *string
			if assetId != "" {
				id = aws.String(assetId)
			}

			req := &iotsitewise.DescribeAssetPropertyInput{
				AssetId:    id,
				PropertyId: aws.String(query.PropertyId),
			}

			resp, err := client.DescribeAssetPropertyWithContext(ctx, req)
			if err != nil {
				return nil, err
			}

			if resp.CompositeModel != nil && *resp.CompositeModel.AssetProperty.Name == l4eAnomalyResultPropertyName {
				anomalyAssetIds = append(anomalyAssetIds, assetId)
			}
		}
	}

	return anomalyAssetIds, nil
}
