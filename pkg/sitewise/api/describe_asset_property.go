package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func GetAssetPropertyDescription(ctx context.Context, sw client.SitewiseAPIClient, query models.DescribeAssetPropertyQuery) (*framer.AssetProperty, error) {

	awsReq := &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    util.GetAssetId(query.BaseQuery),
		PropertyId: util.GetPropertyId(query.BaseQuery),
	}

	resp, err := sw.DescribeAssetProperty(ctx, awsReq)
	if err != nil {
		return nil, err
	}

	return &framer.AssetProperty{
		AssetId:       resp.AssetId,
		AssetModelId:  resp.AssetModelId,
		AssetName:     resp.AssetName,
		AssetProperty: resp.AssetProperty,
	}, nil
}
