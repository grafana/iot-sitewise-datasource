package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type SitewiseResources struct {
	client client.SitewiseClient
}

func NewSitewiseResources(client client.SitewiseClient) *SitewiseResources {
	return &SitewiseResources{
		client: client,
	}
}

func (rp *SitewiseResources) Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error) {

	resp, err := rp.client.DescribeAssetWithContext(ctx, &iotsitewise.DescribeAssetInput{
		AssetId: aws.String(assetId),
	})

	return resp, err
}

func (rp *SitewiseResources) Property(ctx context.Context, assetId string, propertyId string, propertyAlias string) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	if propertyAlias != "" && (assetId == "" && propertyId == "") {
		return &iotsitewise.DescribeAssetPropertyOutput{
			AssetName: aws.String(""),
			AssetProperty: &iotsitewise.Property{
				Name:     aws.String(propertyAlias),
				DataType: aws.String("?"),
			},
		}, nil
	}

	return rp.client.DescribeAssetPropertyWithContext(ctx, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String(assetId),
		PropertyId: aws.String(propertyId),
	})
}

func (rp *SitewiseResources) AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error) {

	resp, err := rp.client.DescribeAssetModelWithContext(ctx, &iotsitewise.DescribeAssetModelInput{
		AssetModelId: aws.String(modelId),
	})

	return resp, err
}
