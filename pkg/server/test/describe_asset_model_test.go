package test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

func TestHandleDescribeAssetModel(t *testing.T) {
	describeAssetModelHappyCase(t).run(t)
}

var describeAssetModelHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseAPIClient{}

	assetModel := testdata.GetIoTSitewiseAssetModelDescription(t, testDataRelativePath("describe-asset-model.json"))

	mockSw.On("DescribeAssetModel", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetModelInput) bool {
		return req.AssetModelId != nil && *req.AssetModelId == testdata.DemoTurbineAssetModelId
	})).Return(&assetModel, nil)

	query := models.DescribeAssetModelQuery{}
	query.AssetModelId = testdata.DemoTurbineAssetModelId

	return &testScenario{
		name: "DescribeAssetModelHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeDescribeAssetModel,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "describe-asset-model",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleDescribeAssetModel
		},
		validationFn: nil,
	}
}
