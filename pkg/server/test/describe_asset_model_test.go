package test

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func TestHandleDescribeAssetModel(t *testing.T) {
	describeAssetModelHappyCase(t).run(t)
}

var describeAssetModelHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	assetModel := testdata.GetIoTSitewiseAssetModelDescription(t, testDataRelativePath("describe-asset-model.json"))

	mockSw.On("DescribeAssetModelWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetModelInput) bool {
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
