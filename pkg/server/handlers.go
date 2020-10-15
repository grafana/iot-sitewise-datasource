package server

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

func (s *Server) handlePropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query := &models.AssetPropertyValueQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	fr, err := s.datasource.HandleGetAssetPropertyValueHistoryQuery(ctx, req, query)

	return framer.FrameResponseWithError(fr, ctx, err)
}

func (s *Server) HandlePropertyValueHistory(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePropertyValueQuery),
	}, nil
}
