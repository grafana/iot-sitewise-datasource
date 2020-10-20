package resource

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/pkg/errors"
)

// CacheDuration is a constant that defines how long to keep cached elements before they are refreshed
const CacheDuration = time.Minute * 5

// CacheCleanupInterval is the interval at which the internal cache is cleaned / garbage collected
const CacheCleanupInterval = time.Minute * 10

// ErrNoValue is returned when a cached value is not available in the local cache
var ErrNoValue = errors.New("no cached value was found with that key")

type cachingProvider struct {
	resources *SitewiseResources
	cache     map[string]cachedResult
}

// cachedResult is a value and a timestamp that defines when the cached value is no longer usable
type cachedResult struct {
	Result    interface{}
	ExpiresAt time.Time
}

func newCachedResult(f interface{}) cachedResult {
	return cachedResult{
		ExpiresAt: time.Now().Add(CacheDuration),
		Result:    f,
	}
}

func (cp *cachingProvider) Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error) {

	return nil, nil
}

func (cp *cachingProvider) Property(ctx context.Context, assetId string, propertyId string) (*iotsitewise.DescribeAssetPropertyOutput, error) {

	return nil, nil
}

func (cp *cachingProvider) AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error) {
	// TODO: implement if needed
	return nil, nil
}
