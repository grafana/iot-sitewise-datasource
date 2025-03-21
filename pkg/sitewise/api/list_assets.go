package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssets(ctx context.Context, sw client.SitewiseAPIClient, query models.ListAssetsQuery) (*framer.Assets, error) {

	var (
		filter       iotsitewisetypes.ListAssetsFilter
		assetModelId *string
	)

	if query.Filter != "" {
		filter = query.Filter
	}

	if query.ModelId != "" {
		assetModelId = aws.String(query.ModelId)
	}

	if assetModelId == nil && filter == "" {
		// only top level filters can be used without a model id
		filter = iotsitewisetypes.ListAssetsFilterTopLevel
	}

	resp, err := sw.ListAssets(ctx, &iotsitewise.ListAssetsInput{
		AssetModelId: assetModelId,
		Filter:       filter,
		MaxResults:   aws.Int32(250),
		NextToken:    getNextToken(query.BaseQuery),
	})

	if err != nil {
		return nil, err
	}

	return &framer.Assets{
		AssetSummaries: resp.AssetSummaries,
		NextToken:      resp.NextToken,
	}, nil
}
