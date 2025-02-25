package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssets(ctx context.Context, client client.SitewiseClient, query models.ListAssetsQuery, taggingApiClient client.TaggingApiClient, includedTagPatterns []map[string][]string) (*framer.Assets, error) {
	// func ListAssets(ctx context.Context, client client.SitewiseClient, query models.ListAssetsQuery, taggingApiClient client.TaggingApiClient) (*framer.Assets, error) {

	var (
		filter       *string
		assetModelId *string
	)

	if query.Filter != "" {
		filter = aws.String(query.Filter)
	}

	if query.ModelId != "" {
		assetModelId = aws.String(query.ModelId)
	}

	if assetModelId == nil && filter == nil {
		// only top level filters can be used without a model id
		filter = aws.String("TOP_LEVEL")
	}

	resp, err := client.ListAssetsWithContext(ctx, &iotsitewise.ListAssetsInput{
		AssetModelId: assetModelId,
		Filter:       filter,
		MaxResults:   aws.Int64(250),
		NextToken:    getNextToken(query.BaseQuery),
	})
	if err != nil {
		return nil, err
	}

	assetSummaries := resp.AssetSummaries

	// get assets arns
	assetArns := make([]*string, 0, len(assetSummaries))
	for _, asset := range assetSummaries {
		assetArns = append(assetArns, asset.Arn)
	}

	// get resources with taggingApiClient
	resources, err := taggingApiClient.GetResourcesPage(ctx, assetArns)
	if err != nil {
		return nil, err
	}

	allowArns := FilterResourcesByTags(resources, includedTagPatterns)
	allowAssets := FilterAssetSummariesByArns(assetSummaries, allowArns)

	return &framer.Assets{
		AssetSummaries: allowAssets,
		NextToken:      resp.NextToken,
	}, nil
}
