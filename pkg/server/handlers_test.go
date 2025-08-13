package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerExecution(t *testing.T) {
	result := iotsitewise.ExecuteQueryOutput{NextToken: aws.String("")}
	mockSw := &mocks.SitewiseAPIClient{}
	mockSw.On("ExecuteQuery", mock.Anything, mock.Anything).Return(&result, nil)
	mockSw.On("ListAssets", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetsOutput{}, nil)
	mockSw.On("ListAssetModels", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetModelsOutput{}, nil)
	mockSw.On("ListAssetProperties", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetPropertiesOutput{}, nil)
	mockSw.On("ListAssociatedAssets", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssociatedAssetsOutput{}, nil)
	mockSw.On("ListTimeSeries", mock.Anything, mock.Anything).Return(&iotsitewise.ListTimeSeriesOutput{}, nil)

	clientGetter := func(context.Context, string) (client.SitewiseAPIClient, error) {
		return mockSw, nil
	}

	server := Server{
		Datasource: &sitewise.Datasource{
			GetClient: clientGetter,
		},
	}
	ctx := context.Background()
	req := &backend.QueryDataRequest{
		Queries: []backend.DataQuery{
			{
				RefID: "A",
				JSON:  []byte(`{"assetIds": ["asset-1"], "rawSQL": "SELECT * FROM table"}`),
			},
		},
	}

	tests := []struct {
		name           string
		server_method  string
		handler_method string
		input_args     []reflect.Value
		called_args    []interface{}
	}{
		{
			name:           "ListAssetModels",
			server_method:  "HandleListAssetModels",
			handler_method: "ListAssetModels",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssets",
			server_method:  "HandleListAssets",
			handler_method: "ListAssets",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListTimeSeries",
			server_method:  "HandleListTimeSeries",
			handler_method: "ListTimeSeries",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssetProperties",
			server_method:  "HandleListAssetProperties",
			handler_method: "ListAssetProperties",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssociatedAssets",
			server_method:  "HandleListAssociatedAssets",
			handler_method: "ListAssociatedAssets",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ExecuteQuery",
			server_method:  "HandleExecuteQuery",
			handler_method: "ExecuteQuery",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			method := reflect.ValueOf(&server).MethodByName(tt.server_method)
			require.True(t, method.IsValid(), "method %s not found", tt.server_method)

			method.Call(tt.input_args)
			client, err := server.Datasource.GetClient(context.Background(), "region")
			if err != nil {
				t.Fatalf("error getting client: %s", err)
			}
			mockClient := client.(*mocks.SitewiseAPIClient)
			mockClient.AssertCalled(t, tt.handler_method, tt.called_args...)
		})
	}
}

func TestProcessQueries_SetsDefaultRegion(t *testing.T) {
	ctx := context.Background()
	defaultRegion := "us-east-2"

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no region set",
			input:    `{"assetIds": ["asset-1"]}`,
			expected: defaultRegion,
		},
		{
			name:     "empty region",
			input:    `{"region": "", "assetIds": ["asset-1"]}`,
			expected: defaultRegion,
		},
		{
			name:     "default region",
			input:    `{"region": "default", "assetIds": ["asset-1"]}`,
			expected: defaultRegion,
		},
		{
			name:     "non-empty region is preserved",
			input:    `{"region": "eu-west-2", "assetIds": ["asset-1"]}`,
			expected: "eu-west-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedQuery backend.DataQuery
			handler := func(_ context.Context, _ *backend.QueryDataRequest, q backend.DataQuery) backend.DataResponse {
				capturedQuery = q
				return backend.DataResponse{}
			}

			req := &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{
					DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
						JSONData: []byte(fmt.Sprintf(`{"defaultRegion": "%s"}`, defaultRegion)),
					},
				},
				Queries: []backend.DataQuery{
					{
						RefID: "A",
						JSON:  []byte(tt.input),
					},
				},
			}

			processQueries(ctx, req, handler)

			var queryMap map[string]interface{}
			err := json.Unmarshal(capturedQuery.JSON, &queryMap)
			require.NoError(t, err)

			require.Equal(t, tt.expected, queryMap["region"])
		})
	}
}
