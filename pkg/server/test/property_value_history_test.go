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
	getPropertyValueHistoryHappyCaseTable(t).run(t)
	getPropertyValueHistoryHappyCaseTimeSeries(t).run(t)
	getPropertyValueBoolean(t).run(t)
	getPropertyValueHistoryFromAliasCaseTable(t).run(t)
	getPropertyValueHistoryFromAliasCaseTimeSeries(t).run(t)
	getPropertyValueBooleanFromAlias(t).run(t)
}

var getPropertyValueHistoryHappyCaseTable testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  testdata.AwsRegion,
			AssetId:    testdata.DemoTurbineAsset1,
			PropertyId: testdata.TurbinePropAvgWindSpeed,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryResponseHappyCaseTable",
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
		goldenFileName: "property-history-values-table",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var getPropertyValueHistoryHappyCaseTimeSeries testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			ResponseFormat: "timeseries",
			AwsRegion:      testdata.AwsRegion,
			AssetId:        testdata.DemoTurbineAsset1,
			PropertyId:     testdata.TurbinePropAvgWindSpeed,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryResponseHappyCaseTimeSeries",
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
		goldenFileName: "property-history-values-timeseries",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var getPropertyValueBoolean testServerScenarioFn = func(t *testing.T) *testScenario {
	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
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

var getPropertyValueHistoryFromAliasCaseTable testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:     testdata.AwsRegion,
			PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryFromAliasResponseHappyCaseTable",
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
		goldenFileName: "property-history-values-from-alias-table",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var getPropertyValueHistoryFromAliasCaseTimeSeries testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			ResponseFormat: "timeseries",
			AwsRegion:      testdata.AwsRegion,
			PropertyAlias:  testdata.TurbinePropWindSpeedAlias,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryFromAliasResponseHappyCaseTimeSeries",
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
		goldenFileName: "property-history-values-from-alias-timeseries",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var getPropertyValueBooleanFromAlias testServerScenarioFn = func(t *testing.T) *testScenario {
	mockSw := &mocks.SitewiseClient{}

	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))

	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:     testdata.AwsRegion,
			PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryFromAliasResponseBoolean",
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
		goldenFileName: "property-history-values-from-alias-boolean",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}
