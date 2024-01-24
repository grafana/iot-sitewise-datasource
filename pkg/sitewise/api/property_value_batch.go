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

	awsReq := valueBatchQueryToInput(modifiedQuery)

	resp, err := client.BatchGetAssetPropertyValueWithContext(ctx, awsReq)

	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	return modifiedQuery,
		&framer.AssetPropertyValueBatch{
			SuccessEntries: resp.SuccessEntries,
			SkippedEntries: resp.SkippedEntries,
			ErrorEntries:   resp.ErrorEntries,
			NextToken:      resp.NextToken,
		},
		nil
}
