package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func (s *Server) HandleHealthCheck(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	if err := s.datasource.HealthCheck(ctx, req); err != nil {
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

func (s *Server) handlePropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseError(err, "failed to unmarshal JSON request into query")
	}

	frames, err := s.datasource.HandleGetAssetPropertyValueHistoryQuery(ctx, req, query)

	if err != nil {
		return DataResponseError(err, "failed to fetch query data")
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) HandlePropertyValueHistory(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePropertyValueQuery),
	}, nil
}
