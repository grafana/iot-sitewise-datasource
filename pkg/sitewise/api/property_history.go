package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

// GetAssetPropertyValueHistory requires either PropertyAlias OR (AssetID and PropertyID) to be set.
// The front end component should ensure that both cannot be sent at the same time by the user.
// If an invalid combo of assetId/propertyId/propertyAlias are sent to the API, an exception will be returned.
// The Framer consumer should bubble up that error to the user.
func historyQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueHistoryInput {
	quality := query.Quality
	if quality == "" || quality == "ANY" {
		quality = iotsitewisetypes.QualityGood
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 20000 {
		query.MaxDataPoints = 20000
	}

	return &iotsitewise.GetAssetPropertyValueHistoryInput{
		StartDate:     from,
		EndDate:       to,
		MaxResults:    aws.Int32(int32(query.MaxDataPoints)),
		NextToken:     getNextToken(query.BaseQuery),
		AssetId:       getFirstAssetId(query.BaseQuery),
		PropertyId:    getFirstPropertyId(query.BaseQuery),
		PropertyAlias: getFirstPropertyAlias(query.BaseQuery),
		TimeOrdering:  query.TimeOrdering,
		Qualities:     []iotsitewisetypes.Quality{quality},
	}
}

func GetAssetPropertyValues(ctx context.Context, sw client.SitewiseAPIClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValueHistory, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, sw, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReq := historyQueryToInput(modifiedQuery)
	resp, err := sw.GetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)

	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	return modifiedQuery,
		&framer.AssetPropertyValueHistory{
			GetAssetPropertyValueHistoryOutput: resp,
			Query:                              modifiedQuery,
		},
		nil
}
