package sitewise

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

// cacheDuration is a constant that defines how long to keep cached elements before they are refreshed
const cacheDuration = time.Minute * 5

// cacheCleanupInterval is the interval at which the internal cache is cleaned / garbage collected
const cacheCleanupInterval = time.Minute * 10

var GetCache = func() func() *cache.Cache {
	var gCache = cache.New(cacheDuration, cacheCleanupInterval) // max size not supported
	return func() *cache.Cache {
		return gCache
	}
}()

func frameResponse(ctx context.Context, query models.BaseQuery, data framer.Framer, client client.SitewiseClient) (data.Frames, error) {
	cp := resource.NewCachingResourceProvider(resource.NewSitewiseResources(client), GetCache())
	rp := resource.NewQueryResourceProvider(cp, query)
	return data.Frames(ctx, rp)
}

func FrameResponse(ctx context.Context, query models.BaseQuery, data framer.Framer, client client.SitewiseClient) (data.Frames, error) {
	return frameResponse(ctx, query, data, client)
}
