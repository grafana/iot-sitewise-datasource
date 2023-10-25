package test

import (
	"context"
	"math"
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

type test struct {
	name                                      string
	query                                     string
	isExpression                              bool
	expectedMaxPages                          int
	expectedMaxResults                        int
	expectedDescribeTimeSeriesWithContextArgs *iotsitewise.DescribeTimeSeriesInput
}

func TestPropertyValueAggregate(t *testing.T) {
	tests := []test{
		{
			name: "query by asset id and property id",
			query: `{
				"region":"us-west-2",
				"assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
				"propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`,
			expectedMaxPages:   1,
			expectedMaxResults: 0,
		},
		{
			name:         "expression query by asset id and property",
			isExpression: true,
			query: `{
				"region":"us-west-2",
				"assetId":"1assetid-aaaa-2222-bbbb-3333cccc4444",
				"propertyId":"11propid-aaaa-2222-bbbb-3333cccc4444",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`,
			expectedMaxPages:   math.MaxInt32,
			expectedMaxResults: math.MaxInt32,
		},
		{
			name: "query by property alias",
			query: `{
				"region":"us-west-2",
				"propertyAlias":"/amazon/renton/1/rpm",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`,
			expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
			expectedMaxPages:   1,
			expectedMaxResults: 0,
		},
		{
			name:         "expression query by property alias",
			isExpression: true,
			query: `{
				"region":"us-west-2",
				"propertyAlias":"/amazon/renton/1/rpm",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`,
			expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
			expectedMaxPages:   math.MaxInt32,
			expectedMaxResults: math.MaxInt32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseClient{}

			if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
				mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
					Alias:      Pointer("/amazon/renton/1/rpm"),
					AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
					PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
				}, nil)
			}

			mockSw.On(
				"BatchGetAssetPropertyAggregatesPageAggregation",
				mock.Anything,
				mock.MatchedBy(func(input *iotsitewise.BatchGetAssetPropertyAggregatesInput) bool {
					entries := *input.Entries[0]

					if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
						return *entries.EntryId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
							*entries.PropertyAlias == "/amazon/renton/1/rpm" &&
							*entries.AggregateTypes[0] == "SUM"
					} else {
						return *entries.EntryId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
							*entries.AssetId == "1assetid-aaaa-2222-bbbb-3333cccc4444" &&
							*entries.PropertyId == "11propid-aaaa-2222-bbbb-3333cccc4444" &&
							*entries.AggregateTypes[0] == "SUM"
					}
				}),
				tc.expectedMaxPages,
				tc.expectedMaxResults,
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

			query := &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{},
				Queries: []backend.DataQuery{
					{
						RefID:     "A",
						QueryType: models.QueryTypePropertyAggregate,
						TimeRange: timeRange,
						JSON:      []byte(tc.query),
					},
				},
			}

			if tc.isExpression {
				query.Headers = map[string]string{"http_X-Grafana-From-Expr": "true"}
			}

			qdr, err := srvr.HandlePropertyAggregate(context.Background(), query)
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
			if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
				mockSw.AssertCalled(t,
					"DescribeTimeSeriesWithContext",
					mock.Anything,
					tc.expectedDescribeTimeSeriesWithContextArgs,
				)
			}
			mockSw.AssertCalled(t, "DescribeAssetPropertyWithContext", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
				AssetId:    Pointer("1assetid-aaaa-2222-bbbb-3333cccc4444"),
				PropertyId: Pointer("11propid-aaaa-2222-bbbb-3333cccc4444"),
			})
		})
	}
}

func TestPropertyValueAggregateWithDisassociatedStream(t *testing.T) {
	tc := test{
		// an disassociated stream will return nil in DescribeTimeSeriesWithContext for assetId and propertyId
		name: "query by property alias of an disassociated stream",
		query: `{
					"region":"us-west-2",
					"propertyAlias":"/amazon/renton/1/rpm",
					"aggregates":["SUM"],
					"resolution":"1m"
				}`,
		expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer("/amazon/renton/1/rpm")},
		expectedMaxPages:   1,
		expectedMaxResults: 0,
	}

	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseClient{}

		if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
			alias := Pointer("/amazon/renton/1/rpm")
			var assetId *string
			var propertyId *string

			mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
				Alias:      alias,
				AssetId:    assetId,
				PropertyId: propertyId,
			}, nil)
		}
		mockSw.On(
			"BatchGetAssetPropertyAggregatesPageAggregation",
			mock.Anything,
			mock.MatchedBy(func(input *iotsitewise.BatchGetAssetPropertyAggregatesInput) bool {
				entries := *input.Entries[0]
				return *entries.EntryId == "61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef" &&
					*entries.PropertyAlias == "/amazon/renton/1/rpm" &&
					*entries.AggregateTypes[0] == "SUM"

			}),
			tc.expectedMaxPages,
			tc.expectedMaxResults,
		).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
			NextToken: Pointer("some-next-token"),
			SuccessEntries: []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{{
				AggregatedValues: []*iotsitewise.AggregatedValue{{
					Timestamp: Pointer(time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)),
					Value:     &iotsitewise.Aggregates{Sum: Pointer(1688.6)},
				}},
				EntryId: aws.String("61e4e1a8ab39463fa0b9418d9be2923e364f40a8b935b69d006b999516cdecef"),
			}},
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
					QueryType: models.QueryTypePropertyAggregate,
					TimeRange: timeRange,
					JSON:      []byte(tc.query),
				},
			},
		}

		if tc.isExpression {
			query.Headers = map[string]string{"http_X-Grafana-From-Expr": "true"}
		}

		qdr, err := srvr.HandlePropertyAggregate(context.Background(), query)
		require.Nil(t, err)
		_, ok := qdr.Responses["A"]
		require.True(t, ok)
		require.NotNil(t, qdr.Responses["A"].Frames[0])

		expectedFrame := data.NewFrame("/amazon/renton/1/rpm",
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
		if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
			mockSw.AssertCalled(t,
				"DescribeTimeSeriesWithContext",
				mock.Anything,
				tc.expectedDescribeTimeSeriesWithContextArgs,
			)
		}
		mockSw.AssertNotCalled(t, "DescribeAssetPropertyWithContext", mock.Anything, mock.Anything)

	})

}

func Pointer[T any](v T) *T { return &v }
