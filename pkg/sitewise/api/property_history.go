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
		propertyId *string
		assetId    *string
		nextToken  *string
		qualities  []*string
	)

	assetId = getAssetId(query.BaseQuery)
	propertyId = getPropertyId(query.BaseQuery)
	nextToken = getNextToken(query.BaseQuery)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	return &iotsitewise.GetAssetPropertyValueHistoryInput{
		AssetId:    assetId,
		StartDate:  from,
		EndDate:    to,
		MaxResults: aws.Int64(250), // should this even be configurable? 250 == max
		NextToken:  nextToken,
		PropertyId: propertyId,
		Qualities:  qualities,
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
		AssetPropertyValueHistory: resp.AssetPropertyValueHistory,
		NextToken:                 resp.NextToken,
	}, nil
}
