package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	dserrors "github.com/grafana/iot-sitewise-datasource/pkg/errors"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
)

// DatasourceHandler is the plugin entrypoint and implements all of the necessary handler functions for dataqueries, healthchecks, and resources.
type DatasourceHandler struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirement
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

func GetDatasourceServeOpts() datasource.ServeOpts {
	handler := &DatasourceHandler{
		im: datasource.NewInstanceManager(server.NewServerInstance),
	}

	return datasource.ServeOpts{
		CheckHealthHandler: handler,
		QueryDataHandler:   handler,
	}
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (cr *DatasourceHandler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	h, err := cr.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	if srvr, ok := h.(*server.Server); ok {
		return server.HandleQueryData(ctx, srvr, req)
	}

	return nil, dserrors.ErrorBadDatasource
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (cr *DatasourceHandler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	h, err := cr.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	if srvr, ok := h.(*server.Server); ok {
		return srvr.HandleHealthCheck(ctx, req)
	}

	return nil, dserrors.ErrorBadDatasource
}
