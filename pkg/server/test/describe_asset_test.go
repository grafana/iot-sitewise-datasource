package test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

func TestHandleDescribeAsset(t *testing.T) {
	describeAssetHappyCase(t).run(t)
}

var describeAssetHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseAPIClient{}

	asset := testdata.GetIoTSitewiseAssetDescription(t, testDataRelativePath("describe-asset.json"))
	assetWithHierarchy := testdata.GetIoTSitewiseAssetDescription(t, testDataRelativePath("describe-asset-top-level.json"))

	mockSw.On("DescribeAsset", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetInput) bool {
		return req.AssetId != nil && *req.AssetId == testdata.DemoTurbineAsset1
	})).Return(&asset, nil)

	mockSw.On("DescribeAsset", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetInput) bool {
		return req.AssetId != nil && *req.AssetId == testdata.DemoWindFarmAssetId
	})).Return(&assetWithHierarchy, nil)

	query := models.DescribeAssetQuery{}
	query.AssetIds = []string{testdata.DemoTurbineAsset1}

	queryTopLevel := models.DescribeAssetQuery{}
	queryTopLevel.AssetIds = []string{testdata.DemoWindFarmAssetId}

	return &testScenario{
		name: "DescribeAssetHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeDescribeAsset,
				JSON:      testdata.SerializeStruct(t, query),
			},
			{
				RefID:     "B",
				QueryType: models.QueryTypeDescribeAsset,
				JSON:      testdata.SerializeStruct(t, queryTopLevel),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "describe-asset",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleDescribeAsset
		},
		validationFn: nil,
	}
}
