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
// `query.MaxDataPoints` is ignored and it always requests with the maximum number of data points the SiteWise API can support
func historyBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueHistoryInput {

	var (
		qualities []*string
	)

	quality := query.Quality

	if quality == "" || quality == "ANY" {
		qualities = aws.StringSlice([]string{"GOOD"})
	} else {
		qualities = aws.StringSlice([]string{quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	entries := make([]*iotsitewise.BatchGetAssetPropertyValueHistoryEntry, 0)

	switch {
	case query.PropertyAlias != "":
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueHistoryEntry{
			StartDate:     from,
			EndDate:       to,
			EntryId:       util.GetEntryId(query.BaseQuery),
			PropertyAlias: util.GetPropertyAlias(query.BaseQuery),
			TimeOrdering:  aws.String(query.TimeOrdering),
			Qualities:     qualities,
		})
	default:
		for _, id := range query.AssetIds {
			var assetId *string
			if id != "" {
				assetId = aws.String(id)
			}
			entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueHistoryEntry{
				StartDate:    from,
				EndDate:      to,
				EntryId:      assetId,
				AssetId:      assetId,
				PropertyId:   aws.String(query.PropertyId),
				TimeOrdering: aws.String(query.TimeOrdering),
				Qualities:    qualities,
			})
		}
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryInput{
		Entries: entries,
		// performance: hardcoded to fetch the maximum number of data points
		MaxResults: aws.Int64(BatchGetAssetPropertyValueHistoryMaxResults),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValues(ctx context.Context, client client.SitewiseClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValueHistoryBatch, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	batchedQueries := batchQueries(modifiedQuery, BatchGetAssetPropertyValueHistoryMaxEntries)
	responses := []*iotsitewise.BatchGetAssetPropertyValueHistoryOutput{}
	for _, q := range batchedQueries {
		awsReq := historyBatchQueryToInput(q)
		resp, err := client.BatchGetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)
		if err != nil {
			return models.AssetPropertyValueQuery{}, nil, err
		}
		responses = append(responses, resp)
	}

	anomalyAssetIds := []string{}
	if query.FlattenL4e {
		anomalyAssetIds, err = filterAnomalyAssetIds(ctx, client, modifiedQuery)
		if err != nil {
			return models.AssetPropertyValueQuery{}, nil, err
		}
	}

	return modifiedQuery,
		&framer.AssetPropertyValueHistoryBatch{
			Responses:       responses,
			Query:           modifiedQuery,
			AnomalyAssetIds: anomalyAssetIds,
			SitewiseClient:  client,
		},
		nil
}
