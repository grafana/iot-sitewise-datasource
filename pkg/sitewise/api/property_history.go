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
func historyQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueHistoryInput {

	var (
		qualities []*string
	)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	//if propertyAlias is set make sure to set the assetId and propertyId to nil
	if query.PropertyAlias != "" {
		query.AssetIds = []string{}
		// nolint:staticcheck
		query.AssetId = ""
		query.PropertyId = ""
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 250 {
		query.MaxDataPoints = 250
	}

	entries := make([]*iotsitewise.BatchGetAssetPropertyValueHistoryEntry, 0)

	for _, id := range query.AssetIds {
		var assetId *string
		if id != "" {
			assetId = aws.String(id)
		}
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueHistoryEntry{
			StartDate:     from,
			EndDate:       to,
			EntryId:       assetId,
			PropertyId:    aws.String(query.PropertyId),
			AssetId:       assetId,
			PropertyAlias: getPropertyAlias(query.BaseQuery),
			TimeOrdering:  aws.String(query.TimeOrdering),
			Qualities:     qualities,
		})
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryInput{
		Entries:    entries,
		MaxResults: aws.Int64(query.MaxDataPoints),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValues(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyValueHistory, error) {
	var (
		maxDps = int(query.MaxDataPoints)
	)

	awsReq := historyQueryToInput(query)
	resp, err := client.BatchGetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyValueHistory{
		BatchGetAssetPropertyValueHistoryOutput: resp,
		Query:                                   query,
	}, nil
}
