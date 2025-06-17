// nolint
package resource

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	"github.com/patrickmn/go-cache"
)

type cachingResourceProvider struct {
	resources *SitewiseResources
	cache     *cache.Cache
}

func NewCachingResourceProvider(resources *SitewiseResources, c *cache.Cache) *cachingResourceProvider {
	return &cachingResourceProvider{
		resources: resources,
		cache:     c,
	}
}

func (cp *cachingResourceProvider) Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error) {
	val, ok := cp.cache.Get(assetId)
	if ok {
		a, ok := val.(iotsitewise.DescribeAssetOutput)
		if ok {
			return &a, nil
		}
	}

	a, err := cp.resources.Asset(ctx, assetId)
	if err != nil {
		return nil, err
	}
	cp.cache.Set(assetId, *a, -1)
	return a, nil
}

func (cp *cachingResourceProvider) Property(ctx context.Context, assetId string, propertyId string, propertyAlias string) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	key := assetId + "/" + propertyId
	if propertyAlias != "" {
		key = propertyAlias
	}
	val, ok := cp.cache.Get(key)
	if ok {
		a, ok := val.(iotsitewise.DescribeAssetPropertyOutput)
		if ok {
			return &a, nil
		}
	}

	a, err := cp.resources.Property(ctx, assetId, propertyId, propertyAlias)
	if err != nil {
		return nil, err
	}
	cp.cache.Set(key, *a, -1)
	return a, nil
}

func (cp *cachingResourceProvider) AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error) {
	val, ok := cp.cache.Get(modelId)
	if ok {
		a, ok := val.(iotsitewise.DescribeAssetModelOutput)
		if ok {
			return &a, nil
		}
	}

	a, err := cp.resources.AssetModel(ctx, modelId)
	if err != nil {
		return nil, err
	}
	cp.cache.Set(modelId, *a, -1)
	return a, nil
}
