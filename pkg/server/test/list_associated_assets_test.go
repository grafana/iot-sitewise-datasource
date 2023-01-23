package test

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/mock"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

func TestHandleListAssociatedAssets(t *testing.T) {
	listAssociatedAssetsChildrenHappyCase(t).run(t)
	listAssociatedAssetsParentHappyCase(t).run(t)
}

var listAssociatedAssetsChildrenHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		hierarchyId := req.HierarchyId
		traversal := req.TraversalDirection
		if assetId == nil || hierarchyId == nil || traversal == nil {
			return false
		}
		return testdata.DemoWindFarmAssetId == *assetId &&
			testdata.TurbineAssetModelHierarchyId == *hierarchyId
	})

	mockSw.On("ListAssociatedAssetsWithContext", mock.Anything, argMatcher).Return(&assets, nil)

	query := models.ListAssociatedAssetsQuery{}
	query.AssetIds = []string{testdata.DemoWindFarmAssetId}
	query.HierarchyId = testdata.TurbineAssetModelHierarchyId

	return &testScenario{
		name: "ListAssociatedAssetsChildHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssociatedAssets,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "list-associated-assets",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssociatedAssets
		},
		validationFn: nil,
	}
}

var listAssociatedAssetsParentHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseClient{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets-parent.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		traversal := req.TraversalDirection
		if assetId == nil || traversal == nil {
			return false
		}
		return testdata.DemoTurbineAsset1 == *assetId &&
			"PARENT" == *traversal
	})

	mockSw.On("ListAssociatedAssetsWithContext", mock.Anything, argMatcher).Return(&assets, nil)

	query := models.ListAssociatedAssetsQuery{}
	query.AssetIds = []string{testdata.DemoTurbineAsset1}

	return &testScenario{
		name: "ListAssociatedAssetsParentHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssociatedAssets,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "list-associated-assets-parent",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssociatedAssets
		},
		validationFn: nil,
	}
}
