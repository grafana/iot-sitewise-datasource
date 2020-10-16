package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

type Datasource interface {
	HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) error
	HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (framer.Framer, error)
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, srvr *Server, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(srvr)
	return m.QueryData(ctx, req)
}
