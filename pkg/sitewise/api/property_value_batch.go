package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func valueBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueInput {
	entries := make([]*iotsitewise.BatchGetAssetPropertyValueEntry, 0)

	switch {
	case query.PropertyAlias != "":
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueEntry{
			EntryId:       util.GetEntryId(query.BaseQuery),
			PropertyAlias: util.GetPropertyAlias(query.BaseQuery),
		})
	default:
		for _, assetId := range query.AssetIds {
			var id *string
			if assetId != "" {
				id = aws.String(assetId)
			}
			entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueEntry{
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

func BatchGetAssetPropertyValue(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValueBatch, error) {
	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	batchedQueries := batchQueries(modifiedQuery, BatchGetAssetPropertyValueMaxEntries)
	responses := []*iotsitewise.BatchGetAssetPropertyValueOutput{}
	for _, q := range batchedQueries {
		req := valueBatchQueryToInput(q)
		resp, err := client.BatchGetAssetPropertyValueWithContext(ctx, req)
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
		&framer.AssetPropertyValueBatch{
			Responses:       responses,
			AnomalyAssetIds: anomalyAssetIds,
			SitewiseClient:  client,
		},
		nil
}
