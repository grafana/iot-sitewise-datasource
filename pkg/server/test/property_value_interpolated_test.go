package test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPropertyValueInterpolatedQuery(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}

	mockSw.On(
		"GetInterpolatedAssetPropertyValuesPageAggregation",
		mock.Anything,
		mock.MatchedBy(func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
			return *input.AssetId == "1assetid-aaaa-2222-bbbb-3333cccc4444"
		}),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		NextToken: Pointer("asset1-next-token"),
		InterpolatedAssetPropertyValues: []*iotsitewise.InterpolatedAssetPropertyValue{
			{
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207100)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(1.1),
				},
			},
		},
	}, nil)

	mockSw.On(
		"GetInterpolatedAssetPropertyValuesPageAggregation",
		mock.Anything,
		mock.MatchedBy(func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
			return *input.AssetId == "2assetid-aaaa-2222-bbbb-3333cccc4444"
		}),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		NextToken: Pointer("asset2-next-token"),
		InterpolatedAssetPropertyValues: []*iotsitewise.InterpolatedAssetPropertyValue{
			{
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(2.2),
				},
			},
		},
	}, nil)

	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		PropertyId: aws.String("11propid-aaaa-2222-bbbb-3333cccc4444"),
	}, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			DataType: Pointer("DOUBLE"),
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
		},
	}, nil)

	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String("2assetid-aaaa-2222-bbbb-3333cccc4444"),
		PropertyId: aws.String("11propid-aaaa-2222-bbbb-3333cccc4444"),
	}, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer("2assetid-aaaa-2222-bbbb-3333cccc4444"),
		AssetName: Pointer("Demo Turbine Asset 2"),
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

	query := &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyInterpolated,
				TimeRange: timeRange,
				JSON: []byte(`{
					"assetIds": ["1assetid-aaaa-2222-bbbb-3333cccc4444", "2assetid-aaaa-2222-bbbb-3333cccc4444"],
					"propertyId": "11propid-aaaa-2222-bbbb-3333cccc4444",
					"resolution": "1m"
				}`),
			},
		},
	}

	qdr, err := srvr.HandleInterpolatedPropertyValue(context.Background(), query)
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.Equal(t, 2, len(qdr.Responses["A"].Frames))
	require.NotNil(t, qdr.Responses["A"].Frames[0])
	require.NotNil(t, qdr.Responses["A"].Frames[1])

	sort.Slice(qdr.Responses["A"].Frames, func(a, b int) bool {
		return qdr.Responses["A"].Frames[a].Name < qdr.Responses["A"].Frames[b].Name
	})

	expectedFrames := data.Frames{
		data.NewFrame("Demo Turbine Asset 1",
			data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 18, 20, 0, time.UTC)}),
			data.NewField("Wind Speed", nil, []float64{1.1}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "asset1-next-token",
				EntryId:    "1assetid-aaaa-2222-bbbb-3333cccc4444",
				Resolution: "1m",
				Aggregates: []string{},
			},
		}),
		data.NewFrame("Demo Turbine Asset 2",
			data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 20, 0, 0, time.UTC)}),
			data.NewField("Wind Speed", nil, []float64{2.2}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "asset2-next-token",
				EntryId:    "2assetid-aaaa-2222-bbbb-3333cccc4444",
				Resolution: "1m",
				Aggregates: []string{},
			},
		}),
	}

	if diff := cmp.Diff(expectedFrames, qdr.Responses["A"].Frames, data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func TestPropertyValueInterpolatedQueryWithPropertyAlias(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}

	mockSw.On(
		"GetInterpolatedAssetPropertyValuesPageAggregation",
		mock.Anything,
		mock.MatchedBy(func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
			return *input.PropertyAlias == "/turbine_1/wind_speed"
		}),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		NextToken: Pointer("asset1-next-token"),
		InterpolatedAssetPropertyValues: []*iotsitewise.InterpolatedAssetPropertyValue{
			{
				Timestamp: &iotsitewise.TimeInNanos{
					OffsetInNanos: Pointer(int64(0)),
					TimeInSeconds: Pointer(int64(1612207100)),
				},
				Value: &iotsitewise.Variant{
					DoubleValue: Pointer(1.1),
				},
			},
		},
	}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}
	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	query := &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyInterpolated,
				TimeRange: timeRange,
				JSON: []byte(`{
					"propertyAlias": "/turbine_1/wind_speed",
					"resolution": "1m"
				}`),
			},
		},
	}

	qdr, err := srvr.HandleInterpolatedPropertyValue(context.Background(), query)
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.Equal(t, 1, len(qdr.Responses["A"].Frames))
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrames := data.Frames{
		data.NewFrame("/turbine_1/wind_speed",
			data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 19, 18, 20, 0, time.UTC)}),
			data.NewField("/turbine_1/wind_speed", nil, []float64{1.1}).SetConfig(&data.FieldConfig{}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "asset1-next-token",
				EntryId:    "/turbine_1/wind_speed",
				Resolution: "1m",
				Aggregates: []string{},
			},
		}),
	}

	if diff := cmp.Diff(expectedFrames, qdr.Responses["A"].Frames, data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}
