package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

func Test_propertyValueAggregateHappyCase(t *testing.T) {
	propAggs := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggs, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

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
				RefID:     "A",
				QueryType: models.QueryTypePropertyAggregate,
				TimeRange: timeRange,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:  testdata.AwsRegion,
						AssetId:    testdata.DemoTurbineAsset1,
						PropertyId: testdata.TurbinePropWindSpeed},
					AggregateTypes: []string{models.AggregateStdDev, models.AggregateMin, models.AggregateAvg, models.AggregateCount, models.AggregateMax, models.AggregateSum},
					Resolution:     "1m",
				}),
			},
		},
	})
	require.Nil(t, err)

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "property-aggregate-values", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}
}

func Test_propertyValueAggregateFromAliasHappyCase(t *testing.T) {
	propAggs := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggs, nil)
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
				RefID:     "A",
				QueryType: models.QueryTypePropertyAggregate,
				TimeRange: timeRange,
				JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
					BaseQuery: models.BaseQuery{
						AwsRegion:     testdata.AwsRegion,
						PropertyAlias: "/amazon/renton/1/rpm",
					},
					AggregateTypes: []string{models.AggregateStdDev, models.AggregateMin, models.AggregateAvg, models.AggregateCount, models.AggregateMax, models.AggregateSum},
					Resolution:     "1m",
				}),
			},
		},
	})
	require.Nil(t, err)

	for i, dr := range qdr.Responses {
		fname := fmt.Sprintf("%s-%s.golden", "property-aggregate-values-from-alias", i)
		experimental.CheckGoldenJSONResponse(t, "../../testdata", fname, &dr, true)
	}

	mockSw.AssertExpectations(t)
}

func Test_propertyValueAggregateFromAliasHappyCase_wip(t *testing.T) {
	// extract minimum number of elements out of the json files ?
	// inline any of the input/output to the test
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer("/amazon/renton/1/rpm"),
		AssetId:    Pointer("e64c9075-9d89-47cb-8ee5-d3251bd253f4"),
		PropertyId: Pointer("3627f45d-710a-47c8-ae6c-4b71f7c9f5eb"),
	}, nil)
	//mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, &iotsitewise.DescribeTimeSeriesInput{
	//	Alias: Pointer("/amazon/renton/1/rpm")}).Return(&iotsitewise.DescribeTimeSeriesOutput{
	//	Alias:      Pointer("/amazon/renton/1/rpm"),
	//	AssetId:    Pointer("e64c9075-9d89-47cb-8ee5-d3251bd253f4"),
	//	PropertyId: Pointer("3627f45d-710a-47c8-ae6c-4b71f7c9f5eb"),
	//}, nil)

	//mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.MatchedBy(func(input *iotsitewise.DescribeTimeSeriesInput) bool { return input.Alias == Pointer("/amazon/renton/1/rpm") })).Return(&iotsitewise.DescribeTimeSeriesOutput{
	//	Alias:      Pointer("/amazon/renton/1/rpm"),
	//	AssetId:    Pointer("e64c9075-9d89-47cb-8ee5-d3251bd253f4"),
	//	PropertyId: Pointer("3627f45d-710a-47c8-ae6c-4b71f7c9f5eb"),
	//}, nil)
	mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		NextToken: Pointer("some-next-token"),
		SuccessEntries: []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{{
			AggregatedValues: []*iotsitewise.AggregatedValue{{
				Timestamp: Pointer(time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)),
				Value: &iotsitewise.Aggregates{
					Average:           Pointer(21.5),
					Count:             Pointer(float64(60)),
					Maximum:           Pointer(22.5),
					Minimum:           Pointer(20.5),
					StandardDeviation: Pointer(0.06),
					Sum:               Pointer(1688.6),
				},
			}},
			EntryId: aws.String("e64c9075-9d89-47cb-8ee5-d3251bd253f4"),
		}},
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

	expectedFrame := data.NewFrame(" /amazon/renton/1/rpm",
		data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)}),
		data.NewField("avg", nil, []float64{21.5}),
		data.NewField("min", nil, []float64{20.5}),
		data.NewField("max", nil, []float64{22.5}),
		data.NewField("sum", nil, []float64{1688.6}),
		data.NewField("count", nil, []float64{60}),
		data.NewField("stddev", nil, []float64{0.06}),
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
	mockSw.AssertCalled(t, "DescribeTimeSeriesWithContext", mock.Anything, &iotsitewise.DescribeTimeSeriesInput{
		Alias: Pointer("/amazon/renton/1/rpm")}) //other way to do API call input assertions
}

func Pointer[T any](v T) *T { return &v }
