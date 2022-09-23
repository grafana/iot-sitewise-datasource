package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
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

	var (
		qualities []*string
	)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	//if propertyAlias is set make sure to set the assetId and propertyId to nil
	if query.PropertyAlias != "" {
		query.AssetId = ""
		query.PropertyId = ""
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 250 {
		query.MaxDataPoints = 250
	}

	return &iotsitewise.GetAssetPropertyValueHistoryInput{
		StartDate:     from,
		EndDate:       to,
		MaxResults:    aws.Int64(query.MaxDataPoints),
		NextToken:     getNextToken(query.BaseQuery),
		AssetId:       getAssetId(query.BaseQuery),
		PropertyId:    getPropertyId(query.BaseQuery),
		PropertyAlias: getPropertyAlias(query.BaseQuery),
		TimeOrdering:  aws.String(query.TimeOrdering),
		Qualities:     qualities,
	}
}

func GetAssetPropertyValues(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyValueHistory, error) {
	var (
		maxDps = int(query.MaxDataPoints)
	)

	awsReq := historyQueryToInput(query)
	resp, err := client.GetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyValueHistory{
		GetAssetPropertyValueHistoryOutput: resp,
		Query:                              query,
	}, nil
}
