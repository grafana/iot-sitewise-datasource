package test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func Test_property_value_aggregate_query_by_asset_id_and_property_id(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On(
		"BatchGetAssetPropertyAggregatesPageAggregation",
		mock.Anything,
		mock.MatchedBy(func(input *iotsitewise.BatchGetAssetPropertyAggregatesInput) bool {
			entries := *input.Entries[0]
			return *entries.EntryId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
				*entries.AssetId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
				*entries.PropertyId == "11propid-aaaa-2222-bbbb-3333cccc4444" &&
				*entries.AggregateTypes[0] == "SUM"
		}),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		NextToken: Pointer("some-next-token"),
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{{
			AggregatedValues: []*iotsitewise.AggregatedValue{{
				Timestamp: Pointer(time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)),
				Value:     &iotsitewise.Aggregates{Sum: Pointer(1688.6)},
			}},
			EntryId: aws.String("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}},
	}, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewise.Property{
			Name: Pointer("Wind Speed"),
		},
	}, nil)

	srvr := &server.Server{Datasource: mockedDatasource(mockSw).(*sitewise.Datasource)}

	sitewise.GetCache = func() *cache.Cache {
		return cache.New(cache.DefaultExpiration, cache.NoExpiration)
	}

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyAggregate,
				TimeRange: timeRange,
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
						 "propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
					   "aggregates":[
						  "SUM"
					   ],
					   "resolution":"1m"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)}),
		data.NewField("sum", nil, []float64{1688.6}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  "some-next-token",
			Resolution: "1m",
			Aggregates: []string{models.AggregateSum},
		},
	})
	if diff := cmp.Diff(expectedFrame, qdr.Responses["A"].Frames[0], data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
	mockSw.AssertCalled(t, "DescribeAssetPropertyWithContext", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
	})
}

func Test_property_value_aggregate_query_by_property_alias(t *testing.T) {
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer("/amazon/renton/1/rpm"),
		AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
	}, nil)
	mockSw.On(
		"BatchGetAssetPropertyAggregatesPageAggregation",
		mock.Anything,
		mock.MatchedBy(func(input *iotsitewise.BatchGetAssetPropertyAggregatesInput) bool {
			entries := *input.Entries[0]
			return *entries.EntryId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
				*entries.PropertyAlias == "/amazon/renton/1/rpm" &&
				*entries.AggregateTypes[0] == "SUM"
		}),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		NextToken: Pointer("some-next-token"),
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{{
			AggregatedValues: []*iotsitewise.AggregatedValue{{
				Timestamp: Pointer(time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)),
				Value:     &iotsitewise.Aggregates{Sum: Pointer(1688.6)},
			}},
			EntryId: aws.String("1assetid-aaaa-2222-bbbb-3333cccc4444"),
		}},
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

	qdr, err := srvr.HandlePropertyAggregate(context.Background(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyAggregate,
				TimeRange: timeRange,
				JSON: []byte(
					`{
					   "region":"us-west-2",
					   "propertyAlias":"/amazon/renton/1/rpm",
					   "aggregates":[
						  "SUM"
					   ],
					   "resolution":"1m"
					}`),
			},
		},
	})
	require.Nil(t, err)
	_, ok := qdr.Responses["A"]
	require.True(t, ok)
	require.NotNil(t, qdr.Responses["A"].Frames[0])

	expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)}),
		data.NewField("sum", nil, []float64{1688.6}),
	).SetMeta(&data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  "some-next-token",
			Resolution: "1m",
			Aggregates: []string{models.AggregateSum},
		},
	})
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
		"DescribeAssetPropertyWithContext",
		mock.Anything,
		&iotsitewise.DescribeAssetPropertyInput{
			AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
			PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
		},
	)
}

func Pointer[T any](v T) *T { return &v }
