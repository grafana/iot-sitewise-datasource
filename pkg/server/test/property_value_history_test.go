package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_get_property_value_history_with_default_aka_table_response_format(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry{
			{
				AssetPropertyValueHistory: []*iotsitewise.AssetPropertyValue{
					{
						Quality: Pointer("GOOD"),
						Timestamp: &iotsitewise.TimeInNanos{
							OffsetInNanos: Pointer(int64(0)),
							TimeInSeconds: Pointer(int64(1612207200)),
						},
						Value: &iotsitewise.Variant{
							DoubleValue: Pointer(float64(23.8)),
						},
					},
				},
				EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

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
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
						 "propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("Wind Speed", nil, []float64{23.8}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		data.NewField("quality", nil, []string{"GOOD"}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{Resolution: "RAW"},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		1,
		100,
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
		},
	)
}

func Test_get_property_value_history_with_time_series_response_format(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry{
			{
				AssetPropertyValueHistory: []*iotsitewise.AssetPropertyValue{
					{
						Quality: Pointer("GOOD"),
						Timestamp: &iotsitewise.TimeInNanos{
							OffsetInNanos: Pointer(int64(0)),
							TimeInSeconds: Pointer(int64(1612207200)),
						},
						Value: &iotsitewise.Variant{
							DoubleValue: Pointer(float64(23.8)),
						},
					},
				},
				EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

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
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
						 "propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
						 "responseFormat":"timeseries"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("Wind Speed", data.Labels{"quality": "GOOD"}, []*float64{Pointer(23.8)}),
	).SetMeta(&data.FrameMeta{
		Type:   data.FrameTypeTimeSeriesWide,
		Custom: models.SitewiseCustomMeta{Resolution: "RAW"},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		1,
		100,
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
		},
	)
}

func Test_getPropertyValueBoolean(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

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
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropAvgWindSpeed,
					},
				}),
			},
		},
	})
	require.Nil(t, err)

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-boolean", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_getPropertyValueHistoryFromAliasCaseTable(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-table", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_getPropertyValueHistoryFromAliasCaseTimeSeries(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
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
						ResponseFormat: "timeseries",
						AwsRegion:      testdata.AwsRegion,
						PropertyAlias:  testdata.TurbinePropWindSpeedAlias,
					},
				}),
			},
		},
	})
	require.Nil(t, err)

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-timeseries", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
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

func Test_get_property_value_history_from_expression_query_with_time_series_response_format(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry{
			{
				AssetPropertyValueHistory: []*iotsitewise.AssetPropertyValue{
					{
						Quality: Pointer("GOOD"),
						Timestamp: &iotsitewise.TimeInNanos{
							OffsetInNanos: Pointer(int64(0)),
							TimeInSeconds: Pointer(int64(1612207200)),
						},
						Value: &iotsitewise.Variant{
							DoubleValue: Pointer(float64(23.8)),
						},
					},
				},
				EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValueHistory(context.Background(), &backend.QueryDataRequest{
		Headers:       map[string]string{"http_X-Grafana-From-Expr": "true"},
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
						 "propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
						 "responseFormat":"timeseries"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("Wind Speed", data.Labels{"quality": "GOOD"}, []*float64{Pointer(23.8)}),
	).SetMeta(&data.FrameMeta{
		Type:   data.FrameTypeTimeSeriesWide,
		Custom: models.SitewiseCustomMeta{Resolution: "RAW"},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		int(math.Inf(1)),
		int(math.Inf(1)),
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
		},
	)
}
