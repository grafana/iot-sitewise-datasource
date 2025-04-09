package test

import (
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

func TestHandleListAssets(t *testing.T) {
	listAssetsHappyCase(t).run(t)
}

var listAssetsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseAPIClient{}

	topLevelAssets := testdata.GetIoTSitewiseAssets(t, testDataRelativePath("list-assets-top-level.json"))
	childAssets := testdata.GetIoTSitewiseAssets(t, testDataRelativePath("list-assets.json"))

	mockSw.On("ListAssets", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		return req.AssetModelId == nil && req.Filter == types.ListAssetsFilterTopLevel
	})).Return(&topLevelAssets, nil)

	mockSw.On("ListAssets", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		if req.AssetModelId == nil {
			return false
		}
		return *req.AssetModelId == testdata.DemoTurbineAssetModelId && req.Filter == types.ListAssetsFilterAll
	})).Return(&childAssets, nil)

	queryTopLevel := models.ListAssetsQuery{
		ModelId: "",
		Filter:  "",
	}

	queryChild := models.ListAssetsQuery{
		ModelId: testdata.DemoTurbineAssetModelId,
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
