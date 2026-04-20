package sitewise

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequiresJsonParsing(t *testing.T) {
	tests := []struct {
		name      string
		queryType string
		expected  bool
	}{
		{
			name:      "PropertyValueHistory requires parsing",
			queryType: models.QueryTypePropertyValueHistory,
			expected:  true,
		},
		{
			name:      "PropertyAggregate requires parsing",
			queryType: models.QueryTypePropertyAggregate,
			expected:  true,
		},
		{
			name:      "PropertyValue requires parsing",
			queryType: models.QueryTypePropertyValue,
			expected:  true,
		},
		{
			name:      "ListAssets does not require parsing",
			queryType: models.QueryTypeListAssets,
			expected:  false,
		},
		{
			name:      "ListAssociatedAssets does not require parsing",
			queryType: models.QueryTypeListAssociatedAssets,
			expected:  false,
		},
		{
			name:      "Empty query type does not require parsing",
			queryType: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := models.BaseQuery{QueryType: tt.queryType}
			result := requiresJsonParsing(query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseJSONFields_MultipleDiagnosticsFromDifferentAssets(t *testing.T) {
	ctx := context.Background()
	mockClient := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)

	testAsset1 := &iotsitewise.DescribeAssetOutput{
		AssetName: aws.String("Generator"),
		AssetProperties: []types.AssetProperty{
			{
				Id:   aws.String("prop-123"),
				Name: aws.String("Vibration"),
			},
			{
				Id:   aws.String("prop-456"),
				Name: aws.String("Voltage"),
			},
		},
	}

	testAsset2 := &iotsitewise.DescribeAssetOutput{
		AssetName: aws.String("AirCompressor"),
		AssetProperties: []types.AssetProperty{
			{
				Id:   aws.String("prop-789"),
				Name: aws.String("Pressure"),
			},
			{
				Id:   aws.String("prop-012"),
				Name: aws.String("Temperature"),
			},
		},
	}

	mockClient.On("DescribeAsset", mock.Anything, mock.MatchedBy(func(input *iotsitewise.DescribeAssetInput) bool {
		return *input.AssetId == "asset-123"
	}), mock.Anything).Return(testAsset1, nil)

	mockClient.On("DescribeAsset", mock.Anything, mock.MatchedBy(func(input *iotsitewise.DescribeAssetInput) bool {
		return *input.AssetId == "asset-456"
	}), mock.Anything).Return(testAsset2, nil)

	resources := resource.NewCachingResourceProvider(resource.NewSitewiseResources(mockClient), c)

	jsonStr := `{
		"timestamp": "2026-02-20T22:30:00.000000",
		"prediction": 1,
		"prediction_reason": "ANOMALY_DETECTED",
		"anomaly_score": 0.81356,
		"diagnostics": [
			{"name": "asset-123\\prop-123", "value": 0.2847},
			{"name": "asset-123\\prop-456", "value": 0.1923},
			{"name": "asset-456\\prop-789", "value": 0.3562},
			{"name": "asset-456\\prop-012", "value": 0.1668}
		]
	}`

	frame := data.NewFrame("test",
		data.NewField("data", nil, []string{jsonStr}),
	)

	result := ParseJSONFields(ctx, data.Frames{frame}, resources)
	assert.Len(t, result, 1)

	resultFieldMap := make(map[string]*data.Field)
	for _, field := range result[0].Fields {
		resultFieldMap[field.Name] = field
	}

	expectedFields := map[string]float64{
		"contrib_Generator_Vibration":       0.2847,
		"contrib_Generator_Voltage":         0.1923,
		"contrib_AirCompressor_Pressure":    0.3562,
		"contrib_AirCompressor_Temperature": 0.1668,
	}

	for fieldName, expectedValue := range expectedFields {
		field, exists := resultFieldMap[fieldName]
		assert.True(t, exists, "Expected field %s to exist", fieldName)
		if exists {
			assert.Equal(t, expectedValue, field.At(0), "Field %s should have value %f", fieldName, expectedValue)
		}
	}

	assert.NotNil(t, resultFieldMap["prediction"])
	assert.Equal(t, 1.0, resultFieldMap["prediction"].At(0))

	assert.NotNil(t, resultFieldMap["prediction_reason"])
	assert.Equal(t, "ANOMALY_DETECTED", resultFieldMap["prediction_reason"].At(0))

	assert.NotNil(t, resultFieldMap["anomaly_score"])
	assert.Equal(t, 0.81356, resultFieldMap["anomaly_score"].At(0))

	mockClient.AssertExpectations(t)
}

func TestParseJSONFields_NonStringField(t *testing.T) {
	ctx := context.Background()
	mockClient := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
	resources := resource.NewCachingResourceProvider(resource.NewSitewiseResources(mockClient), c)

	frame := data.NewFrame("test",
		data.NewField("number", nil, []float64{1.0, 2.0, 3.0}),
	)

	result := ParseJSONFields(ctx, data.Frames{frame}, resources)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].Fields, 1)
	assert.Equal(t, "number", result[0].Fields[0].Name)
}

func TestParseJSONFields_EmptyStringField(t *testing.T) {
	ctx := context.Background()
	mockClient := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
	resources := resource.NewCachingResourceProvider(resource.NewSitewiseResources(mockClient), c)

	frame := data.NewFrame("test",
		data.NewField("data", nil, []string{}),
	)

	result := ParseJSONFields(ctx, data.Frames{frame}, resources)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].Fields, 1)
}

func TestParseJSONFields_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockClient := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
	resources := resource.NewCachingResourceProvider(resource.NewSitewiseResources(mockClient), c)

	frame := data.NewFrame("test",
		data.NewField("data", nil, []string{"not json", "also not json"}),
	)

	result := ParseJSONFields(ctx, data.Frames{frame}, resources)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].Fields, 1)
	assert.Equal(t, "data", result[0].Fields[0].Name)
}

func TestParseJSONFields_PreservesFrameMeta(t *testing.T) {
	ctx := context.Background()
	mockClient := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
	resources := resource.NewCachingResourceProvider(resource.NewSitewiseResources(mockClient), c)

	jsonStr := `{"timestamp": "2026-02-20T22:30:00.000000", "prediction": 0, "prediction_reason": "NO_ANOMALY_DETECTED", "anomaly_score": 0.0, "diagnostics": []}`
	frame := data.NewFrame("test",
		data.NewField("data", nil, []string{jsonStr}),
	)
	frame.Meta = &data.FrameMeta{
		Custom: map[string]interface{}{"key": "value"},
	}

	result := ParseJSONFields(ctx, data.Frames{frame}, resources)

	assert.NotNil(t, result[0].Meta)
	customMap, ok := result[0].Meta.Custom.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", customMap["key"])
}
