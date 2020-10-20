package sitewise

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ListAssetModels(ctx context.Context, client client.Client, query models.ListAssetModelsQuery) (*framer.AssetModels, error) {

	var token *string
	if query.NextToken != "" {
		token = aws.String(query.NextToken)
	}

	resp, err := client.ListAssetModelsWithContext(ctx, &iotsitewise.ListAssetModelsInput{
		MaxResults: aws.Int64(250),
		NextToken:  token,
	})

	if err != nil {
		return nil, err
	}

	return &framer.AssetModels{
		AssetModelSummaries: resp.AssetModelSummaries,
		NextToken:           resp.NextToken,
	}, nil
}
