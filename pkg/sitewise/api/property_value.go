package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueInput {
	entries := make([]*iotsitewise.BatchGetAssetPropertyValueEntry, 0)

	if query.PropertyAlias != "" {
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueEntry{
			EntryId:       getAssetId(query.BaseQuery),
			PropertyAlias: getPropertyAlias(query.BaseQuery),
		})

		return &iotsitewise.BatchGetAssetPropertyValueInput{
			Entries:   entries,
			NextToken: getNextToken(query.BaseQuery),
		}
	}

	for _, assetId := range query.AssetIds {
		var id *string
		if assetId != "" {
			id = aws.String(assetId)
		}
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueEntry{
			EntryId:       id,
			AssetId:       id,
			PropertyId:    aws.String(query.PropertyId),
			PropertyAlias: getPropertyAlias(query.BaseQuery),
		})
	}

	return &iotsitewise.BatchGetAssetPropertyValueInput{
		Entries:   entries,
		NextToken: getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyValue(ctx context.Context, client client.SitewiseClient, query *models.AssetPropertyValueQuery) (*framer.AssetPropertyValue, error) {
	err := getAndSetAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return nil, err
	}

	awsReq := valueQueryToInput(*query)

	resp, err := client.BatchGetAssetPropertyValueWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyValue{
		SuccessEntries: resp.SuccessEntries,
		SkippedEntries: resp.SkippedEntries,
		ErrorEntries:   resp.ErrorEntries,
		NextToken:      resp.NextToken,
	}, nil
}
