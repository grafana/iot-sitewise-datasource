package api_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeExecuteQueryClient struct {
	lastQueryStatement string
}

func (f *fakeExecuteQueryClient) ExecuteQueryWithContext(ctx aws.Context, input *iotsitewise.ExecuteQueryInput, opts ...request.Option) (*iotsitewise.ExecuteQueryOutput, error) {
	f.lastQueryStatement = aws.StringValue(input.QueryStatement)
	retVal := iotsitewise.ExecuteQueryOutput{NextToken: aws.String("bar")}
	return &retVal, nil
}

func TestExecuteQuery(t *testing.T) {
	client := fakeExecuteQueryClient{}
	query := models.ExecuteQuery{
		BaseQuery: models.BaseQuery{AssetIds: []string{"foo"}},
	}
	framer, err := api.ExecuteQuery(context.Background(), &client, query)
	require.NoError(t, err)
	assert.Equal(t, "bar", *framer.NextToken)
}
func TestExecuteQueryReceivesTheGivenSQL(t *testing.T) {
	client := &fakeExecuteQueryClient{}
	query := models.ExecuteQuery{
		BaseQuery: models.BaseQuery{AssetIds: []string{"foo"}},
		Query: sqlutil.Query{
			RawSQL: "SELECT * FROM assets",
		},
	}
	framer, err := api.ExecuteQuery(context.Background(), client, query)
	require.NoError(t, err)
	assert.Equal(t, "bar", *framer.NextToken)
	assert.Equal(t, "SELECT * FROM assets", client.lastQueryStatement)
}
