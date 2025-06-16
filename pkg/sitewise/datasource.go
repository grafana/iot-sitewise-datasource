package sitewise

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/grafana/grafana-aws-sdk/pkg/awsauth"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/proxy"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"

	"github.com/pkg/errors"
)

const EDGE_REGION string = "Edge"

type clientGetterFunc func(ctx context.Context, region string) (client.SitewiseAPIClient, error)
type invokerFunc func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error)

type Datasource struct {
	cfg               models.AWSSiteWiseDataSourceSetting
	edgeAuthenticator *EdgeAuthenticator
	proxyOptions      *proxy.Options
	GetClient         clientGetterFunc
}

type disableHostPrefixMiddleware struct{}

func (m *disableHostPrefixMiddleware) ID() string {
	return "DisableHostPrefixMiddleware"
}

func (m *disableHostPrefixMiddleware) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	ctx = smithyhttp.SetHostnameImmutable(ctx, true)
	return next.HandleInitialize(ctx, in)
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
	proxyOptions, err := settings.ProxyOptionsFromContext(ctx)
	if err != nil {
		backend.Logger.Error("failed to read proxy options", "error", err.Error())
		return nil, err
	}
	ds := &Datasource{
		cfg:          cfg,
		proxyOptions: proxyOptions,
	}

	if cfg.Region == models.EDGE_REGION && cfg.EdgeAuthMode != models.EDGE_AUTH_MODE_DEFAULT {
		ds.edgeAuthenticator = &EdgeAuthenticator{
			Settings: cfg,
		}

		err := ds.Authenticate()
		if err != nil {
			return nil, fmt.Errorf("error getting initial edge credentials (%s)", err.Error())
		}
	}

	return ds, nil
}

func (ds *Datasource) Authenticate() error {
	authInfo, err := ds.edgeAuthenticator.GetAuthInfo()
	if err != nil {
		return err
	}
	if authInfo == nil {
		return nil
	}
	ds.cfg.AuthType = awsds.AuthTypeKeys
	ds.cfg.AccessKey = authInfo.AccessKeyId
	ds.cfg.SecretKey = authInfo.SecretAccessKey
	ds.cfg.SessionToken = authInfo.SessionToken
	return nil
}

func (ds *Datasource) getClient(ctx context.Context, region string) (client.SitewiseAPIClient, error) {
	if ds.GetClient != nil {
		return ds.GetClient(ctx, region)
	}
	if err := ds.Authenticate(); err != nil {
		return nil, err
	}
	httpclient, err := client.GetHTTPClient(ds.cfg)
	if err != nil {
		return nil, err
	}

	awsCfg, err := awsauth.NewConfigProvider().GetConfig(ctx, awsauth.Settings{
		LegacyAuthType:     ds.cfg.AuthType,
		AccessKey:          ds.cfg.AccessKey,
		SecretKey:          ds.cfg.SecretKey,
		SessionToken:       ds.cfg.SessionToken,
		Region:             region,
		CredentialsProfile: ds.cfg.Profile,
		AssumeRoleARN:      ds.cfg.AssumeRoleARN,
		Endpoint:           ds.cfg.Endpoint,
		ExternalID:         ds.cfg.ExternalID,
		UserAgent:          awsds.GetUserAgentString("grafana-iot-sitewise-datasource"),
		HTTPClient:         httpclient,
		ProxyOptions:       ds.proxyOptions,
	})

	if err != nil {
		return nil, err
	}

	return &client.SitewiseClient{Client: iotsitewise.NewFromConfig(awsCfg, func(o *iotsitewise.Options) {
		if ds.cfg.Region == models.EDGE_REGION {
			o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
				return stack.Initialize.Add(&disableHostPrefixMiddleware{}, middleware.Before)
			})
		}
	})}, nil
}

func (ds *Datasource) invoke(ctx context.Context, _ *backend.QueryDataRequest, baseQuery *models.BaseQuery, invoker invokerFunc) (data.Frames, error) {
	sw, err := ds.getClient(ctx, baseQuery.AwsRegion)
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

	sw, err := ds.getClient(ctx, ds.cfg.Region)
	if err != nil {
		return errors.Wrap(err, "unable to load settings")
	}

	// todo: expand health check to test permission boundaries
	_, err = sw.ListAssetModels(ctx, &iotsitewise.ListAssetModelsInput{MaxResults: aws.Int32(1)})
	return errors.Wrap(err, "unable to test ListAssetModels")
}

func (ds *Datasource) HandleInterpolatedPropertyValueQuery(ctx context.Context, _ *backend.QueryDataRequest, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	sw, err := ds.getClient(ctx, query.AwsRegion)
	if err != nil {
		return nil, err
	}
	modifiedQuery, fr, err := api.GetInterpolatedAssetPropertyValues(ctx, sw, *query)
	if err != nil {
		return nil, err
	}
	return frameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
}

func (ds *Datasource) HandleGetAssetPropertyValueHistoryQuery(ctx context.Context, query *models.AssetPropertyValueQuery) (data.Frames, error) {
	sw, err := ds.getClient(ctx, query.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.AwsRegion == EDGE_REGION {
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
	sw, err := ds.getClient(ctx, query.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.AwsRegion == EDGE_REGION {
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
	sw, err := ds.getClient(ctx, query.AwsRegion)
	if err != nil {
		return nil, err
	}

	// Batch API is not available at the edge
	if query.AwsRegion == EDGE_REGION {
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
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ListAssetModels(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssociatedAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ListAssociatedAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetsQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ListAssets(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListTimeSeriesQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListTimeSeriesQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ListTimeSeries(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.DescribeAsset(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleDescribeAssetModelQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.DescribeAssetModelQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.DescribeAssetModel(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleListAssetPropertiesQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ListAssetPropertiesQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ListAssetProperties(ctx, sw, *query)
	})
}

func (ds *Datasource) HandleExecuteQuery(ctx context.Context, req *backend.QueryDataRequest, query *models.ExecuteQuery) (data.Frames, error) {
	return ds.invoke(ctx, req, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseAPIClient) (framer.Framer, error) {
		return api.ExecuteQuery(ctx, sw, *query)
	})
}
