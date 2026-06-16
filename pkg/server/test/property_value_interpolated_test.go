package test

import (
	"context"
	"fmt"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"math"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/google/go-cmp/cmp"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mockGetInterpolatedAssetPropertyValuesPageAggregation(mockSw *mocks.SitewiseAPIClient, nextToken *string, value *float64, matchFn interface{}) {
	mockSw.On(
		"GetInterpolatedAssetPropertyValuesPageAggregation",
		mock.Anything,
		mock.MatchedBy(matchFn),
		mock.Anything,
		mock.Anything,
	).Return(&iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		NextToken: nextToken,
		InterpolatedAssetPropertyValues: []iotsitewisetypes.InterpolatedAssetPropertyValue{
			{
				Timestamp: &iotsitewisetypes.TimeInNanos{
					OffsetInNanos: Pointer(int32(0)),
					TimeInSeconds: Pointer(int64(1612207200)),
				},
				Value: &iotsitewisetypes.Variant{
					DoubleValue: value,
				},
			},
		},
	}, nil)
}

func TestPropertyValueInterpolatedQueryWithPropertyAlias(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}

	mockGetInterpolatedAssetPropertyValuesPageAggregation(mockSw, Pointer("asset1-next-token"), Pointer(1.1), func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
		return *input.AssetId == mockAssetId
	})

	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias:      Pointer(mockPropertyAlias),
		AssetId:    Pointer(mockAssetId),
		PropertyId: Pointer(mockPropertyId),
	}, nil)

	mockSw.On("DescribeAssetProperty", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
		AssetId:    aws.String(mockAssetId),
		PropertyId: aws.String(mockPropertyId),
	}, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
		AssetId:   Pointer(mockAssetId),
		AssetName: Pointer("Demo Turbine Asset 1"),
		AssetProperty: &iotsitewisetypes.Property{
			DataType: iotsitewisetypes.PropertyDataTypeDouble,
			Name:     Pointer("Wind Speed"),
			Unit:     Pointer("m/s"),
			Id:       aws.String(mockPropertyId),
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
				JSON: []byte(fmt.Sprintf(`{
					"propertyAlias": "%s",
					"resolution": "1m"
				}`, mockPropertyAlias)),
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
		data.NewFrame("Demo Turbine Asset 1",
			data.NewField("time", nil, []time.Time{time.Unix(1612207200, 0)}),
			data.NewField("Wind Speed", nil, []float64{1.1}).SetConfig(&data.FieldConfig{Unit: "m/s"}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "asset1-next-token",
				EntryId:    *mockAssetPropertyEntryId,
				Resolution: "1m",
			},
		}),
	}

	if diff := cmp.Diff(expectedFrames, qdr.Responses["A"].Frames, data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func TestPropertyValueInterpolatedQueryWithDisassociatedPropertyAlias(t *testing.T) {
	mockSw := &mocks.SitewiseAPIClient{}

	mockGetInterpolatedAssetPropertyValuesPageAggregation(mockSw, Pointer("asset1-next-token"), Pointer(1.1), func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
		return *input.PropertyAlias == mockPropertyAlias
	})

	mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
		Alias: Pointer(mockPropertyAlias),
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
				JSON: []byte(fmt.Sprintf(`{
					"propertyAlias": "%s",
					"resolution": "1m"
				}`, mockPropertyAlias)),
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
		data.NewFrame(mockPropertyAlias,
			data.NewField("time", nil, []time.Time{time.Unix(1612207200, 0)}),
			data.NewField(mockPropertyAlias, nil, []float64{1.1}).SetConfig(&data.FieldConfig{}),
		).SetMeta(&data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  "asset1-next-token",
				EntryId:    *mockPropertyAliasEntryId,
				Resolution: "1m",
			},
		}),
	}

	if diff := cmp.Diff(expectedFrames, qdr.Responses["A"].Frames, data.FrameTestCompareOptions()...); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

	mockSw.AssertExpectations(t)
}

func TestPropertyValueInterpolatedQueryMultipleIds(t *testing.T) {
	tests := []batch_test{
		{
			name:           "query by multiple assetIds and one propertyId",
			numAssetIds:    10,
			numPropertyIds: 1,
		},
		{
			name:           "query by one assetId and multiple propertyIds",
			numAssetIds:    1,
			numPropertyIds: 10,
		},
		{
			name:               "query by multiple property aliases",
			numPropertyAliases: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSw := &mocks.SitewiseAPIClient{}

			if tc.numPropertyAliases > 0 {
				propertyAliases := generateIds(tc.numPropertyAliases, mockPropertyAlias)
				for p, propertyAlias := range propertyAliases {
					nextToken := Pointer(fmt.Sprintf("next-token-%d", p))
					doubleValue := Pointer(1.1 + float64(p))
					mockGetInterpolatedAssetPropertyValuesPageAggregation(mockSw, nextToken, doubleValue, func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
						return *input.PropertyAlias == propertyAlias
					})
					mockSw.On("DescribeTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.DescribeTimeSeriesOutput{
						Alias: Pointer(propertyAlias),
					}, nil)
				}
			} else {
				assetIds := generateIds(tc.numAssetIds, mockAssetId)
				propertyIds := generateIds(tc.numPropertyIds, mockPropertyId)
				for a, assetId := range assetIds {
					for p, propertyId := range propertyIds {
						diff := a*tc.numPropertyIds + p
						nextToken := Pointer(fmt.Sprintf("next-token-%d", diff))
						doubleValue := Pointer(1.1 + float64(diff))
						mockGetInterpolatedAssetPropertyValuesPageAggregation(mockSw, nextToken, doubleValue, func(input *iotsitewise.GetInterpolatedAssetPropertyValuesInput) bool {
							return *input.AssetId == assetId && *input.PropertyId == propertyId
						})
						mockSw.On("DescribeAssetProperty", mock.Anything, &iotsitewise.DescribeAssetPropertyInput{
							AssetId:    aws.String(assetId),
							PropertyId: aws.String(propertyId),
						}, mock.Anything).Return(&iotsitewise.DescribeAssetPropertyOutput{
							AssetId:   Pointer(assetId),
							AssetName: Pointer(fmt.Sprintf("Demo Turbine Asset %d", a)),
							AssetProperty: &iotsitewisetypes.Property{
								DataType: iotsitewisetypes.PropertyDataTypeDouble,
								Name:     Pointer("Wind Speed"),
								Unit:     Pointer("m/s"),
								Id:       aws.String(propertyId),
							},
						}, nil)
					}
				}
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

			qdr, err := srvr.HandleInterpolatedPropertyValue(context.Background(), &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{},
				Queries: []backend.DataQuery{{
					QueryType: models.QueryTypePropertyInterpolated,
					RefID:     "A",
					TimeRange: timeRange,
					JSON: testdata.SerializeStruct(t, models.AssetPropertyValueQuery{
						BaseQuery:  baseQuery,
						Resolution: "1m",
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

			sort.Slice(qdr.Responses["A"].Frames, func(a, b int) bool {
				return qdr.Responses["A"].Frames[a].Meta.Custom.(models.SitewiseCustomMeta).NextToken < qdr.Responses["A"].Frames[b].Meta.Custom.(models.SitewiseCustomMeta).NextToken
			})

			expectedFrames := data.Frames{}
			for i, f := range qdr.Responses["A"].Frames {
				require.NotNil(t, f)
				expectedNextToken := fmt.Sprintf("next-token-%d", i)
				var entryId string
				var frameName string
				var propertyField *data.Field
				if tc.numPropertyAliases > 0 {
					propertyAlias := fmt.Sprintf("%s%d", mockPropertyAlias, i+1)
					entryId = *util.GetEntryIdFromPropertyAlias(propertyAlias)
					frameName = propertyAlias
					propertyField = data.NewField(propertyAlias, nil, []float64{1.1 + float64(i)}).SetConfig(&data.FieldConfig{})
				} else {
					assetIdx := int(math.Floor(float64(i) / float64(tc.numPropertyIds)))
					assetId := fmt.Sprintf("%s%d", mockAssetId, assetIdx+1)
					propertyId := fmt.Sprintf("%s%d", mockPropertyId, i%tc.numPropertyIds+1)
					entryId = *util.GetEntryIdFromAssetProperty(assetId, propertyId)
					frameName = fmt.Sprintf("Demo Turbine Asset %d", assetIdx)
					propertyField = data.NewField("Wind Speed", nil, []float64{1.1 + float64(i)}).SetConfig(&data.FieldConfig{Unit: "m/s"})
				}
				expectedFrames = append(expectedFrames, data.NewFrame(frameName,
					data.NewField("time", nil, []time.Time{time.Unix(1612207200, 0)}),
					propertyField,
				).SetMeta(&data.FrameMeta{
					Custom: models.SitewiseCustomMeta{
						NextToken:  expectedNextToken,
						EntryId:    entryId,
						Resolution: "1m",
					},
				}))
				// require.Equal(t, entryId, f.Meta.Custom.(models.SitewiseCustomMeta).EntryId)
				// require.Equal(t, expectedNextToken, f.Meta.Custom.(models.SitewiseCustomMeta).NextToken)
			}

			if diff := cmp.Diff(expectedFrames, qdr.Responses["A"].Frames, data.FrameTestCompareOptions()...); diff != "" {
				t.Errorf("Result mismatch (-want +got):\n%s", diff)
			}

			mockSw.AssertExpectations(t)
		})
	}
}
