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

func valueBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueInput {
	entries := make([]iotsitewisetypes.BatchGetAssetPropertyValueEntry, 0)

	switch {
	case query.PropertyAlias != "":
		entries = append(entries, iotsitewisetypes.BatchGetAssetPropertyValueEntry{
			EntryId:       util.GetEntryId(query.BaseQuery),
			PropertyAlias: util.GetPropertyAlias(query.BaseQuery),
		})
	default:
		for _, assetId := range query.AssetIds {
			var id *string
			if assetId != "" {
				id = aws.String(assetId)
			}
			entries = append(entries, iotsitewisetypes.BatchGetAssetPropertyValueEntry{
				EntryId:    id,
				AssetId:    id,
				PropertyId: aws.String(query.PropertyId),
			})
		}
	}

	return &iotsitewise.BatchGetAssetPropertyValueInput{
		Entries:   entries,
		NextToken: getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValue(ctx context.Context, sw client.SitewiseAPIClient, query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValueBatch, error) {
	modifiedQuery, err := getAssetIdAndPropertyId(query, sw, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	batchedQueries := batchQueries(modifiedQuery, BatchGetAssetPropertyValueMaxEntries)
	responses := []*iotsitewise.BatchGetAssetPropertyValueOutput{}
	for _, q := range batchedQueries {
		req := valueBatchQueryToInput(q)
		resp, err := sw.BatchGetAssetPropertyValue(ctx, req)
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
		&framer.AssetPropertyValueBatch{
			Responses:       responses,
			AnomalyAssetIds: anomalyAssetIds,
			SitewiseClient:  sw,
		},
		nil
}
