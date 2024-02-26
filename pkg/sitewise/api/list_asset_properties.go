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

func ListAssetProperties(ctx context.Context, client client.SitewiseClient, query models.ListAssetPropertiesQuery) (*framer.AssetProperties, error) {
	resp, err := client.ListAssetPropertiesWithContext(ctx, &iotsitewise.ListAssetPropertiesInput{
		AssetId: util.GetAssetId(query.BaseQuery),
		Filter: aws.String("ALL"),
		MaxResults: aws.Int64(250),
		NextToken: getNextToken(query.BaseQuery),
	})

	if err != nil {
		return nil, err
	}

	return &framer.AssetProperties {
		AssetPropertySummaries: resp.AssetPropertySummaries,
		NextToken: resp.NextToken,
	}, nil
}
