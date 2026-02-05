package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
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

func mockBatchGetAssetPropertyValue(mockSw *mocks.SitewiseAPIClient, nextToken *string, successEntries []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry, errorEntries []iotsitewisetypes.BatchGetAssetPropertyValueErrorEntry) {
	mockSw.On(
		"BatchGetAssetPropertyValue",
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		NextToken:      nextToken,
		SuccessEntries: successEntries,
		ErrorEntries:   errorEntries,
	}, nil).Once()
}

func mockBatchGetAssetPropertyValueSuccessEntry(entryId *string, valueVariant iotsitewisetypes.Variant, idx int) iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry {
	return iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{
		AssetPropertyValue: &iotsitewisetypes.AssetPropertyValue{
			Quality: iotsitewisetypes.QualityGood,
			Timestamp: &iotsitewisetypes.TimeInNanos{
				OffsetInNanos: Pointer(int32(0)),
				TimeInSeconds: Pointer(int64(1612207200 + idx)),
			},
			Value: &valueVariant,
		},
		EntryId: entryId,
	}
}

func Test_property_value_query_by_asset_id_and_property_id(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewisetypes.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)
	mockDescribeAssetProperty(mockSw)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
		Custom: models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
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

func Test_property_value_query_by_asset_id_and_property_id_of_flatten_L4E_anomaly_result(t *testing.T) {
	assetPropertyIdDiagnosticOne := "44fa33e2-b2db-4724-ba03-48ce28902809"
	assetPropertyIdDiagnosticTwo := "3a985085-ea71-4ae6-9395-b65990f58a05"

	mockSw := &mocks.SitewiseAPIClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewisetypes.Variant{
		StringValue: Pointer("{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "assetId":"%s",
					   "propertyId":"%s",
					   "flattenL4e":true
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
		Custom: models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func Test_property_value_query_by_asset_id_and_property_id_of_struct_type(t *testing.T) {
	structValue := "{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"

	mockSw := &mocks.SitewiseAPIClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewisetypes.Variant{
		StringValue: Pointer(structValue),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
	require.Equal(t, models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId}, frame.Meta.Custom)
	require.Equal(t, 8, len(frame.Fields)) // time, AWS/L4E_ANOMALY_RESULT, 5 parsed fields, quality

	// Create a map for easier field access by name (order-independent)
	fieldMap := make(map[string]*data.Field)
	for _, field := range frame.Fields {
		fieldMap[field.Name] = field
	}

	// Assert time field
	require.Contains(t, fieldMap, "time")
	require.Equal(t, time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC).Unix(), fieldMap["time"].At(0).(time.Time).Unix())

	// Assert AWS/L4E_ANOMALY_RESULT field (original struct value)
	require.Contains(t, fieldMap, "AWS/L4E_ANOMALY_RESULT")
	require.Equal(t, structValue, fieldMap["AWS/L4E_ANOMALY_RESULT"].At(0).(string))

	// Assert quality field
	require.Contains(t, fieldMap, "quality")
	require.Equal(t, "GOOD", fieldMap["quality"].At(0).(string))

	// Assert parsed JSON fields with their values
	require.Contains(t, fieldMap, "prediction")
	require.InDelta(t, float64(0), fieldMap["prediction"].At(0).(float64), 0.0001)

	require.Contains(t, fieldMap, "prediction_reason")
	require.Equal(t, "NO_ANOMALY_DETECTED", fieldMap["prediction_reason"].At(0).(string))

	require.Contains(t, fieldMap, "anomaly_score")
	require.InDelta(t, 0.2674, fieldMap["anomaly_score"].At(0).(float64), 0.0001)

	require.Contains(t, fieldMap, "contrib_Demo Turbine Asset 1_3a985085-ea71-4ae6-9395-b65990f58a05")
	require.InDelta(t, 44.856, fieldMap["contrib_Demo Turbine Asset 1_3a985085-ea71-4ae6-9395-b65990f58a05"].At(0).(float64), 0.001)

	require.Contains(t, fieldMap, "contrib_Demo Turbine Asset 1_44fa33e2-b2db-4724-ba03-48ce28902809")
	require.InDelta(t, 55.144, fieldMap["contrib_Demo Turbine Asset 1_44fa33e2-b2db-4724-ba03-48ce28902809"].At(0).(float64), 0.001)

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
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

func Test_property_value_query_by_alias_associated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer(mockPropertyAlias),
		AssetId:    Pointer(mockAssetId),
		PropertyId: Pointer(mockPropertyId),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewisetypes.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)
	mockDescribeAssetProperty(mockSw)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"%s"
					}`, mockPropertyAlias)),
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
		Custom: models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeries",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
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
func Test_property_value_query_by_alias_disassociated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer(mockPropertyAlias),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockPropertyAliasEntryId, iotsitewisetypes.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"%s"
					}`, mockPropertyAlias)),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField(mockPropertyAlias, nil, []float64{23.8}).SetConfig(&data.FieldConfig{Unit: ""}),
		data.NewField("quality", nil, []string{"GOOD"}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{EntryId: *mockPropertyAliasEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeries",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:       mockPropertyAliasEntryId,
				PropertyAlias: Pointer(mockPropertyAlias),
			}},
		},
	)
	mockSw.AssertNotCalled(t,
		"DescribeAssetProperty",
		mock.Anything,
		mock.Anything,
	)
}
func Test_property_value_query_by_alias_disassociated_stream_with_integer_value(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer(mockPropertyAlias),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockPropertyAliasEntryId, iotsitewisetypes.Variant{
		IntegerValue: Pointer(int32(23)),
	}, 0)
	mockBatchGetAssetPropertyValue(mockSw, nil, []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{successEntry}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(fmt.Sprintf(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"%s"
					}`, mockPropertyAlias)),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField(mockPropertyAlias, nil, []int64{23}).SetConfig(&data.FieldConfig{Unit: ""}),
		data.NewField("quality", nil, []string{"GOOD"}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{EntryId: *mockPropertyAliasEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeries",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:       mockPropertyAliasEntryId,
				PropertyAlias: Pointer(mockPropertyAlias),
			}},
		},
	)
	mockSw.AssertNotCalled(t,
		"DescribeAssetProperty",
		mock.Anything,
		mock.Anything,
	)
}
func Test_property_value_query_with_empty_property_value_results(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("BatchGetAssetPropertyValue", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: nil,
			EntryId:            mockAssetPropertyEntryId,
		}}}, nil)
	mockDescribeAssetProperty(mockSw)
	mockDescribeAsset(mockSw)
	mockDescribeAssetModel(mockSw)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
		data.NewField("time", nil, []time.Time{}),
		data.NewField("Wind Speed", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		data.NewField("quality", nil, []string{}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValue",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []iotsitewisetypes.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
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

func Test_property_value_query_with_batched_queries(t *testing.T) {
	tests := []batch_test{
		{
			name:           "query by multiple assetIds and one propertyId",
			numAssetIds:    api.BatchGetAssetPropertyValueMaxEntries + 1,
			numPropertyIds: 1,
		},
		{
			name:           "query by one assetId and multiple propertyIds",
			numAssetIds:    1,
			numPropertyIds: api.BatchGetAssetPropertyValueMaxEntries + 1,
		},
		{
			name:           "query by multiple assetIds and multiple propertyIds",
			numAssetIds:    api.BatchGetAssetPropertyValueMaxEntries + 1,
			numPropertyIds: api.BatchGetAssetPropertyValueMaxEntries + 1,
		},
		{
			name:               "query by multiple property aliases",
			numPropertyAliases: api.BatchGetAssetPropertyValueMaxEntries + 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseAPIClient{}
			mockedSuccessEntries := []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{}
			numBatch := 0

			if tc.numPropertyAliases > 0 {
				propertyAliases := generateIds(tc.numPropertyAliases, mockPropertyAlias)
				for p, propertyAlias := range propertyAliases {
					// Build the success entry based on the propertyAlias
					entryId := util.GetEntryIdFromPropertyAlias(propertyAlias)
					successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewisetypes.Variant{
						DoubleValue: Pointer(float64(23.8) + float64(p)),
					}, p)
					mockedSuccessEntries = append(mockedSuccessEntries, successEntry)

					isLastBatch := p == tc.numPropertyAliases-1
					// When batch is complete mock the call with the success entries
					if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
						numBatch++
						mockBatchGetAssetPropertyValue(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
						mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{}
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
						successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewisetypes.Variant{
							DoubleValue: Pointer(float64(23.8) + float64(p)),
						}, p)
						mockedSuccessEntries = append(mockedSuccessEntries, successEntry)

						isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
						// When batch is complete mock the call with the success entries
						if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
							numBatch++
							mockBatchGetAssetPropertyValue(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
							mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{}
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

			qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{},
				Queries: []backend.DataQuery{{
					QueryType:     models.QueryTypePropertyValue,
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
				if (i+1)%api.BatchGetAssetPropertyValueMaxEntries == 0 {
					numBatch++
				}
			}

			mockSw.AssertExpectations(t)
		})
	}
}

func Test_property_value_query_with_batched_queries_with_error(t *testing.T) {
	tc := batch_test{
		name:           "batch query with one error",
		numAssetIds:    5,
		numPropertyIds: 5,
	}

	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseAPIClient{}

		mockedSuccessEntries := []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{}
		mockedErrorEntries := []iotsitewisetypes.BatchGetAssetPropertyValueErrorEntry{}
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
					mockedErrorEntries = append(mockedErrorEntries, iotsitewisetypes.BatchGetAssetPropertyValueErrorEntry{
						ErrorCode:    iotsitewisetypes.BatchGetAssetPropertyValueErrorCodeResourceNotFoundException,
						ErrorMessage: Pointer("Asset property not found."),
						EntryId:      entryId,
					})
				} else {
					successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewisetypes.Variant{
						DoubleValue: Pointer(float64(23.8) + float64(p)),
					}, p)
					mockedSuccessEntries = append(mockedSuccessEntries, successEntry)
				}

				isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
				// When batch is complete mock the call with the success entries
				if len(mockedSuccessEntries)+len(mockedErrorEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
					numBatch++
					mockBatchGetAssetPropertyValue(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, mockedErrorEntries)
					// Reset for next batch
					mockedSuccessEntries = []iotsitewisetypes.BatchGetAssetPropertyValueSuccessEntry{}
					mockedErrorEntries = []iotsitewisetypes.BatchGetAssetPropertyValueErrorEntry{}
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

		qdr, err := srvr.HandlePropertyValue(context.Background(), query)
		require.Nil(t, err)
		_, ok := qdr.Responses["A"]
		require.True(t, ok)

		expectedNumFrames := tc.numAssetIds * tc.numPropertyIds
		require.Len(t, qdr.Responses["A"].Frames, expectedNumFrames)

		// Check for the error
		numErrors := 0
		for _, f := range qdr.Responses["A"].Frames {
			if len(f.Meta.Notices) > 0 {
				expectedErrorFrame := data.NewFrame("Demo Turbine Asset 1").SetMeta(&data.FrameMeta{
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
