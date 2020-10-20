package sitewise

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	gaws "github.com/grafana/iot-sitewise-datasource/pkg/common/aws"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type clientGetterFunc func(ctx backend.PluginContext, q models.BaseQuery) (client client.Client, err error)

type Datasource struct {
	GetClient clientGetterFunc
}

func NewDatasource() *Datasource {
	return &Datasource{
		GetClient: func(ctx backend.PluginContext, q models.BaseQuery) (swclient client.Client, err error) {
			swclient, err = client.GetClient(ctx, q.AwsRegion)
			return
		},
	}
}

func (ds *Datasource) HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) error {

	if settings, err := gaws.LoadSettings(*req.PluginContext.DataSourceInstanceSettings); err != nil {
		return errors.Wrap(err, "unable to load settings")
	} else {
		if sw, err := ds.GetClient(req.PluginContext, models.BaseQuery{AwsRegion: settings.DefaultRegion}); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to get client for region: %s", settings.DefaultRegion))
		} else {
			// todo: expand health check to test permission boundaries
			_, err = sw.ListAssetModels(&iotsitewise.ListAssetModelsInput{MaxResults: aws.Int64(1)})
			return errors.Wrap(err, "unable to test ListAssetModels")
		}
	}
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {

	sw, err := ds.GetClient(req.PluginContext, query.BaseQuery)
	if err != nil {
		return nil, err
	}

	fdata, err := GetAssetPropertyValues(ctx, sw, *query)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, query.BaseQuery, fdata, sw)
}
