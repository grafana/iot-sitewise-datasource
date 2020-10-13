package sitewise

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type Datasource interface {
	HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, query *models.AssetPropertyValueQuery, dataQuery backend.DataQuery)
}
