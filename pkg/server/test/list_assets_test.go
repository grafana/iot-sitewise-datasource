package test

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/mock"
)

func TestHandleListAssets(t *testing.T) {
	listAssetsHappyCase(t).run(t)
}

var listAssetsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	topLevelAssets := testdata.GetIoTSitewiseAssets(t, testDataRelativePath("list-assets-top-level.json"))
	childAssets := testdata.GetIoTSitewiseAssets(t, testDataRelativePath("list-assets.json"))

	mockSw.On("ListAssetsWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		return req.AssetModelId == nil && *req.Filter == "TOP_LEVEL"
	})).Return(&topLevelAssets, nil)

	mockSw.On("ListAssetsWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		if req.AssetModelId == nil {
			return false
		}
		return *req.AssetModelId == testdata.TestAssetModelId && *req.Filter == "ALL"
	})).Return(&childAssets, nil)

	queryTopLevel := models.ListAssetsQuery{
		ModelId: "",
		Filter:  "",
	}

	queryChild := models.ListAssetsQuery{
		ModelId: testdata.TestAssetModelId,
		Filter:  "ALL",
	}

	return &testScenario{
		name: "TestListAssetsHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssets,
				JSON:      testdata.SerializeStruct(t, queryTopLevel),
			},
			{
				RefID:     "B",
				QueryType: models.QueryTypeListAssets,
				JSON:      testdata.SerializeStruct(t, queryChild),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "list-assets",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssets
		},
		validationFn: nil,
	}
}
