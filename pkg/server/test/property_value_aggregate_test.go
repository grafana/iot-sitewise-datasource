package test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/stretchr/testify/mock"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

func TestHandlePropertyValueAggregate(t *testing.T) {

	scenarios := []*testScenario{
		propertyValueAggregateHappyCase(t),
	}

	for _, s := range scenarios {
		s.run(t)
	}
}

var propertyValueAggregateHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	propAggs := testdata.GetIoTSitewisePropAggregateVals(t, testDataRelativePath("property-aggregate-values.json"))
	propDesc := testdata.GetIotSitewiseAssetProp(t, testDataRelativePath("describe-asset-property-raw-wind.json"))

	mockSw.On("GetAssetPropertyAggregatesWithContext", mock.Anything, mock.Anything).Return(&propAggs, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  "us-west-2",
			AssetId:    testdata.TestAssetId,
			PropertyId: testdata.TestPropIdRawWin},
		AggregateTypes: []string{models.AggregateStdDev, models.AggregateMin, models.AggregateAvg, models.AggregateCount, models.AggregateMax, models.AggregateSum},
		Resolution:     "1m",
	}

	return &testScenario{
		name: "PropertyValueAggregateHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypePropertyAggregate,
				TimeRange: timeRange,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "property-aggregate-values",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyAggregate
		},
		validationFn: nil,
	}
}
