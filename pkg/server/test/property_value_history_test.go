package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/google/go-cmp/cmp"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw *mocks.SitewiseAPIClient, nextToken *string, successEntries []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry, errorEntries []iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorEntry) {
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		NextToken:      nextToken,
		SuccessEntries: successEntries,
		ErrorEntries:   errorEntries,
	}, nil).Once()
}

func mockBatchGetAssetPropertyValueHistorySuccessEntry(entryId *string, idx int) iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry {
	return iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{
		AssetPropertyValueHistory: []iotsitewisetypes.AssetPropertyValue{
			{
				Quality: iotsitewisetypes.QualityGood,
				Timestamp: &iotsitewisetypes.TimeInNanos{
					OffsetInNanos: Pointer(int32(0)),
					TimeInSeconds: Pointer(int64(1612207200 + idx)),
				},
				Value: &iotsitewisetypes.Variant{
					DoubleValue: Pointer(23.8 + float64(idx)),
				},
			},
		},
		EntryId: entryId,
	}
}

func Test_get_property_value_history_with_default_aka_table_response_format(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}

	successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(mockAssetPropertyEntryId, 0)
	mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{successEntry}, nil)

	mockDescribeAssetProperty(mockSw)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
						 "propertyId":"%s"
					}`, mockAssetId, mockPropertyId)),
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
		Custom: models.SitewiseCustomMeta{Resolution: "RAW", EntryId: *mockAssetPropertyEntryId},
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
		"DescribeAssetProperty",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		},
	)
}

func Test_get_property_value_history_with_time_series_response_format(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}

	successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(mockAssetPropertyEntryId, 0)
	mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{successEntry}, nil)

	mockDescribeAssetProperty(mockSw)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
						 "propertyId":"%s",
						 "responseFormat":"timeseries"
					}`, mockAssetId, mockPropertyId)),
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
		Type:        data.FrameTypeTimeSeriesWide,
		TypeVersion: data.FrameTypeVersion{0, 1},
		Custom:      models.SitewiseCustomMeta{Resolution: "RAW", EntryId: *mockAssetPropertyEntryId},
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
		"DescribeAssetProperty",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		},
	)
}

func Test_getPropertyValueBoolean(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
	assetPropertyIdDiagnosticOne := "44fa33e2-b2db-4724-ba03-48ce28902809"
	assetPropertyIdDiagnosticTwo := "3a985085-ea71-4ae6-9395-b65990f58a05"

	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{
			{
				AssetPropertyValueHistory: []iotsitewisetypes.AssetPropertyValue{
					{
						Quality: iotsitewisetypes.QualityGood,
						Timestamp: &iotsitewisetypes.TimeInNanos{
							OffsetInNanos: Pointer(int32(0)),
							TimeInSeconds: Pointer(int64(1612207200)),
						},
						Value: &iotsitewisetypes.Variant{
							StringValue: Pointer("{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"),
						},
					},
				},
				EntryId: mockAssetPropertyEntryId,
			},
		},
	}, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == mockPropertyId
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(mockAssetId),
		AssetName: Pointer("Demo Turbine Asset 1"),
		CompositeModel: &iotsitewisetypes.CompositeModelProperty{
			Name: Pointer("prediction1"),
			AssetProperty: &iotsitewisetypes.Property{
				Name:     Pointer("AWS/L4E_ANOMALY_RESULT"),
				DataType: iotsitewisetypes.PropertyDataTypeStruct,
			},
		},
	}, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == assetPropertyIdDiagnosticOne
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewisetypes.Property{
			Id:       Pointer(assetPropertyIdDiagnosticOne),
			DataType: iotsitewisetypes.PropertyDataTypeDouble,
			Name:     Pointer("Torque"),
		},
	}, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == assetPropertyIdDiagnosticTwo
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewisetypes.Property{
			Id:       Pointer(assetPropertyIdDiagnosticTwo),
			DataType: iotsitewisetypes.PropertyDataTypeDouble,
			Name:     Pointer("RPM"),
		},
	}, nil)

	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
					   "propertyId":"%s",
					   "flattenL4E": true
					}`, mockAssetId, mockPropertyId)),
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
		Custom: models.SitewiseCustomMeta{Resolution: "RAW", EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func Test_get_property_value_history_with_struct_type(t *testing.T) {
	structValue := "{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"

	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On(
		"BatchGetAssetPropertyValueHistoryPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{
			{
				AssetPropertyValueHistory: []iotsitewisetypes.AssetPropertyValue{
					{
						Quality: iotsitewisetypes.QualityGood,
						Timestamp: &iotsitewisetypes.TimeInNanos{
							OffsetInNanos: Pointer(int32(0)),
							TimeInSeconds: Pointer(int64(1612207200)),
						},
						Value: &iotsitewisetypes.Variant{
							StringValue: Pointer(structValue),
						},
					},
				},
				EntryId: mockAssetPropertyEntryId,
			},
		},
	}, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(mockAssetId),
		AssetName: Pointer("Demo Turbine Asset 1"),
		CompositeModel: &iotsitewisetypes.CompositeModelProperty{
			Name: Pointer("prediction1"),
			AssetProperty: &iotsitewisetypes.Property{
				Name:     Pointer("AWS/L4E_ANOMALY_RESULT"),
				DataType: iotsitewisetypes.PropertyDataTypeStruct,
			},
		},
	}, nil)

	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
					   "propertyId":"%s"
					}`, mockAssetId, mockPropertyId)),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	frame := qdr.Responses["A"].Frames[0]
	require.Equal(t, "Demo Turbine Asset 1", frame.Name)
	require.Equal(t, models.SitewiseCustomMeta{Resolution: "RAW", EntryId: *mockAssetPropertyEntryId}, frame.Meta.Custom)
	require.Equal(t, 8, len(frame.Fields)) // time, AWS/L4E_ANOMALY_RESULT, 5 parsed fields, quality

	// Create a map for easier field access by name (order-independent)
	fieldMap := make(map[string]*data.Field)
	for _, field := range frame.Fields {
		fieldMap[field.Name] = field
	}

	// Assert time field
	require.Contains(t, fieldMap, "time")
	require.Equal(t, 1, fieldMap["time"].Len())
	require.Equal(t, time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC).Unix(), fieldMap["time"].At(0).(time.Time).Unix())

	// Assert AWS/L4E_ANOMALY_RESULT field (original struct value)
	require.Contains(t, fieldMap, "AWS/L4E_ANOMALY_RESULT")
	require.Equal(t, 1, fieldMap["AWS/L4E_ANOMALY_RESULT"].Len())
	require.Equal(t, structValue, fieldMap["AWS/L4E_ANOMALY_RESULT"].At(0).(string))

	// Assert quality field
	require.Contains(t, fieldMap, "quality")
	require.Equal(t, 1, fieldMap["quality"].Len())
	require.Equal(t, "GOOD", fieldMap["quality"].At(0).(string))

	// Assert parsed JSON fields with their values
	require.Contains(t, fieldMap, "prediction")
	require.Equal(t, 1, fieldMap["prediction"].Len())
	require.InDelta(t, float64(0), fieldMap["prediction"].At(0).(float64), 0.0001)

	require.Contains(t, fieldMap, "prediction_reason")
	require.Equal(t, 1, fieldMap["prediction_reason"].Len())
	require.Equal(t, "NO_ANOMALY_DETECTED", fieldMap["prediction_reason"].At(0).(string))

	require.Contains(t, fieldMap, "anomaly_score")
	require.Equal(t, 1, fieldMap["anomaly_score"].Len())
	require.InDelta(t, 0.2674, fieldMap["anomaly_score"].At(0).(float64), 0.0001)

	require.Contains(t, fieldMap, "contrib_Demo Turbine Asset 1_3a985085-ea71-4ae6-9395-b65990f58a05")
	require.Equal(t, 1, fieldMap["contrib_Demo Turbine Asset 1_3a985085-ea71-4ae6-9395-b65990f58a05"].Len())
	require.InDelta(t, 44.856, fieldMap["contrib_Demo Turbine Asset 1_3a985085-ea71-4ae6-9395-b65990f58a05"].At(0).(float64), 0.001)

	require.Contains(t, fieldMap, "contrib_Demo Turbine Asset 1_44fa33e2-b2db-4724-ba03-48ce28902809")
	require.Equal(t, 1, fieldMap["contrib_Demo Turbine Asset 1_44fa33e2-b2db-4724-ba03-48ce28902809"].Len())
	require.InDelta(t, 55.144, fieldMap["contrib_Demo Turbine Asset 1_44fa33e2-b2db-4724-ba03-48ce28902809"].At(0).(float64), 0.001)

	mockSw.AssertExpectations(t)
}

func Test_getPropertyValueHistoryFromAliasCaseTable(t *testing.T) {
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	propVals.SuccessEntries[0].EntryId = util.GetEntryIdFromAssetProperty(*propTimeSeries.AssetId, *propTimeSeries.PropertyId)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	propVals.SuccessEntries[0].EntryId = util.GetEntryIdFromAssetProperty(*propTimeSeries.AssetId, *propTimeSeries.PropertyId)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

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
	propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values-from-alias-boolean.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-is-windy.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	propVals.SuccessEntries[0].EntryId = util.GetEntryIdFromAssetProperty(*propTimeSeries.AssetId, *propTimeSeries.PropertyId)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeriesWithoutPropertyId, nil)

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
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&propTimeSeriesWithoutPropertyId, nil)

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
	mockSw := &mocks.SitewiseAPIClient{}

	successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(mockAssetPropertyEntryId, 0)
	mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{successEntry}, nil)

	mockDescribeAssetProperty(mockSw)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

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
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
						 "propertyId":"%s",
						 "responseFormat":"timeseries"
					}`, mockAssetId, mockPropertyId)),
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
		Type:        data.FrameTypeTimeSeriesWide,
		TypeVersion: data.FrameTypeVersion{0, 1},
		Custom:      models.SitewiseCustomMeta{Resolution: "RAW", EntryId: *mockAssetPropertyEntryId},
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
		"DescribeAssetProperty",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		},
	)
}

func Test_get_property_value_history_with_batched_queries(t *testing.T) {
	tests := []batch_test{
		{
			name:           "query by multiple assetIds and one propertyId",
			numAssetIds:    api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
			numPropertyIds: 1,
		},
		{
			name:           "query by one assetId and multiple propertyIds",
			numAssetIds:    1,
			numPropertyIds: api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
		},
		{
			name:           "query by multiple assetIds and multiple propertyIds",
			numAssetIds:    api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
			numPropertyIds: api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
		},
		{
			name:               "query by multiple property aliases",
			numPropertyAliases: api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseAPIClient{}
			mockedSuccessEntries := []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{}
			numBatch := 0

			if tc.numPropertyAliases > 0 {
				propertyAliases := generateIds(tc.numPropertyAliases, mockPropertyAlias)
				for p, propertyAlias := range propertyAliases {
					// Build the success entry based on the propertyAlias for disassociated data streams
					entryId := util.GetEntryIdFromPropertyAlias(propertyAlias)
					successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(entryId, p)
					mockedSuccessEntries = append(mockedSuccessEntries, successEntry)

					isLastBatch := p == tc.numPropertyAliases-1
					// When batch is complete mock the History call with the success entries
					if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueHistoryMaxEntries || isLastBatch {
						numBatch++
						mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
						mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{}
					}

					mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
						Alias: Pointer(propertyAlias),
					}, nil)
				}
			} else {
				assetIds := generateIds(tc.numAssetIds, mockAssetId)
				propertyIds := generateIds(tc.numPropertyIds, mockPropertyId)
				for a, assetId := range assetIds {
					for p, propertyId := range propertyIds {
						// Build the success entry based on the assetId and propertyId
						entryId := util.GetEntryIdFromAssetProperty(assetId, propertyId)
						successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(entryId, p)
						mockedSuccessEntries = append(mockedSuccessEntries, successEntry)

						isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
						// When batch is complete mock the History call with the success entries
						if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueHistoryMaxEntries || isLastBatch {
							numBatch++
							mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
							mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{}
						}
					}
				}
				mockDescribeAssetProperty(mockSw)
				mockDescribeAsset(mockSw)
				mockDescribeAssetModel(mockSw)
			}

			srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

			sitewise.GetCache = func() *cache.Cache {
				return cache.New(cache.DefaultExpiration, cache.NoExpiration)
			}

			var baseQuery models.BaseQuery
			if tc.numPropertyAliases > 0 {
				baseQuery = models.BaseQuery{
					AwsRegion:       testdata.AwsRegion,
					PropertyAliases: generateIds(tc.numPropertyAliases, mockPropertyAlias),
				}
			} else {
				baseQuery = models.BaseQuery{
					AwsRegion:   testdata.AwsRegion,
					AssetIds:    generateIds(tc.numAssetIds, mockAssetId),
					PropertyIds: generateIds(tc.numPropertyIds, mockPropertyId),
				}
			}

			qdr, err := srvr.HandlePropertyValueHistory(context.Background(), &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{},
				Queries: []backend.DataQuery{{
					QueryType:     models.QueryTypePropertyValueHistory,
					RefID:         "A",
					MaxDataPoints: 100,
					Interval:      1000,
					TimeRange:     timeRange,
					JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
						BaseQuery: baseQuery,
					}),
				}},
			})

			require.Nil(t, err)
			_, ok := qdr.Responses["A"]
			require.True(t, ok)
			var expectedNumFrames int
			if tc.numPropertyAliases > 0 {
				expectedNumFrames = tc.numPropertyAliases
			} else {
				expectedNumFrames = tc.numAssetIds * tc.numPropertyIds
			}
			require.Len(t, qdr.Responses["A"].Frames, expectedNumFrames)

			numBatch = 1
			for i, f := range qdr.Responses["A"].Frames {
				require.NotNil(t, f)
				expectedNextToken := fmt.Sprintf("some-next-token-%d", numBatch)
				var entryId string
				if tc.numPropertyAliases > 0 {
					propertyAlias := fmt.Sprintf("%s%d", mockPropertyAlias, i+1)
					entryId = *util.GetEntryIdFromPropertyAlias(propertyAlias)
				} else {
					assetId := fmt.Sprintf("%s%d", mockAssetId, int(math.Floor(float64(i)/float64(tc.numPropertyIds)))+1)
					propertyId := fmt.Sprintf("%s%d", mockPropertyId, i%tc.numPropertyIds+1)
					entryId = *util.GetEntryIdFromAssetProperty(assetId, propertyId)
				}
				require.Equal(t, entryId, f.Meta.Custom.(models.SitewiseCustomMeta).EntryId)
				require.Equal(t, expectedNextToken, f.Meta.Custom.(models.SitewiseCustomMeta).NextToken)
				// Increment to next batch
				if (i+1)%api.BatchGetAssetPropertyAggregatesMaxEntries == 0 {
					numBatch++
				}
			}

			mockSw.AssertExpectations(t)

			for i, dr := range qdr.Responses {
				fname := fmt.Sprintf("%s-%s.golden", fmt.Sprintf("property-history-values-%s", tc.name), i)
				experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
			}
		})
	}
}

func Test_get_property_value_history_with_batched_queries_with_error(t *testing.T) {
	tc := batch_test{
		name:           "batch history query with one error",
		numAssetIds:    api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
		numPropertyIds: api.BatchGetAssetPropertyValueHistoryMaxEntries + 1,
	}

	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseAPIClient{}

		mockedSuccessEntries := []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{}
		mockedErrorEntries := []iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorEntry{}
		numBatch := 0
		errorIndex := 20

		assetIds := generateIds(tc.numAssetIds, mockAssetId)
		propertyIds := generateIds(tc.numPropertyIds, mockPropertyId)
		for a, assetId := range assetIds {
			for p, propertyId := range propertyIds {
				// Build the success entry based on the assetId and propertyId
				entryId := util.GetEntryIdFromAssetProperty(assetId, propertyId)

				// Build one error entry
				if a*tc.numPropertyIds+p == errorIndex {
					mockedErrorEntries = append(mockedErrorEntries, iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorEntry{
						ErrorCode:    iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorCodeResourceNotFoundException,
						ErrorMessage: Pointer("Asset property not found."),
						EntryId:      entryId,
					})
				} else {
					successEntry := mockBatchGetAssetPropertyValueHistorySuccessEntry(entryId, p)
					mockedSuccessEntries = append(mockedSuccessEntries, successEntry)
				}

				isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
				// When batch is complete mock the History call with the success entries
				if len(mockedSuccessEntries)+len(mockedErrorEntries) == api.BatchGetAssetPropertyValueHistoryMaxEntries || isLastBatch {
					numBatch++
					mockBatchGetAssetPropertyValueHistoryPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, mockedErrorEntries)
					// Reset for next batch
					mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry{}
					mockedErrorEntries = []iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorEntry{}
				}
			}
		}

		mockDescribeAssetProperty(mockSw)
		mockDescribeAsset(mockSw)
		mockDescribeAssetModel(mockSw)

		srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

		sitewise.GetCache = func() *cache.Cache {
			return cache.New(cache.DefaultExpiration, cache.NoExpiration)
		}

		var baseQuery models.BaseQuery
		if tc.numPropertyAliases > 0 {
			baseQuery = models.BaseQuery{
				AwsRegion:       testdata.AwsRegion,
				PropertyAliases: generateIds(tc.numPropertyAliases, mockPropertyAlias),
			}
		} else {
			baseQuery = models.BaseQuery{
				AwsRegion:   testdata.AwsRegion,
				AssetIds:    generateIds(tc.numAssetIds, mockAssetId),
				PropertyIds: generateIds(tc.numPropertyIds, mockPropertyId),
			}
		}

		query := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{},
			Queries: []backend.DataQuery{
				{
					RefID:     "A",
					QueryType: models.QueryTypePropertyAggregate,
					TimeRange: timeRange,
					JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
						BaseQuery: baseQuery,
					}),
				},
			},
		}

		qdr, err := srvr.HandlePropertyValueHistory(context.Background(), query)
		require.Nil(t, err)
		_, ok := qdr.Responses["A"]
		require.True(t, ok)

		expectedNumFrames := tc.numAssetIds * tc.numPropertyIds
		require.Len(t, qdr.Responses["A"].Frames, expectedNumFrames)

		// Check for the error
		numErrors := 0
		for _, f := range qdr.Responses["A"].Frames {
			if len(f.Meta.Notices) > 0 {
				expectedErrorFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed").SetMeta(&data.FrameMeta{
					Notices: []data.Notice{{Severity: data.NoticeSeverityError, Text: "Asset property not found."}},
				},
				)
				if diff := cmp.Diff(expectedErrorFrame, f, data.FrameTestCompareOptions()...); diff != "" {
					t.Errorf("Result mismatch (-want +got):\n%s", diff)
				}
				numErrors++
			}
		}
		require.True(t, numErrors == 1)

		mockSw.AssertExpectations(t)
	})
}
