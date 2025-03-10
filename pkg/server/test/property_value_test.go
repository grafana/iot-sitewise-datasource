package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func mockBatchGetAssetPropertyValueWithContext(mockSw *mocks.SitewiseClient, nextToken *string, successEntries []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry, errorEntries []*iotsitewise.BatchGetAssetPropertyValueErrorEntry) {
	mockSw.On(
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyValueOutput{
		NextToken:      nextToken,
		SuccessEntries: successEntries,
		ErrorEntries:   errorEntries,
	}, nil).Once()
}

func mockBatchGetAssetPropertyValueSuccessEntry(entryId *string, valueVariant iotsitewise.Variant, idx int) iotsitewise.BatchGetAssetPropertyValueSuccessEntry {
	return iotsitewise.BatchGetAssetPropertyValueSuccessEntry{
		AssetPropertyValue: &iotsitewise.AssetPropertyValue{
			Quality: Pointer("GOOD"),
			Timestamp: &iotsitewise.TimeInNanos{
				OffsetInNanos: Pointer(int64(0)),
				TimeInSeconds: Pointer(int64(1612207200 + idx)),
			},
			Value: &valueVariant,
		},
		EntryId: entryId,
	}
}

func Test_property_value_query_by_asset_id_and_property_id(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewise.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)
	mockDescribeAssetPropertyWithContext(mockSw)

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
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
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

	mockSw := &mocks.SitewiseClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewise.Variant{
		StringValue: Pointer("{\"timestamp\":\"2021-02-01T19:20:00.000000\",\"prediction\":0,\"prediction_reason\":\"NO_ANOMALY_DETECTED\",\"anomaly_score\":0.2674,\"diagnostics\":[{\"name\":\"3a985085-ea71-4ae6-9395-b65990f58a05\\\\3a985085-ea71-4ae6-9395-b65990f58a05\",\"value\":0.44856},{\"name\":\"44fa33e2-b2db-4724-ba03-48ce28902809\\\\44fa33e2-b2db-4724-ba03-48ce28902809\",\"value\":0.55144}]}"),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetPropertyInput) bool {
		return req.PropertyId != nil && *req.PropertyId == mockPropertyId
	})).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(mockAssetId),
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

	mockSw := &mocks.SitewiseClient{}
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewise.Variant{
		StringValue: Pointer(structValue),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(mockAssetId),
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
		data.NewField("AWS/L4E_ANOMALY_RESULT", nil, []string{structValue}).SetConfig(&data.FieldConfig{}),
		data.NewField("quality", nil, []string{"GOOD"}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{EntryId: *mockAssetPropertyEntryId},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		},
	)
}

func Test_property_value_query_by_alias_associated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer(mockPropertyAlias),
		AssetId:    Pointer(mockAssetId),
		PropertyId: Pointer(mockPropertyId),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockAssetPropertyEntryId, iotsitewise.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)
	mockDescribeAssetPropertyWithContext(mockSw)

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
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		},
	)
}
func Test_property_value_query_by_alias_disassociated_stream(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer(mockPropertyAlias),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockPropertyAliasEntryId, iotsitewise.Variant{
		DoubleValue: Pointer(23.8),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)

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
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:       mockPropertyAliasEntryId,
				PropertyAlias: Pointer(mockPropertyAlias),
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
		Alias: Pointer(mockPropertyAlias),
	}, nil)
	successEntry := mockBatchGetAssetPropertyValueSuccessEntry(mockPropertyAliasEntryId, iotsitewise.Variant{
		IntegerValue: Pointer(int64(23)),
	}, 0)
	mockBatchGetAssetPropertyValueWithContext(mockSw, nil, []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{&successEntry}, nil)

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
		"DescribeTimeSeriesWithContext",
		mock.Anything,
		&iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
	)
	mockSw.AssertCalled(t,
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:       mockPropertyAliasEntryId,
				PropertyAlias: Pointer(mockPropertyAlias),
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
			EntryId:            mockAssetPropertyEntryId,
		}}}, nil)
	mockDescribeAssetPropertyWithContext(mockSw)

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
		"BatchGetAssetPropertyValueWithContext",
		mock.Anything,
		&iotsitewise.BatchGetAssetPropertyValueInput{
			Entries: []*iotsitewise.BatchGetAssetPropertyValueEntry{{
				EntryId:    mockAssetPropertyEntryId,
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}},
		},
	)
	mockSw.AssertCalled(t,
		"DescribeAssetPropertyWithContext",
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
			mockSw := &mocks.SitewiseClient{}
			mockedSuccessEntries := []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{}
			numBatch := 0

			if tc.numPropertyAliases > 0 {
				propertyAliases := generateIds(tc.numPropertyAliases, mockPropertyAlias)
				for p, propertyAlias := range propertyAliases {
					// Build the success entry based on the propertyAlias
					entryId := util.GetEntryIdFromPropertyAlias(propertyAlias)
					successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewise.Variant{
						DoubleValue: Pointer(float64(23.8) + float64(p)),
					}, p)
					mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)

					isLastBatch := p == tc.numPropertyAliases-1
					// When batch is complete mock the call with the success entries
					if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
						numBatch++
						mockBatchGetAssetPropertyValueWithContext(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
						mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{}
					}

					mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
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
						successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewise.Variant{
							DoubleValue: Pointer(float64(23.8) + float64(p)),
						}, p)
						mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)

						isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
						// When batch is complete mock the call with the success entries
						if len(mockedSuccessEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
							numBatch++
							mockBatchGetAssetPropertyValueWithContext(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
							mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{}
						}
					}
				}
				mockDescribeAssetPropertyWithContext(mockSw)
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
		mockSw := &mocks.SitewiseClient{}

		mockedSuccessEntries := []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{}
		mockedErrorEntries := []*iotsitewise.BatchGetAssetPropertyValueErrorEntry{}
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
					mockedErrorEntries = append(mockedErrorEntries, &iotsitewise.BatchGetAssetPropertyValueErrorEntry{
						ErrorCode:    Pointer("404"),
						ErrorMessage: Pointer("Asset property not found."),
						EntryId:      entryId,
					})
				} else {
					successEntry := mockBatchGetAssetPropertyValueSuccessEntry(entryId, iotsitewise.Variant{
						DoubleValue: Pointer(float64(23.8) + float64(p)),
					}, p)
					mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)
				}

				isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
				// When batch is complete mock the call with the success entries
				if len(mockedSuccessEntries)+len(mockedErrorEntries) == api.BatchGetAssetPropertyValueMaxEntries || isLastBatch {
					numBatch++
					mockBatchGetAssetPropertyValueWithContext(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, mockedErrorEntries)
					// Reset for next batch
					mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{}
					mockedErrorEntries = []*iotsitewise.BatchGetAssetPropertyValueErrorEntry{}
				}
			}
		}

		mockDescribeAssetPropertyWithContext(mockSw)

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
