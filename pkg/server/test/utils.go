package test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
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

func mockDescribeAssetPropertyWithContext(mockSw *mocks.SitewiseClient) {
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)
}
