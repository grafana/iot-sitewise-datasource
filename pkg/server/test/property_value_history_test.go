package test

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/mock"
)

func TestHandlePropertyValueHistory(t *testing.T) {
	getPropertyValueHistoryHappyCase(t).run(t)
	getPropertyValueBoolean(t).run(t)
}

var getPropertyValueHistoryHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))

	mockSw.On("GetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  testdata.AwsRegion,
			AssetId:    testdata.DemoTurbineAsset1,
			PropertyId: testdata.TurbinePropAvgWindSpeed,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryResponseHappyCase",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON:          testdata.SerializeStruct(t, query),
			},
		},
		goldenFileName: "property-history-values",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var getPropertyValueBoolean testServerScenarioFn = func(t *testing.T) *testScenario {
	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))

	mockSw.On("GetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  testdata.AwsRegion,
			AssetId:    testdata.DemoTurbineAsset1,
			PropertyId: testdata.TurbinePropAvgWindSpeed,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryResponseBoolean",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON:          testdata.SerializeStruct(t, query),
			},
		},
		goldenFileName: "property-history-values-boolean",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}
