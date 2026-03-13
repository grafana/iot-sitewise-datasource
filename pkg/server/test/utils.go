package test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/mock"
)

type batch_test struct {
	name               string
	numAssetIds        int
	numPropertyIds     int
	numPropertyAliases int
}

func Pointer[T any](v T) *T { return &v }

func generateIds(numIds int, idString string) []string {
	ids := []string{}
	for i := 1; i <= numIds; i++ {
		ids = append(ids, fmt.Sprintf("%s%d", idString, i))
	}
	return ids
}

func mockDescribeAssetProperty(mockSw *mocks.SitewiseAPIClient) {
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewisetypes.Property{
			DataType: iotsitewisetypes.PropertyDataTypeDouble,
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)
}

func mockDescribeAsset(mockSw *mocks.SitewiseAPIClient) {
	mockSw.On("DescribeAsset", mock.Anything, mock.MatchedBy(func(input *iotsitewise.DescribeAssetInput) bool {
		return input.AssetId != nil
	})).Maybe().Return(func(ctx context.Context, input *iotsitewise.DescribeAssetInput, opts ...func(*iotsitewise.Options)) *iotsitewise.DescribeAssetOutput {
		return &iotsitewise.DescribeAssetOutput{
			AssetId:      input.AssetId,
			AssetName:    Pointer("Demo Turbine Asset 1"),
			AssetModelId: Pointer("1f95cf92-34ff-4975-91a9-e9f2af35b6a5"),
			AssetProperties: []iotsitewisetypes.AssetProperty{
				{
					Id:   Pointer("3a985085-ea71-4ae6-9395-b65990f58a05"),
					Name: Pointer("RPM"),
				},
				{
					Id:   Pointer("44fa33e2-b2db-4724-ba03-48ce28902809"),
					Name: Pointer("Torque"),
				},
			},
		}
	}, nil)
}

func mockDescribeAssetModel(mockSw *mocks.SitewiseAPIClient) {
	mockSw.On("DescribeAssetModel", mock.Anything, mock.MatchedBy(func(input *iotsitewise.DescribeAssetModelInput) bool {
		return input.AssetModelId != nil
	})).Maybe().Return(func(ctx context.Context, input *iotsitewise.DescribeAssetModelInput, opts ...func(*iotsitewise.Options)) *iotsitewise.DescribeAssetModelOutput {
		return &iotsitewise.DescribeAssetModelOutput{
			AssetModelId:   input.AssetModelId,
			AssetModelName: Pointer("Demo Turbine Asset Model"),
			AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
				{
					Id:   Pointer(mockPropertyId),
					Name: Pointer("Wind Speed"),
				},
			},
		}
	}, nil)
}
