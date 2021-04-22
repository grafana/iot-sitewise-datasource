package resource

import (
	"context"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func tdpath(filename string) string {
	return "../testdata/" + filename
}

func setupMocks(t *testing.T) (*mocks.SitewiseClient, *cachingProvider) {
	client := &mocks.SitewiseClient{}
	return client, NewCachingProvider(&SitewiseResources{client})

}

func TestCachingResourceProvider(t *testing.T) {
	t.Run("testGetProperty", testGetProperty)
	t.Run("testGetAsset", testGetAsset)
}

func testGetProperty(t *testing.T) {

	mockSw, cachingProvider := setupMocks(t)
	property := testdata.GetIotSitewiseAssetProp(t, tdpath("describe-asset-property-avg-wind.json"))
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything, mock.Anything).Return(&property, nil)

	prop1, err := cachingProvider.Property(context.Background(), mock.Anything, mock.Anything)
	assert.NoError(t, err)

	newProp := testdata.GetIotSitewiseAssetProp(t, tdpath("describe-asset-property-raw-wind.json"))
	mockSw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything, mock.Anything).Return(&newProp, nil)
	prop2, err := cachingProvider.Property(context.Background(), mock.Anything, mock.Anything)
	assert.NoError(t, err)

	assert.NotEqual(t, prop2, newProp)
	assert.Equal(t, prop1, prop2)
}

func testGetAsset(t *testing.T) {
	mockSw, cachingProvider := setupMocks(t)
	asset := testdata.GetIoTSitewiseAssetDescription(t, tdpath("describe-asset.json"))
	mockSw.On("DescribeAssetWithContext", mock.Anything, mock.Anything, mock.Anything).Return(&asset, nil)

	asset1, err := cachingProvider.Asset(context.Background(), mock.Anything)
	assert.NoError(t, err)

	newAsset := testdata.GetIoTSitewiseAssetDescription(t, tdpath("describe-asset-top-level.json"))
	mockSw.On("DescribeAssetWithContext", mock.Anything, mock.Anything, mock.Anything).Return(&newAsset, nil)
	asset2, err := cachingProvider.Asset(context.Background(), mock.Anything)
	assert.NoError(t, err)

	assert.NotEqual(t, asset2, newAsset)
	assert.Equal(t, asset1, asset2)
}
