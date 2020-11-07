package sitewise

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

func frameResponse(ctx context.Context, query models.BaseQuery, data framer.Framer, client client.SitewiseClient) (data.Frames, error) {
	rp := resource.NewQueryResourceProvider(client, query)
	return data.Frames(ctx, rp)
}
