package server

import (
	"context"
	"fmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"
)

type Server struct {
	Datasource    *sitewise.Datasource
	channelPrefix string
	closeCh       chan struct{}
	queryMux      *datasource.QueryTypeMux
}

// Make sure SampleDatasource implements required interfaces.
// This is important to do since otherwise we will only get a
// not implemented error response from plugin in runtime.
var (
	_ backend.QueryDataHandler      = (*Server)(nil)
	_ backend.CheckHealthHandler    = (*Server)(nil)
	_ instancemgmt.InstanceDisposer = (*Server)(nil)
)

// QueryHandlerFunc is the function signature used for mux.HandleFunc
// Looks like mux.HandleFunc uses backend.QueryHandlerFunc
// type QueryDataHandlerFunc func(ctx context.Context, req *QueryDataRequest) (*QueryDataResponse, error)
type QueryHandlerFunc func(context.Context, *backend.QueryDataRequest, backend.DataQuery) backend.DataResponse

func DataResponseErrorUnmarshal(err error) backend.DataResponse {
	return backend.DataResponse{
		Error: errors.Wrap(err, "failed to unmarshal JSON request into query"),
	}
}

func DataResponseErrorRequestFailed(err error) backend.DataResponse {
	return backend.DataResponse{
		Error: errors.Wrap(err, "failed to fetch query data"),
		ErrorSource: backend.ErrorSourceDownstream,
	}
}

type handler func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func getQueryHandlers(s *Server) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(models.QueryTypePropertyValueHistory, s.lastObservation(s.HandlePropertyValueHistory))
	mux.HandleFunc(models.QueryTypePropertyAggregate, s.lastObservation(s.HandlePropertyAggregate))
	mux.HandleFunc(models.QueryTypePropertyInterpolated, s.lastObservation(s.HandleInterpolatedPropertyValue))
	mux.HandleFunc(models.QueryTypePropertyValue, s.HandlePropertyValue)
	mux.HandleFunc(models.QueryTypeListAssetModels, s.HandleListAssetModels)
	mux.HandleFunc(models.QueryTypeListAssociatedAssets, s.HandleListAssociatedAssets)
	mux.HandleFunc(models.QueryTypeListAssets, s.HandleListAssets)
	mux.HandleFunc(models.QueryTypeDescribeAsset, s.HandleDescribeAsset)
	mux.HandleFunc(models.QueryTypeListAssetProperties, s.HandleListAssetProperties)
	mux.HandleFunc(models.QueryTypeListTimeSeries, s.HandleListTimeSeries)
	mux.HandleFunc(models.QueryTypeExecuteQuery, s.HandleExecuteQuery)

	return mux
}

func NewServerInstance(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds, err := sitewise.NewDatasource(ctx, settings)
	if err != nil {
		return nil, err
	}
	srvr := &Server{
		Datasource:    ds,
		channelPrefix: fmt.Sprintf("ds/%d/", settings.ID),
		closeCh:       make(chan struct{}),
	}
	srvr.queryMux = getQueryHandlers(srvr) // init once
	return srvr, nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (s *Server) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return s.queryMux.QueryData(ctx, req)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (s *Server) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := s.Datasource.HealthCheck(ctx, req); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: backend.HealthStatusOk.String(),
	}, nil
}

func (s *Server) Dispose() {
	close(s.closeCh)
}
