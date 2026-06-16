package api_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeListAssetPropertiesClient struct {
	assetId string
}

func (f *fakeListAssetPropertiesClient) ListAssetProperties(ctx context.Context, input *iotsitewise.ListAssetPropertiesInput, optFns ...func(*iotsitewise.Options)) (*iotsitewise.ListAssetPropertiesOutput, error) {
	f.assetId = *input.AssetId
	retVal := iotsitewise.ListAssetPropertiesOutput{NextToken: aws.String("bar")}
	return &retVal, nil
}

func TestListAssetProperties(t *testing.T) {
	client := fakeListAssetPropertiesClient{}
	query := models.ListAssetPropertiesQuery{
		BaseQuery: models.BaseQuery{AssetIds: []string{"foo"}},
	}
	framer, err := api.ListAssetProperties(context.Background(), &client, query)
	require.NoError(t, err)
	assert.Equal(t, "foo", client.assetId)
	assert.Equal(t, "bar", *framer.NextToken)
}
