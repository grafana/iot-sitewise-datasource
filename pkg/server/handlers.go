package server

import (
	"context"
	"math"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func processQueries(ctx context.Context, req *backend.QueryDataRequest, handler QueryHandlerFunc) *backend.QueryDataResponse {
	res := backend.Responses{}
	for _, v := range req.Queries {
		res[v.RefID] = handler(ctx, req, v)
	}

	return &backend.QueryDataResponse{
		Responses: res,
	}
}

func (s *Server) HandleInterpolatedPropertyValue(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleInterpolatedPropertyValueQuery), nil
}

func (s *Server) HandlePropertyValueHistory(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handlePropertyValueHistoryQuery), nil
}

func (s *Server) HandlePropertyAggregate(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handlePropertyAggregateQuery), nil
}

func (s *Server) HandlePropertyValue(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handlePropertyValueQuery), nil
}

func (s *Server) HandleListAssetModels(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleListAssetModelsQuery), nil
}

func (s *Server) HandleListAssets(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleListAssetsQuery), nil
}

func (s *Server) HandleDescribeAsset(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleDescribeAssetQuery), nil
}

func (s *Server) HandleListTimeSeries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleListTimeSeriesQuery), nil
}

func (s *Server) HandleListAssetProperties(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleListAssetPropertiesQuery), nil
}

func (s *Server) HandleListAssociatedAssets(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleListAssociatedAssetsQuery), nil
}

func (s *Server) HandleDescribeAssetModel(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleDescribeAssetModelQuery), nil
}

func (s *Server) HandleExecuteQuery(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return processQueries(ctx, req, s.handleExecuteQuery), nil
}

func (s *Server) handleInterpolatedPropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleInterpolatedPropertyValueQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handlePropertyValueHistoryQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	// Expressions need to run synchronously so we set MaxPageAggregations
	// and MaxDataPoints to infinity to ensure that the query is not paginated.
	_, isFromExpression := req.Headers["http_X-Grafana-From-Expr"]
	_, isFromAlert := req.Headers["FromAlert"]
	if isFromAlert || isFromExpression {
		query.MaxPageAggregations = math.MaxInt32
		query.MaxDataPoints = math.MaxInt32
	}

	frames, err := s.Datasource.HandleGetAssetPropertyValueHistoryQuery(ctx, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	if len(frames) > 0 && query.ResponseFormat == "timeseries" {
		for i, frame := range frames {
			wide, err := data.LongToWide(frame, &data.FillMissing{Mode: data.FillModeNull, Value: math.NaN()})
			if err == nil {
				frames[i] = wide
			}
		}
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handlePropertyAggregateQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	// Expressions need to run synchronously so we set MaxPageAggregations
	// and MaxDataPoints to infinity to ensure that the query is not paginated.
	_, isFromExpression := req.Headers["http_X-Grafana-From-Expr"]
	_, isFromAlert := req.Headers["FromAlert"]
	if isFromAlert || isFromExpression {
		query.MaxPageAggregations = math.MaxInt32
		query.MaxDataPoints = math.MaxInt32
	}

	frames, err := s.Datasource.HandleGetAssetPropertyAggregateQuery(ctx, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	if len(frames) > 0 && query.ResponseFormat == "timeseries" {
		for _, frame := range frames {
			wide, err := data.LongToWide(frame, &data.FillMissing{Mode: data.FillModeNull, Value: math.NaN()})
			if err == nil {
				frames = []*data.Frame{wide}
			}
		}
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handlePropertyValueQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleGetAssetPropertyValueQuery(ctx, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	if len(frames) > 0 && query.ResponseFormat == "timeseries" {
		for _, frame := range frames {
			wide, err := data.LongToWide(frame, &data.FillMissing{Mode: data.FillModeNull, Value: math.NaN()})
			if err == nil {
				frames = []*data.Frame{wide}
			}
		}
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleListAssetModelsQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListAssetModelsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleListAssetModelsQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleListAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetListAssetsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleListAssetsQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleListAssociatedAssetsQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetListAssociatedAssetsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleListAssociatedAssetsQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleListTimeSeriesQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListTimeSeriesQuery(&q)

	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleListTimeSeriesQuery(ctx, req, query)

	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleDescribeAssetQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetDescribeAssetQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleDescribeAssetQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleListAssetPropertiesQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListAssetPropertiesQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleListAssetPropertiesQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleDescribeAssetModelQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {

	query, err := models.GetDescribeAssetModelQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := s.Datasource.HandleDescribeAssetModelQuery(ctx, req, query)
	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (s *Server) handleExecuteQuery(ctx context.Context, req *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetExecuteQuery(&q)
	if err != nil {
		log.DefaultLogger.FromContext(ctx).Warn("Error unmarshalling query", "error", err)
		return DataResponseErrorUnmarshal(err)
	}

	query.RawSQL, err = sqlutil.Interpolate(&query.Query, s.Datasource.Macros())
	if err != nil {
		log.DefaultLogger.Warn("Error interpolating query", "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, "macro interpolate: "+err.Error())
	}

	frames, err := s.Datasource.HandleExecuteQuery(ctx, req, query)
	if err != nil {
		log.DefaultLogger.FromContext(ctx).Warn("Error executing query", "error", err)
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}
