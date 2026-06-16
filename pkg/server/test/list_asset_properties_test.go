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

func TestHandleListAssetProperties(t *testing.T) {
	listAssetPropertiesHappyCase(t).run(t)
}

var listAssetPropertiesHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {
	mockSw := &mocks.SitewiseAPIClient{}

	assetProperties := testdata.GetIoTSitewiseAssetProperties(t, testDataRelativePath("list-asset-properties.json"))

	mockSw.On("ListAssetProperties", mock.Anything, mock.Anything).Return(&assetProperties, nil)

	query := models.ListAssetPropertiesQuery{
		BaseQuery: models.BaseQuery{AssetId: "123"},
	}

	return &testScenario{
		name:   "TestListAssetPropertiesResponseHappyCase",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssetProperties,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		goldenFileName: "list-asset-properties",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssetProperties
		},
	}
}
