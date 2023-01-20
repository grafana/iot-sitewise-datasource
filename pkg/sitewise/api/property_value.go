package api

import (
	"context"
	"fmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func valueQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyValueInput {

	if query.AssetId != "" {
		query.AssetIds = []string{query.AssetId}
	}

	entries := make([]*iotsitewise.BatchGetAssetPropertyValueEntry, 0)
	for i, assetId := range query.AssetIds {
		var id *string
		if assetId != "" {
			id = aws.String(assetId)
		}
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueEntry{
			EntryId:       aws.String(fmt.Sprintf("%d", i)),
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

func BatchGetAssetPropertyValue(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyValue, error) {

	awsReq := valueQueryToInput(query)

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
