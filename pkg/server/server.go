package server

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"
)

type Server struct {
	datasource Datasource
}

// QueryHandlerFunc is the function signature used for mux.HandleFunc
// Looks like mux.HandleFunc uses backend.QueryHandlerFunc
// type QueryDataHandlerFunc func(ctx context.Context, req *QueryDataRequest) (*QueryDataResponse, error)
type QueryHandlerFunc func(context.Context, *backend.QueryDataRequest, backend.DataQuery) backend.DataResponse

func processQueries(ctx context.Context, req *backend.QueryDataRequest, handler QueryHandlerFunc) *backend.QueryDataResponse {
	res := backend.Responses{}
	for _, v := range req.Queries {
		res[v.RefID] = handler(ctx, req, v)
	}

	return &backend.QueryDataResponse{
		Responses: res,
	}

}

// UnmarshalQuery attempts to unmarshal a query from JSON
//func UnmarshalQuery(b []byte, v interface{}) *backend.DataResponse {
//	if err := json.Unmarshal(b, v); err != nil {
//		return DataResponseError(err, "failed to unmarshal JSON request into query")
//	}
//	return nil
//}

func DataResponseError(err error, message string) backend.DataResponse {
	return backend.DataResponse{
		Error: errors.Wrap(err, message),
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

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func GetQueryHandlers(s *Server) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(models.QueryTypePropertyValueHistory, s.HandlePropertyValueHistory)
	mux.HandleFunc(models.QueryTypeListAssetModels, s.HandleListAssetModels)
	mux.HandleFunc(models.QueryTypeListAssets, s.HandleListAssets)

	return mux
}

func NewServerInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	srvr := &Server{
		datasource: sitewise.NewDatasource(),
	}
	return srvr, nil
}
