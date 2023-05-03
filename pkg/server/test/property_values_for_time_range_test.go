package test

import (
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlePropertyValuesForTimeRange(t *testing.T) {
	var scenario = func(name string, baseQuery models.BaseQuery, expectedResolution string) *testScenario {

		baseQuery.QueryType = models.QueryTypePropertyAggregate
		query := models.AssetPropertyValueQuery{
			BaseQuery:      baseQuery,
			AggregateTypes: []string{"avg"},
			Resolution:     "AUTO",
		}

		mockSw := &mocks.SitewiseClient{}

		propVals := testdata.GetIoTSitewisePropHistoryVals(t, testDataRelativePath("property-history-values.json"))
		propAggregates := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
		propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-avg-wind.json"))
		propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))

		mockSw.On("BatchGetAssetPropertyValueHistoryPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propVals, nil)
		mockSw.On("BatchGetAssetPropertyAggregatesPageAggregation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&propAggregates, nil)
		mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
		mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

		return &testScenario{
			name: name,
			queries: []backend.DataQuery{
				{
					RefID:         "A",
					QueryType:     models.QueryTypePropertyAggregate,
					MaxDataPoints: baseQuery.MaxDataPoints,
					Interval:      baseQuery.Interval,
					TimeRange:     baseQuery.TimeRange,
					JSON:          testdata.SerializeStruct(t, query),
				},
			},
			mockSw:         mockSw,
			goldenFileName: "prop-val-for-time-range-" + strings.ReplaceAll(name, " ", "-"),
			handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
				return srvr.HandlePropertyAggregate
			},
			validationFn: func(t *testing.T, dr *backend.QueryDataResponse) {
				resp := dr.Responses["A"]
				frame := resp.Frames[0]
				actual, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
				assert.True(t, ok, "unable to cast custom metadata")
				assert.Equal(t, expectedResolution, actual.Resolution)
				if expectedResolution == "RAW" {
					assert.Equal(t, "raw", frame.Fields[1].Name)
				}
			},
		}
	}

	scenario("raw data for time range", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		AssetId:       testdata.DemoTurbineAsset1,
		PropertyId:    testdata.TurbinePropWindSpeed,
		TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionRaw).run(t)

	scenario("1m data for time range", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		AssetId:       testdata.DemoTurbineAsset1,
		PropertyId:    testdata.TurbinePropWindSpeed,
		TimeRange:     backend.TimeRange{From: testdata.TwoHours, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionMinute).run(t)

	scenario("1h data for time range", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		AssetId:       testdata.DemoTurbineAsset1,
		PropertyId:    testdata.TurbinePropWindSpeed,
		TimeRange:     backend.TimeRange{From: testdata.OneDay, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionHour).run(t)

	scenario("1d data for time range", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		AssetId:       testdata.DemoTurbineAsset1,
		PropertyId:    testdata.TurbinePropWindSpeed,
		TimeRange:     backend.TimeRange{From: testdata.OneMonth, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionDay).run(t)

	scenario("1m data for reduced max data point", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		AssetId:       testdata.DemoTurbineAsset1,
		PropertyId:    testdata.TurbinePropWindSpeed,
		TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
		MaxDataPoints: 299,
	}, propvals.ResolutionMinute).run(t)

	scenario("raw data for time range from alias", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionRaw).run(t)

	scenario("1m data for time range from alias", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		TimeRange:     backend.TimeRange{From: testdata.TwoHours, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionMinute).run(t)

	scenario("1h data for time range from alias", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		TimeRange:     backend.TimeRange{From: testdata.OneDay, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionHour).run(t)

	scenario("1d data for time range from alias", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		TimeRange:     backend.TimeRange{From: testdata.OneMonth, To: testdata.Now},
		MaxDataPoints: 720,
	}, propvals.ResolutionDay).run(t)

	scenario("1m data for reduced max data point from alias", models.BaseQuery{
		AwsRegion:     testdata.AwsRegion,
		PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
		MaxDataPoints: 299,
	}, propvals.ResolutionMinute).run(t)

}
