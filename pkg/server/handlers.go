package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

func (s *Server) handlePropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseError(err, "failed to unmarshal JSON request into query")
	}

	fr, err := s.datasource.HandleGetAssetPropertyValueHistoryQuery(ctx, req, query)

	if err != nil {
		return DataResponseError(err, "failed to fetch query data")
	}

	return framer.FrameResponse(ctx, fr)
}

func (s *Server) HandlePropertyValueHistory(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePropertyValueQuery),
	}, nil
}
