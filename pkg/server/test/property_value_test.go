package test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func Test_property_value_query_by_asset_id_and_property_id(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(float64(23.8)),
				},
			},
			EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}}}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
			}},
		},
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

func Test_property_value_query_by_asset_id_and_property_id_of_flatten_L4E_anomaly_result(t *testing.T) {
	assetId := "1assetid-aaaa-2222-bbbb-3333cccc4444"
	assetPropertyIdQuery := "11propid-aaaa-2222-bbbb-3333cccc4444"
	assetPropertyIdDiagnosticOne := "44fa33e2-b2db-4724-ba03-48ce28902809"
	assetPropertyIdDiagnosticTwo := "3a985085-ea71-4ae6-9395-b65990f58a05"

	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					StringValue: Pointer("{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"),
				},
			},
			EntryId: Pointer(assetId),
		}}}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
					   "propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
					   "flattenL4e":true
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
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func Test_property_value_query_by_asset_id_and_property_id_of_struct_type(t *testing.T) {
	structValue := "{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"

	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					StringValue: Pointer(structValue),
				},
			},
			EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}}}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
			}},
		},
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

func Test_property_value_query_by_alias_associated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer("/amazon/renton/1/rpm"),
		AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
	}, nil)
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(float64(23.8)),
				},
			},
			EntryId: Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}}}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"/amazon/renton/1/rpm"
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
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:       Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				PropertyAlias: Pointer("/amazon/renton/1/rpm"),
			}},
		},
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
func Test_property_value_query_by_alias_disassociated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer("/amazon/renton/1/rpm"),
	}, nil)
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(float64(23.8)),
				},
			},
			EntryId: Pointer("61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef"),
		}}}, nil)

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
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"/amazon/renton/1/rpm"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("/amazon/renton/1/rpm", nil, []float64{23.8}).SetConfig(&data.FieldConfig{Unit: ""}),
		data.NewField("quality", nil, []string{"GOOD"}),
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:       Pointer("61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef"),
				PropertyAlias: Pointer("/amazon/renton/1/rpm"),
			}},
		},
	)
	mockSw.AssertNotCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		mock.Anything,
	)
}
func Test_property_value_query_by_alias_disassociated_stream_with_integer_value(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer("/amazon/renton/1/rpm"),
	}, nil)
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: &iotsitewise.AssetPropertyValue{
				Quality: Pointer("GOOD"),
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					IntegerValue: Pointer(int64(23)),
				},
			},
			EntryId: Pointer("61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef"),
		}}}, nil)

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
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"/amazon/renton/1/rpm"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
		data.NewField("/amazon/renton/1/rpm", nil, []int64{23}).SetConfig(&data.FieldConfig{Unit: ""}),
		data.NewField("quality", nil, []string{"GOOD"}),
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:       Pointer("61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef"),
				PropertyAlias: Pointer("/amazon/renton/1/rpm"),
			}},
		},
	)
	mockSw.AssertNotCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		mock.Anything,
	)
}
func Test_property_value_query_with_empty_property_value_results(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
			AssetPropertyValue: nil,
			EntryId:            Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}}}, nil)
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

	qdr, err := srvr.HandlePropertyValue(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
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
		data.NewField("time", nil, []time.Time{}),
		data.NewField("Wind Speed", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		data.NewField("quality", nil, []string{}),
	)
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
			}},
		},
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
