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

	expectedFrame := data.NewFrame("Demo Turbine Asset 1",
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

	expectedFrame := data.NewFrame("Demo Turbine Asset 1",
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

func Test_get_property_value_history_with_flatten_l4e(t *testing.T) {
	assetId := "1assetid-aaaa-2222-bbbb-3333cccc4444"
	assetPropertyIdQuery := "11propid-aaaa-2222-bbbb-3333cccc4444"
	assetPropertyIdDiagnosticOne := "44fa33e2-b2db-4724-ba03-48ce28902809"
	assetPropertyIdDiagnosticTwo := "3a985085-ea71-4ae6-9395-b65990f58a05"

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
							StringValue: Pointer("{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"),
						},
					},
				},
				EntryId: Pointer(assetId),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == assetPropertyIdQuery
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(assetId),
		AssetName: Pointer("Demo Turbine Asset 1"),
		CompositeModel: &iotsitewise.CompositeModelProperty{
			Name: Pointer("prediction1"),
			AssetProperty: &iotsitewise.Property{
				Name:     Pointer("AWS/L4E_ANOMALY_RESULT"),
				DataType: Pointer("STRUCT"),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == assetPropertyIdDiagnosticOne
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			Id:       Pointer(assetPropertyIdDiagnosticOne),
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Torque"),
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == assetPropertyIdDiagnosticTwo
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			Id:       Pointer(assetPropertyIdDiagnosticTwo),
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("RPM"),
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
					   "flattenL4E": true
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("quality", nil, []string{"GOOD"}),
		data.NewField("anomaly_score", nil, []float64{0.2674}),
		data.NewField("prediction_reason", nil, []string{"NO_ANOMALY_DETECTED"}),
		data.NewField("RPM", nil, []float64{0.44856}),
		data.NewField("Torque", nil, []float64{0.55144}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{Resolution: "RAW"},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func Test_get_property_value_history_with_struct_type(t *testing.T) {
	structValue := "{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"
	assetId := "1assetid-aaaa-2222-bbbb-3333cccc4444"

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
							StringValue: Pointer(structValue),
						},
					},
				},
				EntryId: Pointer(assetId),
			},
		},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(assetId),
		AssetName: Pointer("Demo Turbine Asset 1"),
		CompositeModel: &iotsitewise.CompositeModelProperty{
			Name: Pointer("prediction1"),
			AssetProperty: &iotsitewise.Property{
				Name:     Pointer("AWS/L4E_ANOMALY_RESULT"),
				DataType: Pointer("STRUCT"),
			},
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

	expectedFrame := data.NewFrame("Demo Turbine Asset 1",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("AWS/L4E_ANOMALY_RESULT", nil, []string{structValue}).SetConfig(&data.FieldConfig{}),
		data.NewField("quality", nil, []string{"GOOD"}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{Resolution: "RAW"},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
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

func Test_getPropertyValueHistoryFromAliasCaseTable_disassociated_stream(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-from-alias-disassociated.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-table-disassociated", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}
func Test_getPropertyValueHistoryFromAliasCaseTable_disassociated_stream_empty_response(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-from-alias-disassociated-empty-response.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-table-disassociated-empty-response", i)
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
func Test_getPropertyValueHistoryFromAliasCaseTimeSeries_disassociated_stream(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-from-alias-disassociated.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-timeseries-disassociated", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}
func Test_getPropertyValueHistoryFromAliasCaseTimeSeries_disassociated_stream_with_empty_response(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-from-alias-disassociated-empty-response.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-timeseries-disassociated-empty-response", i)
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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-boolean-associated-stream", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_getPropertyValueBooleanFromAliasWithDisassociatedStream(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean-disassociated.json"))
	propTimeSeriesWithoutPropertyId := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeriesWithoutPropertyId, nil)

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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-boolean-with-disassociated-stream", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_getPropertyValueBooleanFromAlias_disassociated_stream_with_empty_response(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean-disassociated-empty-response.json"))
	propTimeSeriesWithoutPropertyId := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series-without-property.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeriesWithoutPropertyId, nil)

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
		fname := fmt.Sprintf("%s-%s.golden", "property-history-values-from-alias-boolean-with-disassociated-stream-empty-response", i)
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

	expectedFrame := data.NewFrame("Demo Turbine Asset 1",
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
		math.MaxInt32,
		math.MaxInt32,
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
