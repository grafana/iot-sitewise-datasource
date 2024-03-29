package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func GetAssetPropertyDescription(ctx context.Context, client client.SitewiseClient, query models.DescribeAssetPropertyQuery) (*framer.AssetProperty, error) {

	awsReq := &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    util.GetAssetId(query.BaseQuery),
		PropertyId: util.GetPropertyId(query.BaseQuery),
	}

	resp, err := client.DescribeAssetPropertyWithContext(ctx, awsReq)
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
