package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerExecution(t *testing.T) {
	result := iotsitewise.ExecuteQueryOutput{NextToken: aws.String("")}
	mockSw := &mocks.SitewiseClient{}
	mockSw.On("ExecuteQueryWithContext", mock.Anything, mock.Anything).Return(&result, nil)
	mockSw.On("ListAssetsWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetsOutput{}, nil)
	mockSw.On("ListAssetModelsWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetModelsOutput{}, nil)
	mockSw.On("ListAssetPropertiesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssetPropertiesOutput{}, nil)
	mockSw.On("ListAssociatedAssetsWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.ListAssociatedAssetsOutput{}, nil)
	mockSw.On("ListTimeSeriesWithContext", mock.Anything, mock.Anything).Return(&iotsitewise.ListTimeSeriesOutput{}, nil)

	clientGetter := func(region string) (client.SitewiseClient, error) {
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
				JSON:  []byte(`{"rawSQL": "SELECT * FROM table"}`),
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
			handler_method: "ListAssetModelsWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssets",
			server_method:  "HandleListAssets",
			handler_method: "ListAssetsWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListTimeSeries",
			server_method:  "HandleListTimeSeries",
			handler_method: "ListTimeSeriesWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssetProperties",
			server_method:  "HandleListAssetProperties",
			handler_method: "ListAssetPropertiesWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ListAssociatedAssets",
			server_method:  "HandleListAssociatedAssets",
			handler_method: "ListAssociatedAssetsWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
		{
			name:           "ExecuteQuery",
			server_method:  "HandleExecuteQuery",
			handler_method: "ExecuteQueryWithContext",
			input_args:     []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)},
			called_args:    []interface{}{mock.Anything, mock.Anything},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			method := reflect.ValueOf(&server).MethodByName(tt.server_method)
			require.True(t, method.IsValid(), "method %s not found", tt.server_method)

			method.Call(tt.input_args)
			client, err := server.Datasource.GetClient("region")
			if err != nil {
				t.Fatalf("error getting client: %s", err)
			}
			mockClient := client.(*mocks.SitewiseClient)
			mockClient.AssertCalled(t, tt.handler_method, tt.called_args...)
		})
	}
}
