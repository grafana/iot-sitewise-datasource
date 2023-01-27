package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type queryResourceProvider struct {
	resources *cachingResourceProvider
	baseQuery models.BaseQuery
}

func NewQueryResourceProvider(cachingProvider *cachingResourceProvider, query models.BaseQuery) *queryResourceProvider {
	//cachingResourceProvider := NewCachingResourceProvider(NewSitewiseResources(client))
	return &queryResourceProvider{
		resources: cachingProvider,
		baseQuery: query,
	}
}

func (rp *queryResourceProvider) Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error) {
	assetId := ""

	// use the first asset id if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	return rp.resources.Asset(ctx, assetId)
}

func (rp *queryResourceProvider) Assets(ctx context.Context) (map[string]*iotsitewise.DescribeAssetOutput, error) {
	assets := map[string]*iotsitewise.DescribeAssetOutput{}
	for _, id := range rp.baseQuery.AssetIds {
		asset, err := rp.resources.Asset(ctx, id)
		if err != nil {
			return nil, err
		}
		assets[id] = asset
	}
	return assets, nil
}

func (rp *queryResourceProvider) Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	assetId := ""

	// use the first asset id if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	return rp.resources.Property(ctx, assetId, rp.baseQuery.PropertyId, rp.baseQuery.PropertyAlias)
}

func (rp *queryResourceProvider) Properties(ctx context.Context) (map[string]*iotsitewise.DescribeAssetPropertyOutput, error) {
	properties := map[string]*iotsitewise.DescribeAssetPropertyOutput{}
	for _, id := range rp.baseQuery.AssetIds {
		prop, err := rp.resources.Property(ctx, id, rp.baseQuery.PropertyId, rp.baseQuery.PropertyAlias)
		if err != nil {
			return nil, err
		}
		properties[id] = prop
	}
	return properties, nil
}

func (rp *queryResourceProvider) AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error) {
	assetId := ""

	// use the first asset id if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	asset, err := rp.resources.Asset(ctx, assetId)

	if err != nil {
		return nil, err
	}

	return rp.resources.AssetModel(ctx, *asset.AssetModelId)
}
