package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_propertyValueForTimeRange_raw_data_for_time_range(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionRaw, actual.Resolution)
	if propvals.ResolutionRaw == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-raw-data-for-time-range", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1m_data_for_time_range(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.TwoHours, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionMinute, actual.Resolution)
	if propvals.ResolutionMinute == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1m-data-for-time-range", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1h_data_for_time_range(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.OneDay, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionHour, actual.Resolution)
	if propvals.ResolutionHour == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1h-data-for-time-range", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1d_data_for_time_range(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.OneMonth, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionDay, actual.Resolution)
	if propvals.ResolutionDay == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1d-data-for-time-range", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1m_data_for_reduced_max_data_point(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
				MaxDataPoints: 299,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionMinute, actual.Resolution)
	if propvals.ResolutionMinute == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1m-data-for-reduced-max-data-point", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_raw_data_for_time_range_from_alias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionRaw, actual.Resolution)
	if propvals.ResolutionRaw == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-raw-data-for-time-range-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1m_data_for_time_range_from_alias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.TwoHours, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionMinute, actual.Resolution)
	if propvals.ResolutionMinute == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1m-data-for-time-range-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1h_data_for_time_range_from_alias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.OneDay, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionHour, actual.Resolution)
	if propvals.ResolutionHour == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1h-data-for-time-range-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1d_data_for_time_range_from_alias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.OneMonth, To: testdata.Now},
				MaxDataPoints: 720,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionDay, actual.Resolution)
	if propvals.ResolutionDay == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1d-data-for-time-range-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueForTimeRange_1m_data_for_reduced_max_data_point_from_alias(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	srvr := &server.Server{
		Datasource: mockedDatasource(mockSw).(*sitewise.Datasource),
	}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypePropertyAggregate,
				TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
				MaxDataPoints: 299,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: testdata.TurbinePropWindSpeedAlias,
					},
					AggregateTypes: []string{"avg"},
					Resolution:     "AUTO",
				}),
			},
		},
	})
	require.Nil(t, err)

	resp := qdr.Responses["A"]
	frame := resp.Frames[0]
	actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	assert.True(t, ok, "unable to cast custom metadata")
	assert.Equal(t, propvals.ResolutionMinute, actual.Resolution)
	if propvals.ResolutionMinute == "RAW" {
		assert.Equal(t, "raw", frame.Fields[1].Name)
	}

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "prop-val-for-time-range-1m-data-for-reduced-max-data-point-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}
