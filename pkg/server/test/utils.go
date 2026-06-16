package test

import (
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
	mockSw.On("DescribeAsset", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetOutput{
		AssetId:      Pointer(mockAssetId),
		AssetName:    Pointer("Demo Turbine Asset 1"),
		AssetModelId: Pointer("1f95cf92-34ff-4975-91a9-e9f2af35b6a5"),
	}, nil)
}

func mockDescribeAssetModel(mockSw *mocks.SitewiseAPIClient) {
	mockSw.On("DescribeAssetModel", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetModelOutput{
		AssetModelId:   Pointer("1f95cf92-34ff-4975-91a9-e9f2af35b6a5"),
		AssetModelName: Pointer("Demo Turbine Asset Model"),
		AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
			{
				Id:   Pointer(mockPropertyId),
				Name: Pointer("Wind Speed"),
			},
		},
	}, nil)
}
