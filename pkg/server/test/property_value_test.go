package test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func TestHandlePropertyValue(t *testing.T) {
	getPropertyValueHappyCase(t).run(t)
}

var getPropertyValueHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	propVal := testdata.GetIoTSitewisePropVal(t, testDataRelativePath("property-value.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))

	mockSw.On("GetAssetPropertyValueWithContext", mock.Anything, mock.Anything).Return(&propVal, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  "us-west-2",
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
