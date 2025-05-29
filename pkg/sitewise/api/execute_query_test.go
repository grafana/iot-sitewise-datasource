package api_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

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

func (f *fakeExecuteQueryClient) ExecuteQuery(_ context.Context, input *iotsitewise.ExecuteQueryInput, _ ...func(*iotsitewise.Options)) (*iotsitewise.ExecuteQueryOutput, error) {
	f.executeCount++
	f.lastQueryStatement = *input.QueryStatement
	var retVal = iotsitewise.ExecuteQueryOutput{
		NextToken: aws.String("next-token"),
		Rows: []iotsitewisetypes.Row{
			{
				Data: []iotsitewisetypes.Datum{
					{
						ScalarValue: aws.String("123.45"),
					},
				},
			},
		},
		Columns: []iotsitewisetypes.ColumnInfo{
			{
				Name: aws.String("example_column"),
				Type: &iotsitewisetypes.ColumnType{ScalarType: iotsitewisetypes.ScalarTypeDouble},
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
	assert.Equal(t, "next-token", *framer.NextToken)
	assert.Equal(t, "SELECT * FROM assets", client.lastQueryStatement)
}
