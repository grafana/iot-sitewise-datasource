package sitewise

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssets(ctx context.Context, client client.Client, query models.ListAssetsQuery) (*framer.Assets, error) {

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

	return &framer.Assets{
		AssetSummaries: resp.AssetSummaries,
		NextToken:      resp.NextToken,
	}, nil
}
