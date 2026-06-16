package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func DescribeAsset(ctx context.Context, sw client.SitewiseAPIClient, query models.DescribeAssetQuery) (*framer.AssetDescription, error) {
	awsReq := &iotsitewise.DescribeAssetInput{AssetId: util.GetAssetId(query.BaseQuery)}

	resp, err := sw.DescribeAsset(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.AssetDescription{
		AssetArn:            resp.AssetArn,
		AssetCreationDate:   resp.AssetCreationDate,
		AssetHierarchies:    resp.AssetHierarchies,
		AssetId:             resp.AssetId,
		AssetLastUpdateDate: resp.AssetLastUpdateDate,
		AssetModelId:        resp.AssetModelId,
		AssetName:           resp.AssetName,
		AssetProperties:     resp.AssetProperties,
		AssetStatus:         resp.AssetStatus,
	}, nil
}
