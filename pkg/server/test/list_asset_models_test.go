package test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/mock"
)

func TestHandleListAssetModels(t *testing.T) {
	listAssetModelsHappyCase(t).run(t)
}

var listAssetModelsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseAPIClient{}

	assetModels := testdata.GetIoTSitewiseAssetModels(t, testDataRelativePath("list-asset-models.json"))

	mockSw.On("ListAssetModels", mock.Anything, mock.Anything).Return(&assetModels, nil)

	query := models.ListAssetModelsQuery{
		BaseQuery: models.BaseQuery{},
	}

	return &testScenario{
		name:   "TestListAssetModelsResponseHappyCase",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypeListAssetModels,
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     backend.TimeRange{},
				JSON:          testdata.SerializeStruct(t, query),
			},
		},
		goldenFileName: "list-asset-models",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssetModels
		},
	}
}
