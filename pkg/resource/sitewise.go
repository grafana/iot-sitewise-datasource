package resource

import (
	"context"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type SitewiseResources struct {
	sw client.SitewiseAPIClient
}

func NewSitewiseResources(sw client.SitewiseAPIClient) *SitewiseResources {
	return &SitewiseResources{
		sw: sw,
	}
}

func (rp *SitewiseResources) Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error) {

	resp, err := rp.sw.DescribeAsset(ctx, &iotsitewise.DescribeAssetInput{
		AssetId: aws.String(assetId),
	})

	return resp, err
}

func (rp *SitewiseResources) Property(ctx context.Context, assetId string, propertyId string, propertyAlias string) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	if propertyAlias != "" && (assetId == "" && propertyId == "") {
		return &iotsitewise.DescribeAssetPropertyOutput{
			AssetName: aws.String(""),
			AssetProperty: &iotsitewisetypes.Property{
				Name:     aws.String(propertyAlias),
				DataType: "?",
			},
		}, nil
	}

	return rp.sw.DescribeAssetProperty(ctx, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String(assetId),
		PropertyId: aws.String(propertyId),
	})
}

func (rp *SitewiseResources) AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error) {

	resp, err := rp.sw.DescribeAssetModel(ctx, &iotsitewise.DescribeAssetModelInput{
		AssetModelId: aws.String(modelId),
	})

	return resp, err
}
