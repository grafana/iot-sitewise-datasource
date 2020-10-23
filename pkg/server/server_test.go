package server

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"

	"github.com/stretchr/testify/mock"

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

type testServerScenarioFn func(t *testing.T) *testScenario

type testScenario struct {
	name           string
	queries        []backend.DataQuery
	mockSw         *mocks.Client
	goldenFileName string
	handlerFn      func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc
	validationFn   func(t *testing.T, dr *backend.QueryDataResponse)
}

func (ts *testScenario) run(t *testing.T) {
	runTestScenario(t, ts)
}

var getPropertyValueHistoryHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	propVals := testutil.GetIoTSitewisePropHistoryVals(t, "property-history-values.json")
	propDesc := testutil.GetIotSitewiseAssetProp(t, "describe-asset-property-avg-wind.json")

	mockSw.On("GetAssetPropertyValueHistoryWithContext", mock.Anything, mock.Anything).Return(&propVals, nil)
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&propDesc, nil)

	swQuery := models.AssetPropertyValueQuery{
		BaseQuery: models.BaseQuery{
			AwsRegion:  "us-west-2",
			AssetId:    testutil.TestAssetId,
			PropertyId: testutil.TestPropIdAvgWind,
		},
	}

	return &testScenario{
		name:   "PropertyValueHistoryResponseHappyCase",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				QueryType:     models.QueryTypePropertyValueHistory,
				RefID:         "A",
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     timeRange,
				JSON:          testutil.SerializeStruct(t, swQuery),
			},
		},
		goldenFileName: "property-history-values",
		handlerFn: func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc {
			return srvr.HandlePropertyValueHistory
		},
	}
}

var listAssetModelsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	assetModels := testutil.GetIoTSitewiseAssetModels(t, "list-asset-models.json")

	mockSw.On("ListAssetModelsWithContext", mock.Anything, mock.Anything).Return(&assetModels, nil)

	query := models.ListAssetModelsQuery{
		BaseQuery: models.BaseQuery{},
		NextToken: "",
	}

	return &testScenario{
		name:   "TestListAssetModelsResponseHappyCase",
		mockSw: mockSw,
		queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     models.QueryTypeListAssetModels,
				MaxDataPoints: 100,
				Interval:      1000,
				TimeRange:     backend.TimeRange{},
				JSON:          testutil.SerializeStruct(t, query),
			},
		},
		goldenFileName: "list-asset-models",
		handlerFn: func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssetModels
		},
	}
}

var listAssetsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	topLevelAssets := testutil.GetIoTSitewiseAssets(t, "list-assets-top-level.json")
	childAssets := testutil.GetIoTSitewiseAssets(t, "list-assets.json")

	mockSw.On("ListAssetsWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		return req.AssetModelId == nil && *req.Filter == "TOP_LEVEL"
	})).Return(&topLevelAssets, nil)

	mockSw.On("ListAssetsWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.ListAssetsInput) bool {
		if req.AssetModelId == nil {
			return false
		}
		return *req.AssetModelId == testutil.TestAssetModelId && *req.Filter == "ALL"
	})).Return(&childAssets, nil)

	queryTopLevel := models.ListAssetsQuery{
		ModelId: "",
		Filter:  "",
	}

	queryChild := models.ListAssetsQuery{
		ModelId: testutil.TestAssetModelId,
		Filter:  "ALL",
	}

	return &testScenario{
		name: "TestListAssetsHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssets,
				JSON:      testutil.SerializeStruct(t, queryTopLevel),
			},
			{
				RefID:     "B",
				QueryType: models.QueryTypeListAssets,
				JSON:      testutil.SerializeStruct(t, queryChild),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "list-assets",
		handlerFn: func(t *testing.T, srvr *Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssets
		},
		validationFn: nil,
	}
}

var describeAssetHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	asset := testutil.GetIoTSitewiseAssetDescription(t, "describe-asset.json")
	assetWithHierarchy := testutil.GetIoTSitewiseAssetDescription(t, "describe-asset-top-level.json")

	mockSw.On("DescribeAssetWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetInput) bool {
		return req.AssetId != nil && *req.AssetId == testutil.TestAssetId
	})).Return(&asset, nil)

	mockSw.On("DescribeAssetWithContext", mock.Anything, mock.MatchedBy(func(req *iotsitewise.DescribeAssetInput) bool {
		return req.AssetId != nil && *req.AssetId == testutil.TestTopLevelAssetId
	})).Return(&assetWithHierarchy, nil)

	query := models.DescribeAssetQuery{}
	query.AssetId = testutil.TestAssetId

	queryTopLevel := models.DescribeAssetQuery{}
	queryTopLevel.AssetId = testutil.TestTopLevelAssetId

	return &testScenario{
		name: "DescribeAssetHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeDescribeAsset,
				JSON:      testutil.SerializeStruct(t, query),
			},
			{
				RefID:     "B",
				QueryType: models.QueryTypeDescribeAsset,
				JSON:      testutil.SerializeStruct(t, queryTopLevel),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "describe-asset",
		handlerFn: func(s *testing.T, srvr *Server) backend.QueryDataHandlerFunc {
			return srvr.HandleDescribeAsset
		},
		validationFn: nil,
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

		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{},
			Queries:       scenario.queries,
		}

		srvr := &Server{
			datasource: mockedDatasource(scenario.mockSw),
		}

		qdr, err := scenario.handlerFn(t, srvr)(ctx, req)

		// this should always be nil, as the error is wrapped in the QueryDataResponse
		if err != nil {
			t.Fatal(err)
		}

		if scenario.validationFn != nil {
			scenario.validationFn(t, qdr)
		}

		// write out the golden for all data responses
		for i, dr := range qdr.Responses {
			fname := fmt.Sprintf("../testdata/%s-%s.golden.txt", scenario.goldenFileName, i)

			// temporary fix for golden files https://github.com/grafana/grafana-plugin-sdk-go/issues/213
			for _, fr := range dr.Frames {
				if fr.Meta != nil {
					fr.Meta.Custom = nil
				}
			}

			if err := experimental.CheckGoldenDataResponse(fname, &dr, true); err != nil {
				if !strings.Contains(err.Error(), "no such file or directory") {
					t.Fatal(err)
				}
			}
		}
	})
}

func TestHandlePropertyValueHistory(t *testing.T) {
	getPropertyValueHistoryHappyCase(t).run(t)
}

func TestHandleListAssetModels(t *testing.T) {
	listAssetModelsHappyCase(t).run(t)
}

func TestHandleListAssets(t *testing.T) {
	listAssetsHappyCase(t).run(t)
}

func TestHandleDescribeAsset(t *testing.T) {
	describeAssetHappyCase(t).run(t)
}
