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
		results            []*iotsitewise.AssociatedAssetsSummary
	)

	seenAssetIds := make(map[string]bool)

	for _, assetId := range query.BaseQuery.AssetIds {
		assetIdPtr := aws.String(assetId)

		// Recursively load children
		if query.LoadAllChildren {
			asset, err := client.DescribeAsset(&iotsitewise.DescribeAssetInput{
				AssetId: assetIdPtr,
			})

			if err != nil {
				return nil, err
			}

			for _, h := range asset.AssetHierarchies {
				var nextToken *string = nil

				for {
					resp, err := client.ListAssociatedAssetsWithContext(ctx, &iotsitewise.ListAssociatedAssetsInput{
						AssetId:     assetIdPtr,
						HierarchyId: h.Id,
						MaxResults:  MaxSitewiseResults,
						NextToken:   nextToken,
					})

					if err != nil {
						return nil, err
					}

					for _, assetSummary := range resp.AssetSummaries {
						assetId := aws.StringValue(assetSummary.Id)
						if !seenAssetIds[assetId] {
							results = append(results, assetSummary)
							seenAssetIds[assetId] = true
						}
					}

					if resp.NextToken == nil {
						break
					}

					nextToken = resp.NextToken
				}
			}
		} else {
			if query.HierarchyId != "" {
				hierarchyId = aws.String(query.HierarchyId)
				traversalDirection = aws.String("CHILD")
			} else {
				traversalDirection = aws.String("PARENT")
			}

			var nextToken *string = nil

			for {
				resp, err := client.ListAssociatedAssetsWithContext(ctx, &iotsitewise.ListAssociatedAssetsInput{
					AssetId:            assetIdPtr,
					HierarchyId:        hierarchyId,
					MaxResults:         MaxSitewiseResults,
					NextToken:          nextToken,
					TraversalDirection: traversalDirection,
				})

				if err != nil {
					return nil, err
				}

				for _, assetSummary := range resp.AssetSummaries {
					assetId := aws.StringValue(assetSummary.Id)
					if !seenAssetIds[assetId] {
						results = append(results, assetSummary)
						seenAssetIds[assetId] = true
					}
				}

				if resp.NextToken == nil {
					break
				}

				nextToken = resp.NextToken
			}
		}
	}

	return &framer.AssociatedAssets{
		AssetSummaries: results,
	}, nil
}
