package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func TestHandlePropertyValue(t *testing.T) {
	getPropertyValueHappyCase(t).run(t)
	getPropertyValueFromPropertyAliasCase(t).run(t)
	getPropertyValueEmptyCase(t).run(t)
}

var getPropertyValueHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVal := testdata.GetIoTSitewisePropVal(t, testDataRelativePath("property-value.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))

	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&propVal, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  testdata.AwsRegion,
			AssetId:    testdata.DemoTurbineAsset1,
			PropertyId: testdata.TurbinePropWindSpeed,
		},
	}

	return &testScenario{
		name: "GetPropertyValueHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "property-value",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValue
		},
		validationFn: nil,
	}
}

var getPropertyValueFromPropertyAliasCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVal := testdata.GetIoTSitewisePropVal(t, testDataRelativePath("property-value.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))
	propTimeSeries := testdata.GetIoTSitewiseTimeSeries(t, testDataRelativePath("describe-time-series.json"))

	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&propVal, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)
	mockSw.On("DescribeTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&propTimeSeries, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:     testdata.AwsRegion,
			PropertyAlias: testdata.TurbinePropWindSpeedAlias,
		},
	}

	return &testScenario{
		name: "GetPropertyValueFromAliasCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "property-value-from-alias",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValue
		},
		validationFn: nil,
	}
}

var getPropertyValueEmptyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVal := iotsitewise.BatchGetAssetPropertyValueOutput{SuccessEntries: []*iotsitewise.BatchGetAssetPropertyValueSuccessEntry{{
		AssetPropertyValue: nil,
		EntryId:            aws.String(testdata.DemoTurbineAsset1),
	}}} // empty prop value response
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))

	mockSw.On("BatchGetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&propVal, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  testdata.AwsRegion,
			AssetId:    testdata.DemoTurbineAsset1,
			PropertyId: testdata.TurbinePropWindSpeed,
		},
	}

	return &testScenario{
		name: "GetPropertyValueHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyValue,
				TimeRange: timeRange,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "property-value-empty",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValue
		},
		validationFn: nil,
	}
}
