package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

type Datasource interface {
	HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (framer.Framer, error)
}
