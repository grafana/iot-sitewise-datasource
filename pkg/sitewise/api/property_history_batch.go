package api

import (
	"context"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

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

	qualities := []iotsitewisetypes.Quality{query.Quality}

	if query.Quality == "" || query.Quality == "ANY" {
		qualities[0] = iotsitewisetypes.QualityGood
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	entries := make([]iotsitewisetypes.BatchGetAssetPropertyValueHistoryEntry, 0)

	switch {
	case query.PropertyAlias != "":
		entries = append(entries, iotsitewisetypes.BatchGetAssetPropertyValueHistoryEntry{
			StartDate:     from,
			EndDate:       to,
			EntryId:       util.GetEntryId(query.BaseQuery),
			PropertyAlias: util.GetPropertyAlias(query.BaseQuery),
			TimeOrdering:  query.TimeOrdering,
			Qualities:     qualities,
		})
	default:
		for _, id := range query.AssetIds {
			var assetId *string
			if id != "" {
				assetId = aws.String(id)
			}
			entries = append(entries, iotsitewisetypes.BatchGetAssetPropertyValueHistoryEntry{
				StartDate:    from,
				EndDate:      to,
				EntryId:      assetId,
				AssetId:      assetId,
				PropertyId:   aws.String(query.PropertyId),
				TimeOrdering: query.TimeOrdering,
				Qualities:    qualities,
			})
		}
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryInput{
		Entries: entries,
		// performance: hardcoded to fetch the maximum number of data points
		MaxResults: aws.Int32(BatchGetAssetPropertyValueHistoryMaxResults),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValues(ctx context.Context, sw client.SitewiseAPIClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValueHistoryBatch, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, sw, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	batchedQueries := batchQueries(modifiedQuery, BatchGetAssetPropertyValueHistoryMaxEntries)
	responses := []*iotsitewise.BatchGetAssetPropertyValueHistoryOutput{}
	for _, q := range batchedQueries {
		awsReq := historyBatchQueryToInput(q)
		resp, err := sw.BatchGetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)
		if err != nil {
			return models.AssetPropertyValueQuery{}, nil, err
		}
		responses = append(responses, resp)
	}

	anomalyAssetIds := []string{}
	if query.FlattenL4e {
		anomalyAssetIds, err = filterAnomalyAssetIds(ctx, sw, modifiedQuery)
		if err != nil {
			return models.AssetPropertyValueQuery{}, nil, err
		}
	}

	return modifiedQuery,
		&framer.AssetPropertyValueHistoryBatch{
			Responses:       responses,
			Query:           modifiedQuery,
			AnomalyAssetIds: anomalyAssetIds,
			SitewiseClient:  sw,
		},
		nil
}
