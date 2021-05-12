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
	HandleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error)
	HandleDescribeAssetModelQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetModelQuery) (data.Frames, error)
}
