package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func DescribeAssetModel(ctx context.Context, client client.SitewiseClient, query models.DescribeAssetModelQuery) (*framer.AssetModelDescription, error) {

	awsReq := &iotsitewise.DescribeAssetModelInput{AssetModelId: aws.String(query.AssetModelId)}

	resp, err := client.DescribeAssetModelWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.AssetModelDescription{
		AssetModelArn:             resp.AssetModelArn,
		AssetModelCompositeModels: resp.AssetModelCompositeModels,
		AssetModelCreationDate:    resp.AssetModelCreationDate,
		AssetModelDescription:     resp.AssetModelDescription,
		AssetModelHierarchies:     resp.AssetModelHierarchies,
		AssetModelId:              resp.AssetModelId,
		AssetModelLastUpdateDate:  resp.AssetModelLastUpdateDate,
		AssetModelName:            resp.AssetModelName,
		AssetModelProperties:      resp.AssetModelProperties,
		AssetModelStatus:          resp.AssetModelStatus,
	}, nil
}
