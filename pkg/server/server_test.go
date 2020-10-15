package server

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/testutil"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

var (
	timeRange = backend.TimeRange{
		From: time.Now().Add(time.Hour * -3),
		To:   time.Now(),
	}
)

type testScenario struct {
	name         string
	queries      []backend.DataQuery
	propVals     iotsitewise.GetAssetPropertyValueHistoryOutput
	property     iotsitewise.DescribeAssetPropertyOutput
	handlerFn    func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc
	validationFn func(t *testing.T, dr *backend.QueryDataResponse, err error)
}

var propertyValueHistoryResponseScenario = func(t *testing.T) *testScenario {

	query := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion: "us-west-2",
		},
		AssetId:    testutil.TestAssetId,
		PropertyId: testutil.TestPropIdAvgWind,
	}

	qbytes, err := json.Marshal(query)
	if err != nil {
		t.Fatal(err)
	}

	return &testScenario{
		name: "PropertyValueHistoryResponseHappyCase",
		queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON:          qbytes,
			},
		},
		propVals: testutil.GetIoTSitewisePropHistoryVals(t, "property-history-values.json"),
		property: testutil.GetIotSitewiseAssetProp(t, "describe-asset-property-avg-wind.json"),
		handlerFn: func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
		validationFn: func(t *testing.T, qdr *backend.QueryDataResponse, err error) {

			assert.NoError(t, err)
			assert.Len(t, qdr.Responses, 1)

			dr, found := qdr.Responses["A"]

			assert.True(t, found, "could not find expected data response")
			assert.NoError(t, dr.Error)
			assert.Len(t, dr.Frames, 1)

			// does it have the expected asset property
			assert.Equal(t, dr.Frames[0].Name, testutil.TestPropertyName)
			// are there the expected number of fields
			assert.Len(t, dr.Frames[0].Fields, 2)
			// do both fields have data
			assert.True(t, dr.Frames[0].Fields[0].Len() > 1)
			assert.True(t, dr.Frames[0].Fields[1].Len() > 1)
		},
	}

}

func testScenarios(t *testing.T) []*testScenario {
	return []*testScenario{
		propertyValueHistoryResponseScenario(t),
	}
}

func mockedDatasource(swmock *mocks.Client) Datasource {
	return &sitewise.Datasource{
		GetClient: func(_ backend.PluginContext, _ models.BaseQuery) (client client.Client, err error) {
			client = swmock
			return
		},
	}
}

func runTestScenario(t *testing.T, scenario *testScenario) {

	t.Run(scenario.name, func(t *testing.T) {

		ctx := context.Background()
		swmock := &mocks.Client{}

		swmock.On("GetAssetPropertyValueHistoryWithContext", mock.Anything, mock.Anything).Return(&scenario.propVals, nil)
		swmock.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&scenario.property, nil)

		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{},
			Queries:       scenario.queries,
		}

		srvr := &Server{
			datasource: mockedDatasource(swmock),
		}

		dr, err := scenario.handlerFn(t, srvr)(ctx, req)

		scenario.validationFn(t, dr, err)

	})

}

func TestDataResponse(t *testing.T) {

	for _, v := range testScenarios(t) {
		runTestScenario(t, v)
	}

}
