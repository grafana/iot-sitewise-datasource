package resource

import (
	"context"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func tdpath(filename string) string {
	return "../testdata/" + filename
}

func setupMocks() (*mocks.SitewiseAPIClient, *CachingResourceProvider) {
	client := &mocks.SitewiseAPIClient{}
	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)
	return client, NewCachingResourceProvider(&SitewiseResources{client}, c)
}

func TestCachingResourceProvider(t *testing.T) {
	t.Run("testGetProperty", testGetProperty)
	t.Run("testGetAsset", testGetAsset)
	t.Run("testGetAssetModel", testGetAssetModel)
	t.Run("testGetPropertyWithAlias", testGetPropertyWithAlias)
	t.Run("testAssetError", testAssetError)
	t.Run("testPropertyError", testPropertyError)
	t.Run("testAssetModelError", testAssetModelError)
}

func testGetProperty(t *testing.T) {

	mockSw, cachingProvider := setupMocks()
	property := testdata.GetIotSitewiseAssetProp(t, tdpath("describe-asset-property-avg-wind.json"))
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything, mock.Anything).
		Return(&property, nil).
		Once()

	prop1, err := cachingProvider.Property(context.Background(), mock.Anything, mock.Anything, mock.Anything)
	assert.NoError(t, err)

	prop2, err := cachingProvider.Property(context.Background(), mock.Anything, mock.Anything, mock.Anything)
	assert.NoError(t, err)

	assert.Equal(t, prop1, prop2)
	mockSw.AssertExpectations(t)
}

func testGetAsset(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	asset := testdata.GetIoTSitewiseAssetDescription(t, tdpath("describe-asset.json"))
	mockSw.On("DescribeAsset", mock.Anything, mock.Anything, mock.Anything).
		Return(&asset, nil).
		Once()

	asset1, err := cachingProvider.Asset(context.Background(), mock.Anything)
	assert.NoError(t, err)

	asset2, err := cachingProvider.Asset(context.Background(), mock.Anything)
	assert.NoError(t, err)

	assert.Equal(t, asset1, asset2)
	mockSw.AssertExpectations(t)
}

func testGetAssetModel(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	assetModel := testdata.GetIoTSitewiseAssetModelDescription(t, tdpath("describe-asset-model.json"))
	mockSw.On("DescribeAssetModel", mock.Anything, mock.Anything, mock.Anything).
		Return(&assetModel, nil).
		Once()

	model1, err := cachingProvider.AssetModel(context.Background(), mock.Anything)
	assert.NoError(t, err)

	model2, err := cachingProvider.AssetModel(context.Background(), mock.Anything)
	assert.NoError(t, err)

	assert.Equal(t, model1, model2)
	mockSw.AssertExpectations(t)
}

func testGetPropertyWithAlias(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	property := testdata.GetIotSitewiseAssetProp(t, tdpath("describe-asset-property-avg-wind.json"))
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything, mock.Anything).
		Return(&property, nil).
		Once()

	// First call with alias - should cache with alias as key
	prop1, err := cachingProvider.Property(context.Background(), "asset123", "prop456", "/alias/test")
	assert.NoError(t, err)
	assert.NotNil(t, prop1)

	// Second call with same alias - should return cached value without calling API again
	prop2, err := cachingProvider.Property(context.Background(), "asset123", "prop456", "/alias/test")
	assert.NoError(t, err)

	assert.Equal(t, prop1, prop2)
	mockSw.AssertExpectations(t)
}

func testAssetError(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	mockSw.On("DescribeAsset", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	asset, err := cachingProvider.Asset(context.Background(), "asset123")
	assert.Error(t, err)
	assert.Nil(t, asset)
}

func testPropertyError(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	mockSw.On("DescribeAssetProperty", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	prop, err := cachingProvider.Property(context.Background(), "asset123", "prop456", "")
	assert.Error(t, err)
	assert.Nil(t, prop)
}

func testAssetModelError(t *testing.T) {
	mockSw, cachingProvider := setupMocks()
	mockSw.On("DescribeAssetModel", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	model, err := cachingProvider.AssetModel(context.Background(), "model123")
	assert.Error(t, err)
	assert.Nil(t, model)
}
