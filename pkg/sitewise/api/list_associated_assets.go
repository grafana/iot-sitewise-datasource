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

func ListAssociatedAssets(ctx context.Context, client client.SitewiseClient, query models.ListAssociatedAssetsQuery) (*framer.AssociatedAssets, error) {

	var (
		hierarchyId        *string
		traversalDirection *string
		assetId            *string = util.GetAssetId(query.BaseQuery)
		results            []*iotsitewise.AssociatedAssetsSummary
	)

	// Recursively load children
	if query.LoadAllChildren {
		asset, err := client.DescribeAsset(&iotsitewise.DescribeAssetInput{
			AssetId: assetId,
		})

		if err != nil {
			return nil, err
		}

		for _, h := range asset.AssetHierarchies {
			// For this code path, we need to handle this internally, since it's a union of multiple queries
			var nextToken *string = nil

			for {
				resp, err := client.ListAssociatedAssetsWithContext(ctx, &iotsitewise.ListAssociatedAssetsInput{
					AssetId:     assetId,
					HierarchyId: h.Id,
					MaxResults:  MaxSitewiseResults,
					NextToken:   nextToken,
				})

				if err != nil {
					return nil, err
				}

				results = append(results, resp.AssetSummaries...)

				if resp.NextToken == nil {
					break
				}

				nextToken = resp.NextToken
			}
		}

		return &framer.AssociatedAssets{
			AssetSummaries: results,
		}, nil

	} else {
		if query.HierarchyId != "" {
			hierarchyId = aws.String(query.HierarchyId)
			traversalDirection = aws.String("CHILD")
		} else {
			traversalDirection = aws.String("PARENT")
		}

		resp, err := client.ListAssociatedAssetsWithContext(ctx, &iotsitewise.ListAssociatedAssetsInput{
			AssetId:            assetId,
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
}
