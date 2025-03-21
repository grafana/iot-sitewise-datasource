package resource

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type SitewiseResources struct {
	client client.SitewiseAPIClient
}

func NewSitewiseResources(client client.SitewiseAPIClient) *SitewiseResources {
	return &SitewiseResources{
		client: client,
	}
}

func (rp *SitewiseResources) Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error) {

	resp, err := rp.client.DescribeAsset(ctx, &iotsitewise.DescribeAssetInput{
		AssetId: aws.String(assetId),
	})

	return resp, err
}

func (rp *SitewiseResources) Property(ctx context.Context, assetId string, propertyId string, propertyAlias string) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	if propertyAlias != "" && (assetId == "" && propertyId == "") {
		defaultOutput := &iotsitewise.DescribeAssetPropertyOutput{
			AssetName: aws.String(""),
			AssetProperty: &iotsitewisetypes.Property{
				Name:     aws.String(propertyAlias),
				DataType: "?",
			},
		}

		resp, err := rp.client.DescribeTimeSeries(ctx, &iotsitewise.DescribeTimeSeriesInput{
			Alias: aws.String(propertyAlias),
		})
		if err != nil {
			return defaultOutput, err
		}

		if resp.AssetId != nil && resp.PropertyId != nil {
			return rp.client.DescribeAssetProperty(ctx, &iotsitewise.DescribeAssetPropertyInput{
				AssetId:    resp.AssetId,
				PropertyId: resp.PropertyId,
			})
		}

		return defaultOutput, nil
	}

	return rp.client.DescribeAssetProperty(ctx, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String(assetId),
		PropertyId: aws.String(propertyId),
	})
}

func (rp *SitewiseResources) AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error) {

	resp, err := rp.client.DescribeAssetModel(ctx, &iotsitewise.DescribeAssetModelInput{
		AssetModelId: aws.String(modelId),
	})

	return resp, err
}
