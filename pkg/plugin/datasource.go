package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
	"github.com/pkg/errors"
)

func NewSitewiseDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	serverInstance, err := createSitewiseDatasourcePluginServer(settings)
	if err != nil {
		return nil, err
	}
	return serverInstance, nil
}

type DatasourceServerInstance struct {
	GetClient clientGetterFunc
	Settings  backend.DataSourceInstanceSettings
	streamMu  sync.RWMutex
	streams   map[string]SitewiseQuery
}

func createSitewiseDatasourcePluginServer(settings backend.DataSourceInstanceSettings) (*DatasourceServerInstance, error) {
	getClient, err := createClientGetterFunc(settings)
	if err != nil {
		return nil, err
	}
	return &DatasourceServerInstance{
		GetClient: getClient,
		Settings:  settings,
		streams:   make(map[string]SitewiseQuery),
	}, nil
}

func createClientGetterFunc(settings backend.DataSourceInstanceSettings) (clientGetterFunc, error) {
	cfg := models.AWSSiteWiseDataSourceSetting{}

	err := cfg.Load(settings)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	sessions := awsds.NewSessionCache()
	clientGetter := func(region string) (swclient client.SitewiseClient, err error) {
		swclient, err = client.GetClient(region, cfg, sessions.GetSession)
		return
	}

	if cfg.Region == models.EDGE_REGION && cfg.EdgeAuthMode != models.EDGE_AUTH_MODE_DEFAULT {
		edgeAuthenticator := sitewise.EdgeAuthenticator{ //DummyAuthenticator{
			Settings: cfg,
		}

		var mu sync.Mutex
		authInfo, err := edgeAuthenticator.Authenticate()
		if err != nil {
			return nil, fmt.Errorf("error getting initial edge credentials (%s)", err.Error())
		}
		cfg.AuthType = awsds.AuthTypeKeys // Force key auth
		cfg.AccessKey = authInfo.AccessKeyId
		cfg.SecretKey = authInfo.SecretAccessKey
		cfg.SessionToken = authInfo.SessionToken

		clientGetter = func(region string) (swclient client.SitewiseClient, err error) {
			mu.Lock()
			if time.Now().After(authInfo.SessionExpiryTime) {
				log.DefaultLogger.Debug("edge credentials expired. updating credentials now.")
				authInfo, err = edgeAuthenticator.Authenticate()
				if err != nil {
					mu.Unlock()
					return nil, fmt.Errorf("error updating edge credentials (%s)", err.Error())
				}
				cfg.AccessKey = authInfo.AccessKeyId
				cfg.SecretKey = authInfo.SecretAccessKey
				cfg.SessionToken = authInfo.SessionToken
			}
			cfgCopy := cfg
			mu.Unlock()
			swclient, err = client.GetClient(region, cfgCopy, sessions.GetSession)
			return
		}
	}
	return clientGetter, nil
}

type clientGetterFunc func(region string) (client client.SitewiseClient, err error)
type invokerFunc func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error)

var (
	_ backend.CheckHealthHandler    = (*DatasourceServerInstance)(nil)
	_ backend.QueryDataHandler      = (*DatasourceServerInstance)(nil)
	_ backend.StreamHandler         = (*DatasourceServerInstance)(nil) // Streaming data source needs to implement this
	_ instancemgmt.InstanceDisposer = (*DatasourceServerInstance)(nil)
)

func (ds *DatasourceServerInstance) CheckHealth(ctx context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return nil, nil
}

/* 
Notes:
we need access to headers in our queries to know if something is from expressions, right now we're not storing headers
let's format request object into individual queries that contain all that we need to make descions

*/
func (ds *DatasourceServerInstance) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		sitewiseQuery, err := readQuery(q)
		if err != nil {
			response.Responses[q.RefID] = backend.DataResponse{
				Error: err,
			}
		}

		// make each query
		queryRes := ds.runQuery(ctx, sitewiseQuery)
		// attach it to the response
		response.Responses[q.RefID] = queryRes

		// if you're not streaming then return the data and move on to next query
		// if there's no data in this response then return the data and move on to next query,
		// apparently we don't stream when there's no data in a response?
		if !sitewiseQuery.IsStreaming || len(queryRes.Frames) == 0 {
			continue
			// TODO: double check if we need to do anything else to handle pagination
		}

		// if the results are paged, attach the next token from the response
		// to the queryMetaData so that RunStream can start from the next page
		if customMeta := loadMetaFromResponse(queryRes); customMeta != nil {
			sitewiseQuery.NextToken = customMeta.NextToken
		}

		// if we are streaming and the streaming channel does not exist yet
		// create a response channel and attach it to the query response frames as meta data
		if response.Responses[q.RefID].Frames[0].Meta == nil {
			response.Responses[q.RefID].Frames[0].Meta = &data.FrameMeta{}
		}
		queryUID := uuid.New().String()
		response.Responses[q.RefID].Frames[0].Meta.Channel = fmt.Sprintf("ds/%s/%s", ds.Settings.UID, queryUID)

		// if there's no next token, update query to move forward in time and store that query in the stream
		if sitewiseQuery.NextToken == "" {
			if ts := getFromTimestamp(queryRes); ts != nil {
				sitewiseQuery.TimeRange.From = *ts
			}
			sitewiseQuery.TimeRange.To = time.Now().Add(sitewiseQuery.IntervalStreaming)
		}

		// stash the query in the stream map for use in RunStream
		ds.streamMu.Lock()
		ds.streams[queryUID] = sitewiseQuery
		ds.streamMu.Unlock()
	}
	return response, nil
}

func (ds *DatasourceServerInstance) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	status := backend.SubscribeStreamStatusNotFound

	ds.streamMu.RLock()
	if _, ok := ds.streams[req.Path]; ok {
		status = backend.SubscribeStreamStatusOK
	}
	ds.streamMu.RUnlock()

	return &backend.SubscribeStreamResponse{
		Status: status,
	}, nil
}

func (ds *DatasourceServerInstance) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ds.streamMu.Lock()
	query, ok := ds.streams[req.Path]
	if !ok {
		ds.streamMu.Unlock()
		return fmt.Errorf("not found")
	}
	delete(ds.streams, req.Path)
	ds.streamMu.Unlock()

	resChannel := make(chan *backend.DataResponse)
	// I think this is the problem maybe? which is that we're using query
	go ds.RequestLoop(ctx, query, resChannel)

	for {
		select {
		case <-ctx.Done():
			return nil
		case res := <-resChannel:
			if res == nil {
				return nil
			}
			if res.Error != nil {
				return res.Error
			}
			for _, frame := range res.Frames {
				if err := sender.SendFrame(frame, data.IncludeAll); err != nil {
					return err
				}
			}
		}
	}
}

func (ds *DatasourceServerInstance) RequestLoop(ctx context.Context, query SitewiseQuery, resChannel chan *backend.DataResponse) {
	// stop the request loop if the context is cancelled
	select {
	case <-ctx.Done():
		resChannel <- nil
		return
	default:
	}

	res := ds.runQuery(ctx, query)
	resChannel <- &res
	if res.Error != nil {
		resChannel <- nil
		return
	}

	customMeta := loadMetaFromResponse(res)
	// if the results are paged, request the next page
	if customMeta != nil {
		query.NextToken = customMeta.NextToken
		ds.RequestLoop(ctx, query, resChannel)
		return
	}

	// we've hit the last page, and this isn't a streaming response,
	// so we can close the channel and return
	if !query.IsStreaming {
		resChannel <- nil
		return
	}

	// reset the next token for the streaming query
	query.NextToken = ""

	if ts := getFromTimestamp(res); ts != nil {
		query.TimeRange.From = *ts
	}

	time.Sleep(query.IntervalStreaming)
	query.TimeRange.To = time.Now()
	ds.RequestLoop(ctx, query, resChannel)
}

func (ds *DatasourceServerInstance) PublishStream(_ context.Context, _ *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (ds *DatasourceServerInstance) Dispose() {

}

func (ds *DatasourceServerInstance) invoke(ctx context.Context, baseQuery *models.BaseQuery, invoker invokerFunc) (data.Frames, error) {
	sw, err := ds.GetClient(baseQuery.AwsRegion)
	if err != nil {
		return nil, err
	}

	fr, err := invoker(ctx, sw)
	if err != nil {
		return nil, err
	}

	return sitewise.FrameResponse(ctx, *baseQuery, fr, sw)
}

func (ds *DatasourceServerInstance) runQuery(ctx context.Context, q SitewiseQuery) backend.DataResponse {
	remarshalledQueryJSON, err := json.Marshal(q.JSON)
	if err != nil {
		return backend.DataResponse{
			Error: fmt.Errorf("issue parsing query in run query"),
		}
	}
	bq := backend.DataQuery{
		RefID:         q.RefID,
		QueryType:     q.QueryType,
		MaxDataPoints: q.MaxDataPoints,
		Interval:      q.Interval,
		TimeRange:     q.TimeRange,
		JSON:          remarshalledQueryJSON,
	}
	req := backend.QueryDataRequest{
		Headers: [],
	}

	switch q.QueryType {
	// case models.QueryTypePropertyValueHistory:
	// 	return s.lastObservation(s.HandlePropertyValueHistory)
	case models.QueryTypePropertyAggregate:
		return ds.lastObservation(ds.handlePropertyAggregate(ctx, ))
	// case models.QueryTypePropertyInterpolated:
	// 	return s.lastObservation(s.HandleInterpolatedPropertyValue)
	case models.QueryTypePropertyValue:
		return ds.handlePropertyValue(ctx, bq)
	case models.QueryTypeListAssetModels:
		return ds.handleListAssetModels(ctx, bq)
	// case models.QueryTypeListAssociatedAssets:
	// 	return s.HandleListAssociatedAssets
	case models.QueryTypeListAssets:
		return ds.handleListAssets(ctx, bq)
	case models.QueryTypeDescribeAsset:
		return ds.handleDescribeAsset(ctx, bq)
	}
	return backend.DataResponse{
		Error: fmt.Errorf("unknown query type"),
	}
}

func (ds *DatasourceServerInstance) handleListAssets(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListAssetsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := ds.invoke(ctx, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssets(ctx, sw, *query)
	})

	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (ds *DatasourceServerInstance) handleListAssetModels(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListAssetModelsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := ds.invoke(ctx, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssetModels(ctx, sw, *query)
	})

	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (ds *DatasourceServerInstance) handlePropertyValue(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetAssetPropertyValueQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	sw, err := ds.GetClient(query.BaseQuery.AwsRegion)
	if err != nil {
		return backend.DataResponse{
			Error: errors.Wrap(err, "failed to get client"),
		}
	}

	modifiedQuery, fr, err := api.BatchGetAssetPropertyValue(ctx, sw, *query)
	if err != nil {
		return backend.DataResponse{
			Error: errors.Wrap(err, "failed to call BatchGetAssetPropertyValue"),
		}
	}

	frames, err := sitewise.FrameResponse(ctx, modifiedQuery.BaseQuery, fr, sw)
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

func (ds *DatasourceServerInstance) handleDescribeAsset(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetDescribeAssetQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := ds.invoke(ctx, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.DescribeAsset(ctx, sw, *query)
	})

	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

func (ds *DatasourceServerInstance) handleListAssociatedAssetsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query, err := models.GetListAssociatedAssetsQuery(&q)
	if err != nil {
		return DataResponseErrorUnmarshal(err)
	}

	frames, err := ds.invoke(ctx, &query.BaseQuery, func(ctx context.Context, sw client.SitewiseClient) (framer.Framer, error) {
		return api.ListAssociatedAssets(ctx, sw, *query)
	})

	if err != nil {
		return DataResponseErrorRequestFailed(err)
	}

	return backend.DataResponse{
		Frames: frames,
		Error:  nil,
	}
}

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

type SitewiseQuery struct {
	backend.DataQuery
	IsStreaming       bool                     `json:"isStreaming,omitempty"`
	NextToken         string                   `json:"nextToken,omitempty"`
	QueryType         models.SiteWiseQueryType `json:"queryType,omitempty"`
	IntervalStreaming time.Duration            `json:"_"`
}

func readQuery(query backend.DataQuery) (SitewiseQuery, error) {
	sitewiseQuery := SitewiseQuery{}
	if err := json.Unmarshal(query.JSON, &sitewiseQuery); err != nil {
		return sitewiseQuery, fmt.Errorf("could not read query: %w", err)
	}

	sitewiseQuery.RefID = query.RefID
	sitewiseQuery.MaxDataPoints = query.MaxDataPoints
	sitewiseQuery.Interval = query.Interval
	sitewiseQuery.TimeRange = query.TimeRange
	sitewiseQuery.JSON = query.JSON

	return sitewiseQuery, nil
}

type SitewiseResponseMetaData struct {
	NextToken string `json:"nextToken,omitempty"`
}

func loadMetaFromResponse(res backend.DataResponse) *SitewiseResponseMetaData {
	for _, frame := range res.Frames {
		if frame.Meta == nil || frame.Meta.Custom == nil {
			continue
		}
		meta, ok := frame.Meta.Custom.(SitewiseResponseMetaData)
		// skip frame if NextToken is not set
		if ok && meta.NextToken != "" {
			return &meta
		}
	}
	return nil
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

type handler func(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)

// used when Expand Time Range is true
// basically makes additional requests for values outside of the time range to make graphs look smoother
func (ds *DatasourceServerInstance) lastObservation(h handler) handler {
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

			lastValueRes, err := ds.lastValueQuery(ctx, query, "DESCENDING")
			if err != nil {
				log.DefaultLogger.Debug("failed to fetch last observation", "error", err)
			}
			if r, ok := resp.Responses[refID]; err == nil && ok {
				resp.Responses[refID] = mergeLastValueResponse(r, lastValueRes)
			}

			nextValueRes, err := ds.lastValueQuery(ctx, query, "ASCENDING")
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
func (ds *DatasourceServerInstance) lastValueQuery(ctx context.Context, query backend.DataQuery, timeOrdering string) (backend.DataResponse, error) {
	query.MaxDataPoints = 1
	if timeOrdering == "DESCENDING" {
		query.TimeRange.To = query.TimeRange.From.Add(-1 * time.Second)
		query.TimeRange.From = query.TimeRange.From.Add(-8760 * time.Hour) // 1 year ago
	} else if timeOrdering == "ASCENDING" {
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
