package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssociatedAssets(ctx context.Context, client client.SitewiseClient, query models.ListAssociatedAssetsQuery) (*framer.AssociatedAssets, error) {

	var (
		hierarchyId        *string
		traversalDirection *string
	)

	if query.HierarchyId != "" {
		hierarchyId = aws.String(query.HierarchyId)
		traversalDirection = aws.String("CHILD")
	} else {
		traversalDirection = aws.String("PARENT")
	}

	if query.PropertyAlias != "" {
		resp, err := client.DescribeTimeSeriesWithContext(ctx, &iotsitewise.DescribeTimeSeriesInput{
			Alias: aws.String(query.PropertyAlias),
		})
		if err != nil {
			return nil, err
		}
		assetsIds := []string{*resp.AssetId}
		query.AssetIds = assetsIds
		query.AssetId = *resp.AssetId
	}

	resp, err := client.ListAssociatedAssetsWithContext(ctx, &iotsitewise.ListAssociatedAssetsInput{
		AssetId:            getAssetId(query.BaseQuery),
		HierarchyId:        hierarchyId,
		MaxResults:         MaxSitewiseResults,
		NextToken:          getNextToken(query.BaseQuery),
		TraversalDirection: traversalDirection,
	})

	if err != nil {
		return nil, err
	}

	return &framer.AssociatedAssets{
		AssetSummaries: resp.AssetSummaries,
		NextToken:      resp.NextToken,
	}, nil
}
