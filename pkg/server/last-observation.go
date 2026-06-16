package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
)

func (s *Server) lastObservation(h handler) handler {
	return func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
		resp := &backend.QueryDataResponse{
			Responses: make(map[string]backend.DataResponse),
		}

		queries := make(map[string]backend.DataQuery, len(req.Queries))
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
			resp.Responses[refID] = res

			// ensure this is the last page of data
			if len(res.Frames) > 0 {
				if meta, ok := res.Frames[0].Meta.Custom.(models.SitewiseCustomMeta); ok && meta.NextToken != "" {
					continue
				}
			}

			// ensure that this is a supported query type, and that the user requested last observation
			assetQuery, err := models.GetAssetPropertyValueQuery(&query)
			if err != nil || !assetQuery.LastObservation {
				continue
			}

			lastValueRes, err := s.lastValueQuery(ctx, query, iotsitewisetypes.TimeOrderingDescending)
			if err != nil {
				log.DefaultLogger.Debug("failed to fetch last observation", "error", err)
			}
			if r, ok := resp.Responses[refID]; err == nil && ok {
				resp.Responses[refID] = mergeLastValueResponse(r, lastValueRes)
			}

			nextValueRes, err := s.lastValueQuery(ctx, query, iotsitewisetypes.TimeOrderingAscending)
			if err != nil {
				log.DefaultLogger.Debug("failed to fetch next observation", "error", err)
			}
			if r, ok := resp.Responses[refID]; err == nil && ok {
				resp.Responses[refID] = mergeLastValueResponse(r, nextValueRes)
			}
		}

		return resp, nil
	}
}

func (s *Server) lastValueQuery(ctx context.Context, query backend.DataQuery, timeOrdering iotsitewisetypes.TimeOrdering) (backend.DataResponse, error) {
	query.MaxDataPoints = 1
	switch timeOrdering {
	case iotsitewisetypes.TimeOrderingDescending:
		query.TimeRange.To = query.TimeRange.From.Add(-1 * time.Second)
		query.TimeRange.From = query.TimeRange.From.Add(-8760 * time.Hour) // 1 year ago

	case iotsitewisetypes.TimeOrderingAscending:
		query.TimeRange.From = query.TimeRange.To.Add(time.Second)
		query.TimeRange.To = time.Now()

	}
	assetQuery, err := models.GetAssetPropertyValueQuery(&query)
	if err != nil {
		return backend.DataResponse{}, err
	}
	assetQuery.NextToken = ""
	assetQuery.TimeOrdering = timeOrdering
	assetQuery.LastObservation = false
	assetQuery.MaxDataPoints = 1
	assetQuery.MaxPageAggregations = 1
	assetQuery.TimeRange = query.TimeRange
	assetQuery.Resolution = propvals.ResolutionMinute

	log.DefaultLogger.Debug("last observation query", "timeOrdering", timeOrdering, "timeRange", assetQuery.TimeRange)
	query.JSON, err = json.Marshal(&assetQuery)
	if err != nil {
		return backend.DataResponse{}, err
	}

	res, err := s.QueryData(ctx, &backend.QueryDataRequest{Queries: []backend.DataQuery{query}})
	if err != nil {
		return backend.DataResponse{}, err
	}

	dataRes, ok := res.Responses[query.RefID]
	if !ok || dataRes.Error != nil || len(dataRes.Frames) == 0 || dataRes.Frames[0].Rows() == 0 {
		return backend.DataResponse{}, fmt.Errorf("no response for query %s", query.RefID)
	}

	log.DefaultLogger.Debug("last observation response", "timeOrdering", timeOrdering, "timeRange", assetQuery.TimeRange, "time", getFirstTime(dataRes))

	frame := dataRes.Frames[0].EmptyCopy()
	frame.AppendRow(dataRes.Frames[0].RowCopy(0)...)

	return backend.DataResponse{Frames: []*data.Frame{frame}}, nil
}

func mergeLastValueResponse(originalRes, lastValueRes backend.DataResponse) backend.DataResponse {
	// always return the original response if either response has an error
	if hasError(originalRes) || hasError(lastValueRes) {
		log.DefaultLogger.Debug("has error", "original", hasError(originalRes), "lastValue", hasError(lastValueRes))
		return originalRes
	}

	if isEmpty(lastValueRes) {
		log.DefaultLogger.Debug("last value response is empty")
		return originalRes
	}

	if isEmpty(originalRes) {
		log.DefaultLogger.Debug("original response is empty")
		return lastValueRes
	}

	if !fieldsMatch(originalRes, lastValueRes) {
		log.DefaultLogger.Debug("fields do not match")
		return originalRes
	}

	originalTime := getFirstTime(originalRes)
	lastValueTime := getFirstTime(lastValueRes)
	if originalTime.After(lastValueTime) {
		log.DefaultLogger.Debug("original time is after last value time")
		originalRes.Frames[0].InsertRow(0, lastValueRes.Frames[0].RowCopy(0)...)
		return originalRes
	}

	log.DefaultLogger.Debug("last value time is after original time")

	originalRes.Frames[0].AppendRow(lastValueRes.Frames[0].RowCopy(0)...)
	return originalRes
}

func hasError(r backend.DataResponse) bool {
	return r.Error != nil
}

func isEmpty(r backend.DataResponse) bool {
	return len(r.Frames) == 0 || r.Frames[0].Rows() == 0
}

func fieldsMatch(originalRes, lastValueRes backend.DataResponse) bool {
	if len(originalRes.Frames[0].Fields) != len(lastValueRes.Frames[0].Fields) {
		return false
	}

	for i := 0; i < len(originalRes.Frames[0].Fields); i++ {
		if originalRes.Frames[0].Fields[i].Name != lastValueRes.Frames[0].Fields[i].Name {
			return false
		}
	}

	return true
}

func getFirstTime(r backend.DataResponse) time.Time {
	if len(r.Frames) > 0 && r.Frames[0].Rows() > 0 {
		for _, f := range r.Frames[0].Fields {
			if f.Type() == data.FieldTypeTime {
				if t, ok := f.At(0).(time.Time); ok {
					return t
				}
			}
			if f.Type() == data.FieldTypeNullableTime {
				if t, ok := f.At(0).(*time.Time); ok {
					return *t
				}
			}
		}
	}
	return time.Time{}
}
