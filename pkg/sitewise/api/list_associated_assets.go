package api

import (
	"context"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssociatedAssets(ctx context.Context, client client.SitewiseAPIClient, query models.ListAssociatedAssetsQuery) (*framer.AssociatedAssets, error) {

	var (
		hierarchyId *string
		results     []iotsitewisetypes.AssociatedAssetsSummary
	)

	seenAssetIds := make(map[string]bool)

	for _, assetId := range query.AssetIds {
		assetIdPtr := aws.String(assetId)

		// Recursively load children
		if query.LoadAllChildren {
			asset, err := client.DescribeAsset(ctx, &iotsitewise.DescribeAssetInput{
				AssetId: assetIdPtr,
			})

			if err != nil {
				return nil, err
			}

			for _, h := range asset.AssetHierarchies {
				var nextToken *string = nil

				for {
					resp, err := client.ListAssociatedAssets(ctx, &iotsitewise.ListAssociatedAssetsInput{
						AssetId:     assetIdPtr,
						HierarchyId: h.Id,
						MaxResults:  MaxSitewiseResults,
						NextToken:   nextToken,
					})

					if err != nil {
						return nil, err
					}

					for _, assetSummary := range resp.AssetSummaries {
						assetId := *assetSummary.Id
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
			traversalDirection := iotsitewisetypes.TraversalDirectionParent
			if query.HierarchyId != "" {
				hierarchyId = aws.String(query.HierarchyId)
				traversalDirection = iotsitewisetypes.TraversalDirectionChild
			}

			var nextToken *string = nil

			for {
				resp, err := client.ListAssociatedAssets(ctx, &iotsitewise.ListAssociatedAssetsInput{
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
					assetId := *assetSummary.Id
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
