package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"

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

func Test_getPropertyValueBooleanFromAlias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValueHistory(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
				}),
			},
		},
	})
	require.Nil(t, err)

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-boolean", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}
