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

	// All unique properties are collected in AssetPropertyEntries and assigned to
	// a BatchGetAssetPropertyValueEntry
	for _, entry := range query.AssetPropertyEntries {
		batchEntry := iotsitewise.BatchGetAssetPropertyValueEntry{}
		if entry.AssetId != "" && entry.PropertyId != "" {
			batchEntry.AssetId = aws.String(entry.AssetId)
			batchEntry.PropertyId = aws.String(entry.PropertyId)
			batchEntry.EntryId = util.GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
		} else {
			// If there is no assetId or propertyId, then we use the propertyAlias
			batchEntry.PropertyAlias = aws.String(entry.PropertyAlias)
			batchEntry.EntryId = util.GetEntryIdFromPropertyAlias(entry.PropertyAlias)
		}
		entries = append(entries, &batchEntry)
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
