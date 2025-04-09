package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssetModels(ctx context.Context, sw client.SitewiseAPIClient, query models.ListAssetModelsQuery) (*framer.AssetModels, error) {

	resp, err := sw.ListAssetModels(ctx, &iotsitewise.ListAssetModelsInput{
		MaxResults: aws.Int32(250),
		NextToken:  getNextToken(query.BaseQuery),
	})

	if err != nil {
		return nil, err
	}

	return &framer.AssetModels{
		AssetModelSummaries: resp.AssetModelSummaries,
		NextToken:           resp.NextToken,
	}, nil
}
