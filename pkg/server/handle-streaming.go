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

type handler func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)

// finds any queries that should be streaming and stores the next streaming request in the streams map
func (s *Server) handleStreaming(h handler) handler {
	return func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
		// handle the request and get back a response regardless of if query is streaming
		response, err := h(ctx, req)
		if err != nil {
			return nil, err
		}

		queries := mapQueriesToRefId(req.Queries)

		// for each response, check if the original query is streaming,
		// if it is attach a streaming channel to the response
		// and store next query in the streams map to be processed by runStream
		for refID, res := range response.Responses {
			query := queries[refID]
			isStreaming, err := isStreaming(query)
			if err != nil {
				return nil, err
			}
			if isStreaming {
				queryUID := modifyResponseToHaveStreamingChannel(query, res, s.Settings.UID)
				nextQuery := getNextQuery(query, res, queryUID, req.Headers)
				s.streamMu.Lock()
				s.streams[queryUID] = nextQuery
				s.streamMu.Unlock()
			}
		}

		return response, nil
	}
}

func mapQueriesToRefId(unmappedQueries []backend.DataQuery) map[string]backend.DataQuery {
	queries := make(map[string]backend.DataQuery, len(unmappedQueries))
	for _, q := range unmappedQueries {
		queries[q.RefID] = q
	}
	return queries
}

type MaybeStreamingQuery struct {
	IsStreaming bool `json:"isStreaming,omitempty"`
}

func isStreaming(query backend.DataQuery) (bool, error) {
	maybeStreamingQuery := MaybeStreamingQuery{}
	if err := json.Unmarshal(query.JSON, &maybeStreamingQuery); err != nil {
		return false, fmt.Errorf("could not determine if query was streaming or not: %w", err)
	}
	return maybeStreamingQuery.IsStreaming, nil
}

func modifyResponseToHaveStreamingChannel(query backend.DataQuery, response backend.DataResponse, dsUID string) string {
	// create a channel and attach it to the response
	if response.Frames[0].Meta == nil {
		response.Frames[0].Meta = &data.FrameMeta{}
	}
	queryUID := uuid.New().String()
	response.Frames[0].Meta.Channel = fmt.Sprintf("ds/%s/%s", dsUID, queryUID)

	return queryUID
}

func getNextQuery(query backend.DataQuery, res backend.DataResponse, queryUID string, headers map[string]string) backend.QueryDataRequest {
	// create a new query request, copying over properties from the last query
	nextStreamingQuery := backend.QueryDataRequest{}
	nextStreamingQuery.Queries = append(nextStreamingQuery.Queries, query)
	nextStreamingQuery.Headers = headers

	nextToken := ""
	if custom, ok := res.Frames[0].Meta.Custom.(map[string]interface{}); ok {
		if nt, ok := custom["NextToken"]; ok {
			nextToken = nt.(string)
		}
	}

	// if there is a next token, update the query to have that next token
	if nextToken != "" {
		var oldJSON map[string]interface{}
		err := json.Unmarshal(query.JSON, &oldJSON)
		if err != nil {
			fmt.Println("Error Unmarshaling JSON:", err)
			return nextStreamingQuery
		}
		oldJSON["NextToken"] = nextToken
		newQuerysJSON, err := json.Marshal(oldJSON)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return nextStreamingQuery
		}
		nextStreamingQuery.Queries[0].JSON = newQuerysJSON
		// TODO: do we also need to delete the next token from the response??
	} else {
		// no next token means we move forward in time
		nextStreamingQuery.Queries[0].TimeRange.From = nextStreamingQuery.Queries[0].TimeRange.To
		nextStreamingQuery.Queries[0].TimeRange.To = nextStreamingQuery.Queries[0].TimeRange.To.Add(5 * time.Second)
		// TODO: rather than hardcoded 5 seconds above, change to intervalStreaming, specify that in the ui?
	}

	return nextStreamingQuery
}
