package test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/stretchr/testify/mock"
)

func TestHandleListAssociatedAssets(t *testing.T) {
	listAssociatedAssetsChildrenHappyCase(t).run(t)
	listAssociatedAssetsParentHappyCase(t).run(t)
	listAssociatedAssetsMultipleAssetIdsCase(t).run(t)
}

var listAssociatedAssetsChildrenHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.SitewiseAPIClient{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		hierarchyId := req.HierarchyId
		if assetId == nil || hierarchyId == nil {
			return false
		}
		return testdata.DemoWindFarmAssetId == *assetId &&
			testdata.TurbineAssetModelHierarchyId == *hierarchyId
	})

	mockSw.On("ListAssociatedAssets", mock.Anything, argMatcher).Return(&assets, nil)

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

	mockSw := &mocks.SitewiseAPIClient{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets-parent.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		traversal := req.TraversalDirection
		if assetId == nil {
			return false
		}
		return testdata.DemoTurbineAsset1 == *assetId &&
			iotsitewisetypes.TraversalDirectionParent == traversal
	})

	mockSw.On("ListAssociatedAssets", mock.Anything, argMatcher).Return(&assets, nil)

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

var listAssociatedAssetsMultipleAssetIdsCase testServerScenarioFn = func(t *testing.T) *testScenario {
	mockSw := &mocks.SitewiseAPIClient{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets.json"))
	assets2 := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets-multiple.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		hierarchyId := req.HierarchyId
		if assetId == nil || hierarchyId == nil {
			return false
		}
		return testdata.DemoWindFarmAssetId == *assetId &&
			testdata.TurbineAssetModelHierarchyId == *hierarchyId
	})

	argMatcher2 := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		hierarchyId := req.HierarchyId
		if assetId == nil || hierarchyId == nil {
			return false
		}
		return testdata.DemoWindFarmAssetId2 == *assetId &&
			testdata.TurbineAssetModelHierarchyId == *hierarchyId
	})

	mockSw.On("ListAssociatedAssets", mock.Anything, argMatcher).Return(&assets, nil)
	mockSw.On("ListAssociatedAssets", mock.Anything, argMatcher2).Return(&assets2, nil)

	query := models.ListAssociatedAssetsQuery{}
	query.AssetIds = []string{testdata.DemoWindFarmAssetId, testdata.DemoWindFarmAssetId2}
	query.HierarchyId = testdata.TurbineAssetModelHierarchyId

	return &testScenario{
		name: "ListAssociatedAssetsMultipleAssetIdsHappyCase",
		queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: models.QueryTypeListAssociatedAssets,
				JSON:      testdata.SerializeStruct(t, query),
			},
		},
		mockSw:         mockSw,
		goldenFileName: "list-associated-assets-multiple",
		handlerFn: func(srvr *server.Server) backend.QueryDataHandlerFunc {
			return srvr.HandleListAssociatedAssets
		},
		validationFn: nil,
	}
}
