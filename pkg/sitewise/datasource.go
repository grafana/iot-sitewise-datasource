package sitewise

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
	"github.com/pkg/errors"
)

type clientGetterFunc func(region string) (client client.SitewiseClient, err error)
type invokerFunc func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error)

type Datasource struct {
	GetClient clientGetterFunc
	closeCh   chan struct{}
	closeOnce sync.Once
}

func NewDatasource(settings backend.DataSourceInstanceSettings) (*Datasource, error) {
	cfg := models.AWSSiteWiseDataSourceSetting{}

	err := cfg.Load(settings)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})
	sessions := awsds.NewSessionCache()
	clientGetter := func(region string) (swclient client.SitewiseClient, err error) {
		swclient, err = client.GetClient(region, cfg, sessions.GetSession)
		return
	}

	if cfg.Region == models.EDGE_REGION && cfg.EdgeAuthMode != models.EDGE_AUTH_MODE_DEFAULT {
		edgeAuthenticator := EdgeAuthenticator{
			Settings: cfg,
		}

		var mu sync.Mutex
		var waitTime time.Duration

		updateAuth := func() error {
			authInfo, err := edgeAuthenticator.Authenticate()
			if err == nil {
				mu.Lock()
				cfg.AccessKey = authInfo.AccessKeyId
				cfg.SecretKey = authInfo.SecretAccessKey
				cfg.SessionToken = authInfo.SessionToken
				cfg.AuthType = awsds.AuthTypeKeys
				mu.Unlock()
				waitTime = time.Until(authInfo.SessionExpiryTime)
				log.DefaultLogger.Debug("should wait for: ", "time:", waitTime)
			} else {
				waitTime = waitTime + time.Second*10
			}
			return err
		}

		err = updateAuth()
		if err != nil {
			return &Datasource{}, err
		}

		go func() {
			for {
				log.DefaultLogger.Debug("wait time until next credential fetch: ", "time:", waitTime)
				select {
				case <-time.After(waitTime):
					log.DefaultLogger.Debug("updating edge auth credentials now")
					_ = updateAuth()
				case <-done:
					return
				}
			}
		}()

		clientGetter = func(region string) (swclient client.SitewiseClient, err error) {
			// TODO (new PR?) check on demand
			mu.Lock()
			cfgCopy := cfg
			mu.Unlock()
			swclient, err = client.GetClient(region, cfgCopy, sessions.GetSession)
			return
		}
	}

	return &Datasource{
		GetClient: clientGetter,
		closeCh:   done,
	}, nil
}

func (ds *Datasource) invoke(ctx context.Context, req *backend.QueryDataRequest, baseQuery models.BaseQuery, invoker invokerFunc) (data.Frames, error) {
	sw, err := ds.GetClient(baseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	fr, err := invoker(ctx, sw)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, baseQuery, fr, sw)
}

func (ds *Datasource) Dispose() {
	ds.closeOnce.Do(func() {
		close(ds.closeCh)
	})
}

func (ds *Datasource) HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) error {

	sw, err := ds.GetClient("") // Default region
	if err != nil {
		return errors.Wrap(err, "unable to load settings")
	}

	// todo: expand health check to test permission boundaries
	_, err = sw.ListAssetModelsWithContext(ctx, &iotsitewise.ListAssetModelsInput{MaxResults: aws.Int64(1)})
	return errors.Wrap(err, "unable to test ListAssetModels")
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.GetAssetPropertyValues(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleGetAssetPropertyAggregateQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.GetAssetPropertyValuesForTimeRange(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleGetAssetPropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.GetAssetPropertyValue(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetModelsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssetModels(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssociatedAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.DescribeAsset(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetModelQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetModelQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.DescribeAssetModel(ctx, sw, *query)
	})
}
