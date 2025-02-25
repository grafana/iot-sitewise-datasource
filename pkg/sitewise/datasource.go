package sitewise

import (
	"context"
	"fmt"
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

const EDGE_REGION string = "Edge"

type clientGetterFunc func(region string) (client client.SitewiseClient, err error)
type invokerFunc func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error)

type Datasource struct {
	GetClient clientGetterFunc
	// Refactor to work with region like clientGetterFunc
	RGTaggingClient     client.TaggingApiClient
	IncludedTagPatterns []map[string][]string
}

func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (*Datasource, error) {
	cfg := models.AWSSiteWiseDataSourceSetting{}

	err := cfg.Load(settings)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	sessions := awsds.NewSessionCache()
	authSettings := awsds.ReadAuthSettings(ctx)
	clientGetter := func(region string) (swclient client.SitewiseClient, err error) {
		swclient, err = client.GetClient(region, cfg, sessions.GetSessionWithAuthSettings, authSettings)
		return
	}

	if cfg.Region == models.EDGE_REGION && cfg.EdgeAuthMode != models.EDGE_AUTH_MODE_DEFAULT {
		edgeAuthenticator := EdgeAuthenticator{ //DummyAuthenticator{
			Settings: cfg,
		}

		var mu sync.Mutex
		authInfo, err := edgeAuthenticator.Authenticate()
		if err != nil {
			return nil, fmt.Errorf("error getting initial edge credentials (%s)", err.Error())
		}
		cfg.AuthType = awsds.AuthTypeKeys // Force key auth
		cfg.AccessKey = authInfo.AccessKeyId
		cfg.SecretKey = authInfo.SecretAccessKey
		cfg.SessionToken = authInfo.SessionToken

		clientGetter = func(region string) (swclient client.SitewiseClient, err error) {
			mu.Lock()
			if time.Now().After(authInfo.SessionExpiryTime) {
				log.DefaultLogger.Debug("edge credentials expired. updating credentials now.")
				authInfo, err = edgeAuthenticator.Authenticate()
				if err != nil {
					mu.Unlock()
					return nil, fmt.Errorf("error updating edge credentials (%s)", err.Error())
				}
				cfg.AccessKey = authInfo.AccessKeyId
				cfg.SecretKey = authInfo.SecretAccessKey
				cfg.SessionToken = authInfo.SessionToken
			}
			cfgCopy := cfg
			mu.Unlock()
			swclient, err = client.GetClient(region, cfgCopy, sessions.GetSessionWithAuthSettings, authSettings)
			return
		}
	}

	// Construct a ResourceGroupsTaggingAPI client
	rgTaggingClient, err := client.GetTaggingApiClient(cfg.Region, cfg, sessions.GetSessionWithAuthSettings, authSettings)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		GetClient:           clientGetter,
		RGTaggingClient:     rgTaggingClient,
		IncludedTagPatterns: cfg.IncludedTagPatterns,
	}, nil
}

func (ds *Datasource) invoke(ctx context.Context, _ *backend.QueryDataRequest, baseQuery *models.BaseQuery, invoker invokerFunc) (data.Frames, error) {
	sw, err := ds.GetClient(baseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	fr, err := invoker(ctx, sw)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, *baseQuery, fr, sw)
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

func (ds *Datasource) HandleInterpolatedPropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.GetInterpolatedAssetPropertyValues(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	sw, err := ds.GetClient(query.BaseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.BaseQuery.AwsRegion == EDGE_REGION {
		modifiedQuery, fr, err := api.GetAssetPropertyValues(ctx, sw, *query)
		if err != nil {
			return nil, err
		}

		return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
	}

	modifiedQuery, fr, err := api.BatchGetAssetPropertyValues(ctx, sw, *query)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
}

func (ds *Datasource) HandleGetAssetPropertyAggregateQuery(ctx context.Context, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	sw, err := ds.GetClient(query.BaseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.BaseQuery.AwsRegion == EDGE_REGION {
		modifiedQuery, fr, err := api.GetAssetPropertyValuesForTimeRange(ctx, sw, *query)
		if err != nil {
			return nil, err
		}

		return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
	}

	modifiedQuery, fr, err := api.BatchGetAssetPropertyValuesForTimeRange(ctx, sw, *query)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
}

func (ds *Datasource) HandleGetAssetPropertyValueQuery(ctx context.Context, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	sw, err := ds.GetClient(query.BaseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.BaseQuery.AwsRegion == EDGE_REGION {
		modifiedQuery, fr, err := api.GetAssetPropertyValue(ctx, sw, *query)
		if err != nil {
			return nil, err
		}

		return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
	}

	modifiedQuery, fr, err := api.BatchGetAssetPropertyValue(ctx, sw, *query)
	if err != nil {
		return nil, err
	}

	return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
}

func (ds *Datasource) HandleListAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetModelsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssetModels(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssociatedAssets(ctx, sw, *query, ds.RGTaggingClient, ds.IncludedTagPatterns)
		// return api.ListAssociatedAssets(ctx, sw, *query, ds.RGTaggingClient)
	})
}

func (ds *Datasource) HandleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssets(ctx, sw, *query, ds.RGTaggingClient, ds.IncludedTagPatterns)
		// return api.ListAssets(ctx, sw, *query, ds.RGTaggingClient)
	})
}

func (ds *Datasource) HandleListTimeSeriesQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListTimeSeriesQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListTimeSeries(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.DescribeAsset(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetModelQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetModelQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.DescribeAssetModel(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetPropertiesQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetPropertiesQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssetProperties(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleExecuteQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ExecuteQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ExecuteQuery(ctx, sw, *query)
	})
}
