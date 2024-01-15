package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"
)

type Server struct {
	Datasource    *sitewise.Datasource
	Settings      backend.DataSourceInstanceSettings
	channelPrefix string
	closeCh       chan struct{}
	queryMux      *datasource.QueryTypeMux
	streamMu      sync.RWMutex
	streams       map[string]backend.QueryDataRequest
}

// Make sure SampleDatasource implements required interfaces.
// This is important to do since otherwise we will only get a
// not implemented error response from plugin in runtime.
var (
	_ backend.QueryDataHandler      = (*Server)(nil)
	_ backend.CheckHealthHandler    = (*Server)(nil)
	_ instancemgmt.InstanceDisposer = (*Server)(nil)
	_ backend.StreamHandler         = (*Server)(nil)
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

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func getQueryHandlers(s *Server) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(models.QueryTypePropertyValueHistory, s.handleStreaming(s.lastObservation(s.HandlePropertyValueHistory)))
	mux.HandleFunc(models.QueryTypePropertyAggregate, s.handleStreaming(s.lastObservation(s.HandlePropertyAggregate)))
	mux.HandleFunc(models.QueryTypePropertyInterpolated, s.handleStreaming(s.lastObservation(s.HandleInterpolatedPropertyValue)))
	mux.HandleFunc(models.QueryTypePropertyValue, s.handleStreaming(s.HandlePropertyValue))
	mux.HandleFunc(models.QueryTypeListAssetModels, s.handleStreaming(s.HandleListAssetModels))
	mux.HandleFunc(models.QueryTypeListAssociatedAssets, s.handleStreaming(s.HandleListAssociatedAssets))
	mux.HandleFunc(models.QueryTypeListAssets, s.handleStreaming(s.HandleListAssets))
	mux.HandleFunc(models.QueryTypeDescribeAsset, s.handleStreaming(s.HandleDescribeAsset))

	return mux
}

func NewServerInstance(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds, err := sitewise.NewDatasource(settings)
	backend.Logger.Info("SETTTTTTTTINGs", settings)
	if err != nil {
		return nil, err
	}
	srvr := &Server{
		Datasource:    ds,
		Settings:      settings,
		channelPrefix: fmt.Sprintf("ds/%d/", settings.ID),
		closeCh:       make(chan struct{}),
		streams:       make(map[string]backend.QueryDataRequest),
	}
	srvr.queryMux = getQueryHandlers(srvr) // init once
	return srvr, nil
}

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

func (s *Server) PublishStream(_ context.Context, _ *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (s *Server) handleStreamingQuery(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// there should always be one query as this was hard coded in handleStreaming or in streamingLoop
	query := req.Queries[0]
	switch query.QueryType {
	case models.QueryTypePropertyValueHistory:
		return s.lastObservation(s.HandlePropertyValueHistory)(ctx, req)
	case models.QueryTypePropertyAggregate:
		return s.lastObservation(s.HandlePropertyAggregate)(ctx, req)
	case models.QueryTypePropertyInterpolated:
		return s.lastObservation(s.HandleInterpolatedPropertyValue)(ctx, req)
	case models.QueryTypePropertyValue:
		return s.HandlePropertyValue(ctx, req)
	case models.QueryTypeListAssetModels:
		return s.HandleListAssetModels(ctx, req)
	case models.QueryTypeListAssociatedAssets:
		return s.HandleListAssociatedAssets(ctx, req)
	case models.QueryTypeListAssets:
		return s.HandleListAssets(ctx, req)
	case models.QueryTypeDescribeAsset:
		return s.HandleDescribeAsset(ctx, req)
	default:
		return nil, fmt.Errorf("unknown query type %s", query.QueryType)
	}
}

func (s *Server) streamingLoop(ctx context.Context, queryRequest *backend.QueryDataRequest, resChannel chan *backend.QueryDataResponse) {
	select {
	case <-ctx.Done():
		backend.Logger.Info("context closed in streaming loop?")
		resChannel <- nil
		return
	default:
	}

	// get the last query from the stream, and then send the response to the res channel
	res, err := s.handleStreamingQuery(ctx, queryRequest)

	if err != nil {
		backend.Logger.Info("got a error", err)
		resChannel <- nil
		return
	}
	resChannel <- res

	// if the results are paged, request the next page
	nextToken := findNextToken(*res)
	if nextToken != "" {
		lastQueriesJson := queryRequest.Queries[0].JSON
		var data map[string]interface{}
		err := json.Unmarshal(lastQueriesJson, &data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		data["NextToken"] = nextToken
		newQuerysJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		newQueries := make([]backend.DataQuery, len(queryRequest.Queries))
		newQueries[0] = queryRequest.Queries[0]
		newQueries[0].JSON = newQuerysJSON
		newRequest := backend.QueryDataRequest{
			PluginContext: queryRequest.PluginContext,
			Headers:       queryRequest.Headers,
			Queries:       newQueries,
		}

		s.streamingLoop(ctx, &newRequest, resChannel)
		return
	}

	// the results are either not paged, or we're on the last page. So now we create a query that moves forward in time
	newRequest := backend.QueryDataRequest{
		PluginContext: queryRequest.PluginContext,
		Headers:       queryRequest.Headers,
		Queries:       queryRequest.Queries,
	}
	refId := queryRequest.Queries[0].RefID
	resValue := res.Responses[refId]
	if ts := getFromTimestamp(resValue); ts != nil {
		newRequest.Queries[0].TimeRange.From = *ts
	}
	// TODO figure out interval streaming
	// newRequest.Queries[0].TimeRange.To = time.Now().Add(queryRequest.IntervalStreaming)
	time.Sleep(5 * time.Second)
	newRequest.Queries[0].TimeRange.To = newRequest.Queries[0].TimeRange.To.Add(5 * time.Second)
	s.streamingLoop(ctx, &newRequest, resChannel)
}

func (s *Server) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.streamMu.Lock()
	storedRequest, ok := s.streams[req.Path]
	if !ok {
		s.streamMu.Unlock()
		return fmt.Errorf("not found")
	}
	delete(s.streams, req.Path)
	s.streamMu.Unlock()

	resChannel := make(chan *backend.QueryDataResponse)

	go s.streamingLoop(ctx, &storedRequest, resChannel)

	for {
		select {
		// when context is cancelled, close channel?
		case <-ctx.Done():
			return nil
		// when we get a new response in res channel
		case qdr := <-resChannel:
			// and it's nothing, close channel?
			if qdr == nil {
				return nil
			}
			// again always assume there is one query
			res := qdr.Responses[storedRequest.Queries[0].RefID]
			// or it's an error, return an error
			if res.Error != nil {
				return res.Error
			}
			// otherwise send each frame
			for _, frame := range res.Frames {
				if err := sender.SendFrame(frame, data.IncludeAll); err != nil {
					return err
				}
			}
		}
	}
}

func (s *Server) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	status := backend.SubscribeStreamStatusNotFound

	s.streamMu.RLock()
	if _, ok := s.streams[req.Path]; ok {
		status = backend.SubscribeStreamStatusOK
	}
	s.streamMu.RUnlock()

	return &backend.SubscribeStreamResponse{
		Status: status,
	}, nil
}

type FrameWithCustomMeta struct {
	NextToken string `json:"nextToken,omitempty"`
}

// feels like a lot of work? should we just store the query and ehaders and not a response which has responses
func findNextToken(res backend.QueryDataResponse) string {
	for _, res := range res.Responses {
		for _, frame := range res.Frames {
			if frame.Meta == nil || frame.Meta.Custom == nil {
				continue
			}
			meta, ok := frame.Meta.Custom.(FrameWithCustomMeta)
			// skip frame if NextToken is not set
			if ok && meta.NextToken != "" {
				return meta.NextToken
			}
		}
	}

	return ""
}
