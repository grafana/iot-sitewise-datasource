package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type Datasource interface {
	HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) error
	HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error)
	HandleGetAssetPropertyAggregateQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error)
	HandleGetAssetPropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error)
	HandleListAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetModelsQuery) (data.Frames, error)
	HandleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetsQuery) (data.Frames, error)
	HandleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetQuery) (data.Frames, error)
	HandleListAssociatedAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error)
}

// HandleQueryData handles the `QueryData` request for the Sitewise Datasource
func HandleQueryData(ctx context.Context, srvr *Server, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(srvr)
	return m.QueryData(ctx, req)
}
