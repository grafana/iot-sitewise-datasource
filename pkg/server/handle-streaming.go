package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type MaybeStreamingQuery struct {
	IsStreaming       bool          `json:"isStreaming,omitempty"`
	NextToken         string        `json:"nextToken,omitempty"`
	IntervalStreaming time.Duration `json:"_"`
}

// takes a request, finds any queries that are streaming and have no next token, and stores a new request into the stream, to be consumed by runstream
func (s *Server) handleStreaming(h handler) handler {
	return func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
		respWithMaybeStreamingChannels := &backend.QueryDataResponse{
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
			respWithMaybeStreamingChannels.Responses[refID] = res

			if res.Error != nil {
				continue
			}

			originalIndividualQuery := queries[refID]

			maybeStreamingQuery := MaybeStreamingQuery{}
			if err := json.Unmarshal(originalIndividualQuery.JSON, &maybeStreamingQuery); err != nil {
				return nil, fmt.Errorf("could not determine if query was streaming or not: %w", err)
			}

			if maybeStreamingQuery.IsStreaming {
				backend.Logger.Info("found a streaming query")

				// create a channel and attach it to the response
				if respWithMaybeStreamingChannels.Responses[originalIndividualQuery.RefID].Frames[0].Meta == nil {
					respWithMaybeStreamingChannels.Responses[originalIndividualQuery.RefID].Frames[0].Meta = &data.FrameMeta{}
				}
				queryUID := uuid.New().String()
				backend.Logger.Info("creating a streaming channel", "channel", fmt.Sprintf("ds/%s/%s", s.Settings.UID, queryUID))
				respWithMaybeStreamingChannels.Responses[originalIndividualQuery.RefID].Frames[0].Meta.Channel = fmt.Sprintf("ds/%s/%s", s.Settings.UID, queryUID)

				// create a new query for the streaming channel and copy over the original query
				requestForStreaming := backend.QueryDataRequest{}
				requestForStreaming.Queries = append(requestForStreaming.Queries, originalIndividualQuery)
				requestForStreaming.Headers = req.Headers

				nextToken := ""
				if custom, ok := res.Frames[0].Meta.Custom.(map[string]interface{}); ok {
					if nt, ok := custom["NextToken"]; ok {
						nextToken = nt.(string)
					}
				}

				if nextToken != "" {
					requestForStreaming.Queries[0].JSON.MarshalJSON()
				// TODO: Shouldnt' we update the next token in request for streaming if it exists?
				// if there's no next token, update that new query to move forward in time and store that query in the stream
				if nextToken == "" {
					backend.Logger.Info("updating query to be added to stream to move forward in time")
					if ts := getFromTimestamp(res); ts != nil {
						// TODO: BUT should it be the same from??
						requestForStreaming.Queries[0].TimeRange.From = *ts
					}
					requestForStreaming.Queries[0].TimeRange.To = requestForStreaming.Queries[0].TimeRange.To.Add(5 * time.Second)
				}

				// stash the query in the stream map for use in RunStream
				s.streamMu.Lock()
				s.streams[queryUID] = requestForStreaming
				s.streamMu.Unlock()
			}
		}

		return respWithMaybeStreamingChannels, nil
	}
}

func getFromTimestamp(res backend.DataResponse) *time.Time {
	var lastTimestamp *time.Time

	for _, frame := range res.Frames {
		for _, field := range frame.Fields {
			if field.Len() == 0 {
				continue
			}

			ts := time.Unix(0, 0)
			switch field.Type() {
			case data.FieldTypeTime:
				if t, ok := field.At(field.Len() - 1).(time.Time); ok {
					ts = t
				}
			case data.FieldTypeNullableTime:
				if t, ok := field.At(field.Len() - 1).(*time.Time); ok && t != nil {
					ts = *t
				}
			}

			if lastTimestamp == nil || ts.After(*lastTimestamp) {
				lastTimestamp = &ts
			}
		}
	}

	return lastTimestamp
}
