package sitewise

import (
	"context"
	"fmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"

	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	gaws "github.com/grafana/iot-sitewise-datasource/pkg/common/aws"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type clientGetterFunc func(ctx backend.PluginContext, q models.BaseQuery) (client client.SitewiseClient, err error)
type invokerFunc func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error)

type Datasource struct {
	GetClient clientGetterFunc
}

func NewDatasource() *Datasource {
	return &Datasource{
		GetClient: func(ctx backend.PluginContext, q models.BaseQuery) (swclient client.SitewiseClient, err error) {
			swclient, err = client.GetClient(ctx, q.AwsRegion)
			return
		},
	}
}

func (ds *Datasource) invoke(ctx context.Context, req *backend.QueryDataRequest, baseQuery models.BaseQuery, invoker invokerFunc) (data.Frames, error) {
	sw, err := ds.GetClient(req.PluginContext, baseQuery)
	if err != nil {
		return nil, err
	}

	fr, err := invoker(ctx, sw)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, baseQuery, fr, sw)
}

func (ds *Datasource) HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) error {

	if settings, err := gaws.LoadSettings(*req.PluginContext.DataSourceInstanceSettings); err != nil {
		return errors.Wrap(err, "unable to load settings")
	} else {
		if sw, err := ds.GetClient(req.PluginContext, models.BaseQuery{AwsRegion: settings.DefaultRegion}); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to get client for region: %s", settings.DefaultRegion))
		} else {
			// todo: expand health check to test permission boundaries
			_, err = sw.ListAssetModelsWithContext(ctx, &iotsitewise.ListAssetModelsInput{MaxResults: aws.Int64(1)})
			return errors.Wrap(err, "unable to test ListAssetModels")
		}
	}
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return GetAssetPropertyValues(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleGetAssetPropertyAggregateQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return GetAssetPropertyAggregates(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleGetAssetPropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return GetAssetPropertyValue(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetModelsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return ListAssetModels(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return ListAssociatedAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return ListAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return DescribeAsset(ctx, sw, *query)
	})
}
