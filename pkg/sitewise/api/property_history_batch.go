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
func historyBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueHistoryInput {

	var (
		qualities []*string
	)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 20000 {
		query.MaxDataPoints = 20000
	}

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
		Entries:    entries,
		MaxResults: aws.Int64(query.MaxDataPoints),
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

	awsReq := historyBatchQueryToInput(modifiedQuery)
	resp, err := client.BatchGetAssetPropertyValueHistoryPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
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
			BatchGetAssetPropertyValueHistoryOutput: resp,
			Query:                                   modifiedQuery,
			AnomalyAssetIds:                         anomalyAssetIds,
			SitewiseClient:                          client,
		},
		nil
}
