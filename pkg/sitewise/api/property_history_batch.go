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
// `query.MaxDataPoints` is ignored and it always requests with the maximum number of data points the SiteWise API can support
func historyBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueHistoryInput {
	qualities := []iotsitewisetypes.Quality{iotsitewisetypes.QualityGood}
	if query.Quality != "" && query.Quality != "ANY" {
		qualities[0] = query.Quality
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	entries := make([]iotsitewisetypes.BatchGetAssetPropertyValueHistoryEntry, 0)

	// All unique properties are collected in AssetPropertyEntries and assigned to
	// a BatchGetAssetPropertyValueHistoryEntry
	for _, entry := range query.AssetPropertyEntries {
		historyEntry := iotsitewisetypes.BatchGetAssetPropertyValueHistoryEntry{
			StartDate:    from,
			EndDate:      to,
			TimeOrdering: query.TimeOrdering,
			Qualities:    qualities,
		}
		if entry.AssetId != "" && entry.PropertyId != "" {
			historyEntry.AssetId = aws.String(entry.AssetId)
			historyEntry.PropertyId = aws.String(entry.PropertyId)
			historyEntry.EntryId = util.GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
		} else {
			// If there is no assetId or propertyId, then we use the propertyAlias
			historyEntry.PropertyAlias = aws.String(entry.PropertyAlias)
			historyEntry.EntryId = util.GetEntryIdFromPropertyAlias(entry.PropertyAlias)
		}
		entries = append(entries, historyEntry)
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryInput{
		Entries: entries,
		// performance: hardcoded to fetch the maximum number of data points
		MaxResults: aws.Int32(BatchGetAssetPropertyValueHistoryMaxResults),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValues(ctx context.Context, client client.SitewiseAPIClient,
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
