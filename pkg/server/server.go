package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	}
}

type handler func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)

func (s *Server) fallbackToLastObservation(h handler) handler {
	return func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
		resp := &backend.QueryDataResponse{
			Responses: make(map[string]backend.DataResponse),
		}

		queries := make(map[string]backend.DataQuery, 0)
		for _, q := range req.Queries {
			queries[q.RefID] = q
		}

		origResp, err := h(ctx, req)
		if err != nil {
			return nil, err
		}

		for refID, res := range origResp.Responses {
			if res.Error != nil {
				resp.Responses[refID] = res
				continue
			}

			query := queries[refID]
			assetQuery, err := models.GetAssetPropertyValueQuery(&query)
			if err != nil || !assetQuery.LastObservation || assetQuery.NextToken != "" {
				resp.Responses[refID] = res
				continue
			}
			switch query.QueryType {
			case models.QueryTypePropertyValueHistory:
				fallthrough
			case models.QueryTypePropertyAggregate:
				fallthrough
			case models.QueryTypePropertyInterpolated:
			default:
				continue
			}

			if len(res.Frames) == 0 || res.Frames[0].Rows() == 0 {
				// No data, fallback to last observation
				lastValueRes, err := s.lastValueQuery(ctx, query)
				if err != nil {
					return nil, err
				}
				resp.Responses[refID] = lastValueRes
				continue
			}

			resp.Responses[refID] = res
		}

		return resp, nil
	}
}

func (s *Server) lastValueQuery(ctx context.Context, query backend.DataQuery) (backend.DataResponse, error) {
	query.QueryType = models.QueryTypePropertyValueHistory
	query.TimeRange.From = time.Unix(0, 0)
	query.MaxDataPoints = 1

	assetQuery, err := models.GetAssetPropertyValueQuery(&query)
	if err != nil {
		return backend.DataResponse{}, err
	}

	assetQuery.TimeOrdering = "DESCENDING"
	query.JSON, err = json.Marshal(&assetQuery)
	if err != nil {
		return backend.DataResponse{}, err
	}

	res, err := s.QueryData(ctx, &backend.QueryDataRequest{Queries: []backend.DataQuery{query}})
	if err != nil {
		return backend.DataResponse{}, err
	}

	for refID, r := range res.Responses {
		if refID == query.RefID && len(r.Frames) > 0 && r.Frames[0].Rows() > 0 {
			firstRow := r.Frames[0].RowCopy(0)
			firstRow[0] = query.TimeRange.To
			r.Frames[0].AppendRow(firstRow...)
			if meta, ok := r.Frames[0].Meta.Custom.(models.SitewiseCustomMeta); ok {
				meta.NextToken = ""
				r.Frames[0].Meta.Custom = meta
			}
			return r, nil
		}
	}

	return backend.DataResponse{}, fmt.Errorf("no response for refID %s", query.RefID)
}

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func getQueryHandlers(s *Server) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(models.QueryTypePropertyValueHistory, s.fallbackToLastObservation(s.HandlePropertyValueHistory))
	mux.HandleFunc(models.QueryTypePropertyAggregate, s.fallbackToLastObservation(s.HandlePropertyAggregate))
	mux.HandleFunc(models.QueryTypePropertyInterpolated, s.fallbackToLastObservation(s.HandleInterpolatedPropertyValue))
	mux.HandleFunc(models.QueryTypePropertyValue, s.HandlePropertyValue)
	mux.HandleFunc(models.QueryTypeListAssetModels, s.HandleListAssetModels)
	mux.HandleFunc(models.QueryTypeListAssociatedAssets, s.HandleListAssociatedAssets)
	mux.HandleFunc(models.QueryTypeListAssets, s.HandleListAssets)
	mux.HandleFunc(models.QueryTypeDescribeAsset, s.HandleDescribeAsset)

	return mux
}

func NewServerInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds, err := sitewise.NewDatasource(settings)
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
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
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
