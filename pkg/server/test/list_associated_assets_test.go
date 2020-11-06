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
	listAssociatedAssetsHappyCase(t).run(t)
}

var listAssociatedAssetsHappyCase testServerScenarioFn = func(t *testing.T) *testScenario {

	mockSw := &mocks.Client{}

	assets := testdata.GetIoTSitewiseAssociatedAssets(t, testDataRelativePath("list-associated-assets.json"))

	argMatcher := mock.MatchedBy(func(req *iotsitewise.ListAssociatedAssetsInput) bool {
		assetId := req.AssetId
		hierarchyId := req.HierarchyId
		if assetId == nil || hierarchyId == nil {
			return false
		}
		return testdata.TestTopLevelAssetId == *assetId && testdata.TestTopLevelAssetHierarchyId == *hierarchyId
	})

	mockSw.On("ListAssociatedAssetsWithContext", mock.Anything, argMatcher).Return(&assets, nil)

	query := models.ListAssociatedAssetsQuery{}
	query.AssetId = testdata.TestTopLevelAssetId
	query.HierarchyId = testdata.TestTopLevelAssetHierarchyId

	return &testScenario{
		name: "ListAssociatedAssetsHappyCase",
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
