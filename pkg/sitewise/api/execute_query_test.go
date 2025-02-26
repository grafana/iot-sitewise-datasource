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
	executeCount       int
	lastQueryStatement string
}

func (f *fakeExecuteQueryClient) ExecuteQueryWithContext(ctx aws.Context, input *iotsitewise.ExecuteQueryInput, opts ...request.Option) (*iotsitewise.ExecuteQueryOutput, error) {
	f.executeCount++
	f.lastQueryStatement = aws.StringValue(input.QueryStatement)
	var retVal = iotsitewise.ExecuteQueryOutput{
		NextToken: aws.String("bar"),
		Rows: []*iotsitewise.Row{
			{
				Data: []*iotsitewise.Datum{
					{
						ScalarValue: aws.String("123.45"),
					},
				},
			},
		},
		Columns: []*iotsitewise.ColumnInfo{
			{
				Name: aws.String("example_column"),
				Type: &iotsitewise.ColumnType{ScalarType: aws.String("DOUBLE")},
			},
		},
	}
	if f.executeCount > 1 {
		retVal.NextToken = nil
		retVal.Rows = nil
	}
	return &retVal, nil
}

func TestExecuteQueryReturnsTheRows(t *testing.T) {
	client := fakeExecuteQueryClient{}
	query := models.ExecuteQuery{
		BaseQuery: models.BaseQuery{AssetIds: []string{"foo"}},
	}
	framer, err := api.ExecuteQuery(context.Background(), &client, query)
	require.NoError(t, err)
	assert.NotNil(t, framer.Rows)
	assert.Len(t, framer.Rows, 1)
	assert.NotNil(t, framer.Columns)
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
	assert.Equal(t, "", aws.StringValue(framer.NextToken))
	assert.Equal(t, "SELECT * FROM assets", client.lastQueryStatement)
}
