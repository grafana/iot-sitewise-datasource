package sitewise

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

type clientGetterFunc func(ctx backend.PluginContext, q models.BaseQuery) (client client.Client, err error)

type Datasource struct {
	getClient clientGetterFunc
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (framer.Framer, error) {

	sw, err := ds.getClient(req.PluginContext, query.BaseQuery)
	if err != nil {
		return nil, err
	}

	fdata, err := GetAssetPropertyValues(ctx, sw, *query)
	if err != nil {
		return nil, err
	}

	return framePropertyValueResponse(query, fdata, sw), nil
}
