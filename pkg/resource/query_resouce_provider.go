package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

type queryResourceProvider struct {
	resources *cachingProvider
	baseQuery models.BaseQuery
}

func NewQueryResourceProvider(client client.SitewiseClient, query models.BaseQuery) *queryResourceProvider {
	cachingProvider := NewCachingProvider(NewSitewiseResources(client))
	return &queryResourceProvider{
		resources: cachingProvider,
		baseQuery: query,
	}
}

func (rp *queryResourceProvider) Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error) {
	return rp.resources.Asset(ctx, rp.baseQuery.AssetId)
}

func (rp *queryResourceProvider) Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	return rp.resources.Property(ctx, rp.baseQuery.AssetId, rp.baseQuery.PropertyId)
}

func (rp *queryResourceProvider) AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error) {

	asset, err := rp.resources.Asset(ctx, rp.baseQuery.AssetId)

	if err != nil {
		return nil, err
	}

	return rp.resources.AssetModel(ctx, *asset.AssetModelId)
}
