package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
	"github.com/patrickmn/go-cache"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

type property_value_aggregate_test struct {
	name                                      string
	query                                     string
	isExpression                              bool
	expectedMaxPages                          int
	expectedMaxResults                        int
	expectedDescribeTimeSeriesWithContextArgs *iotsitewise.DescribeTimeSeriesInput
	numAssetIds                               int
	numPropertyIds                            int
	numPropertyAliases                        int
}

func mockBatchGetAssetPropertyAggregatesSuccessEntry(entryId *string, idx int) iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry {
	return iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{
		AggregatedValues: []*iotsitewise.AggregatedValue{{
			Timestamp: Pointer(time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)),
			Value:     &iotsitewise.Aggregates{Sum: Pointer(1688.6 + float64(idx))},
		}},
		EntryId: entryId,
	}
}

func mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw *mocks.SitewiseClient, nextToken *string, successEntries []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry, errorEntries []*iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry) {
	mockSw.On(
		"BatchGetAssetPropertyAggregatesPageAggregation",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		NextToken:      nextToken,
		SuccessEntries: successEntries,
		ErrorEntries:   errorEntries,
	}, nil).Once()
}

func TestPropertyValueAggregate(t *testing.T) {
	tests := []property_value_aggregate_test{
		{
			name: "query by asset id and property id",
			query: fmt.Sprintf(`{
				"region":"us-west-2",
				"assetId":"%s",
				"propertyId":"%s",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`, mockAssetId, mockPropertyId),
			expectedMaxPages:   1,
			expectedMaxResults: 0,
		},
		{
			name:         "expression query by asset id and property",
			isExpression: true,
			query: fmt.Sprintf(`{
				"region":"us-west-2",
				"assetId":"%s",
				"propertyId":"%s",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`, mockAssetId, mockPropertyId),
			expectedMaxPages:   math.MaxInt32,
			expectedMaxResults: math.MaxInt32,
		},
		{
			name: "query by property alias",
			query: fmt.Sprintf(`{
				"region":"us-west-2",
				"propertyAlias":"%s",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`, mockPropertyAlias),
			expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
			expectedMaxPages:   1,
			expectedMaxResults: 0,
		},
		{
			name:         "expression query by property alias",
			isExpression: true,
			query: fmt.Sprintf(`{
				"region":"us-west-2",
				"propertyAlias":"%s",
				"aggregates":["SUM"],
				"resolution":"1m"
			}`, mockPropertyAlias),
			expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
			expectedMaxPages:   math.MaxInt32,
			expectedMaxResults: math.MaxInt32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseClient{}

			mockDescribeAssetPropertyWithContext(mockSw)
			successEntry := mockBatchGetAssetPropertyAggregatesSuccessEntry(mockAssetPropertyEntryId, 0)
			mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw, Pointer("some-next-token"), []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{&successEntry}, nil)

			if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
				mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
					Alias:      Pointer(mockPropertyAlias),
					AssetId:    Pointer(mockAssetId),
					PropertyId: Pointer(mockPropertyId),
				}, nil)
			}

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
					EntryId:    *mockAssetPropertyEntryId,
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
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			})
		})
	}
}

func TestPropertyValueAggregateWithDisassociatedStream(t *testing.T) {
	tc := property_value_aggregate_test{
		// an disassociated stream will return nil in DescribeTimeSeriesWithContext for assetId and propertyId
		name: "query by property alias of an disassociated stream",
		query: fmt.Sprintf(`{
					"region":"us-west-2",
					"propertyAlias":"%s",
					"aggregates":["SUM"],
					"resolution":"1m"
				}`, mockPropertyAlias),
		expectedDescribeTimeSeriesWithContextArgs: &iotsitewise.DescribeTimeSeriesInput{Alias: Pointer(mockPropertyAlias)},
		expectedMaxPages:   1,
		expectedMaxResults: 0,
	}

	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseClient{}

		if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
			alias := Pointer(mockPropertyAlias)
			var assetId *string
			var propertyId *string

			mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
				Alias:      alias,
				AssetId:    assetId,
				PropertyId: propertyId,
			}, nil)
		}

		successEntry := mockBatchGetAssetPropertyAggregatesSuccessEntry(mockPropertyAliasEntryId, 0)
		mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw, Pointer("some-next-token"), []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{&successEntry}, nil)

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

		expectedFrame := data.NewFrame(mockPropertyAlias,
			data.NewField("time", nil, []time.Time{time.Date(2021, 2, 1, 16, 27, 0, 0, time.UTC)}),
			data.NewField("sum", nil, []float64{1688.6}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "some-next-token",
				EntryId:    *mockPropertyAliasEntryId,
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

func TestPropertyValueAggregate_with_error(t *testing.T) {
	tc := property_value_aggregate_test{
		name: "query by asset id and property id",
		query: fmt.Sprintf(`{
			"region":"us-west-2",
			"assetId":"%s",
			"propertyId":"%s",
			"aggregates":["SUM"],
			"resolution":"1m"
		}`, mockAssetId, mockPropertyId),
		expectedMaxPages:   1,
		expectedMaxResults: 0,
	}
	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseClient{}

		if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
			mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
				Alias:      Pointer(mockPropertyAlias),
				AssetId:    Pointer(mockAssetId),
				PropertyId: Pointer(mockPropertyId),
			}, nil)
		}

		mockSw.On(
			"BatchGetAssetPropertyAggregatesPageAggregation",
			mock.Anything,
			mock.MatchedBy(func(input *iotsitewise.BatchGetAssetPropertyAggregatesInput) bool {
				entries := *input.Entries[0]

				if tc.expectedDescribeTimeSeriesWithContextArgs != nil {
					return *entries.EntryId == *mockAssetPropertyEntryId &&
						*entries.PropertyAlias == mockPropertyAlias &&
						*entries.AggregateTypes[0] == "SUM"
				} else {
					return *entries.EntryId == *mockAssetPropertyEntryId &&
						*entries.AssetId == mockAssetId &&
						*entries.PropertyId == mockPropertyId &&
						*entries.AggregateTypes[0] == "SUM"
				}
			}),
			tc.expectedMaxPages,
			tc.expectedMaxResults,
		).Return(&iotsitewise.BatchGetAssetPropertyAggregatesOutput{
			NextToken: Pointer("some-next-token"),
			ErrorEntries: []*iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry{{
				ErrorCode:    Pointer("404"),
				ErrorMessage: Pointer("Asset property not found."),
				EntryId:      mockAssetPropertyEntryId,
			}},
		}, nil)

		mockDescribeAssetPropertyWithContext(mockSw)

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

		expectedFrame := data.NewFrame("Demo Turbine Asset 1 Wind Speed").SetMeta(&data.FrameMeta{
			Notices: []data.Notice{{Severity: data.NoticeSeverityError, Text: "Asset property not found."}},
		},
		)
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
			AssetId:    Pointer(mockAssetId),
			PropertyId: Pointer(mockPropertyId),
		})

	})

}

func TestPropertyValueAggregate_batched(t *testing.T) {
	tests := []property_value_aggregate_test{
		{
			name:           "query by multiple assetIds and one propertyId",
			numAssetIds:    api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
			numPropertyIds: 1,
		},
		{
			name:           "query by one assetId and multiple propertyIds",
			numAssetIds:    1,
			numPropertyIds: api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
		},
		{
			name:           "query by multiple assetIds and multiple propertyIds",
			numAssetIds:    api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
			numPropertyIds: api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
		},
		{
			name:               "query by multiple property aliases",
			numPropertyAliases: api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseClient{}
			mockedSuccessEntries := []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{}
			numBatch := 0

			if tc.numPropertyAliases > 0 {
				propertyAliases := generateIds(tc.numPropertyAliases, mockPropertyAlias)
				for p, propertyAlias := range propertyAliases {
					// Build the success entry based on the propertyAlias for disassociated data streams
					entryId := util.GetEntryIdFromPropertyAlias(propertyAlias)
					successEntry := mockBatchGetAssetPropertyAggregatesSuccessEntry(entryId, p)
					mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)

					isLastBatch := p == tc.numPropertyAliases-1
					// When batch is complete mock the History call with the success entries
					if len(mockedSuccessEntries) == api.BatchGetAssetPropertyAggregatesMaxEntries || isLastBatch {
						numBatch++
						mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
						mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{}
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
						successEntry := mockBatchGetAssetPropertyAggregatesSuccessEntry(entryId, p)
						mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)

						isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
						// When batch is complete mock the History call with the success entries
						if len(mockedSuccessEntries) == api.BatchGetAssetPropertyAggregatesMaxEntries || isLastBatch {
							numBatch++
							mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, nil)
							mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{}
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

			qdr, err := srvr.HandlePropertyAggregate(context.Background(), query)
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
		})
	}
}

func TestPropertyValueAggregate_batched_with_error(t *testing.T) {
	tc := property_value_aggregate_test{
		name:           "batch aggregate query with one error",
		numAssetIds:    api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
		numPropertyIds: api.BatchGetAssetPropertyAggregatesMaxEntries + 1,
	}

	t.Run(tc.name, func(t *testing.T) {
		mockSw := &mocks.SitewiseClient{}

		mockedSuccessEntries := []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{}
		mockedErrorEntries := []*iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry{}
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
					mockedErrorEntries = append(mockedErrorEntries, &iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry{
						ErrorCode:    Pointer("404"),
						ErrorMessage: Pointer("Asset property not found."),
						EntryId:      entryId,
					})
				} else {
					successEntry := mockBatchGetAssetPropertyAggregatesSuccessEntry(entryId, p)
					mockedSuccessEntries = append(mockedSuccessEntries, &successEntry)
				}

				isLastBatch := a == tc.numAssetIds-1 && p == tc.numPropertyIds-1
				// When batch is complete mock the History call with the success entries
				if len(mockedSuccessEntries)+len(mockedErrorEntries) == api.BatchGetAssetPropertyAggregatesMaxEntries || isLastBatch {
					numBatch++
					mockBatchGetAssetPropertyAggregatesPageAggregation(mockSw, Pointer(fmt.Sprintf("some-next-token-%d", numBatch)), mockedSuccessEntries, mockedErrorEntries)
					// Reset for next batch
					mockedSuccessEntries = []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry{}
					mockedErrorEntries = []*iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry{}
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

		qdr, err := srvr.HandlePropertyAggregate(context.Background(), query)
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
