package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func ListAssetProperties(ctx context.Context, client iotsitewise.ListAssetPropertiesAPIClient, query models.ListAssetPropertiesQuery) (*framer.AssetProperties, error) {
	resp, err := client.ListAssetProperties(ctx, &iotsitewise.ListAssetPropertiesInput{
		AssetId:    util.GetAssetId(query.BaseQuery),
		Filter:     iotsitewisetypes.ListAssetPropertiesFilterAll,
		MaxResults: aws.Int32(250),
		NextToken:  getNextToken(query.BaseQuery),
	})

	if err != nil {
		return nil, err
	}

	return &framer.AssetProperties{
		AssetPropertySummaries: resp.AssetPropertySummaries,
		NextToken:              resp.NextToken,
	}, nil
}
