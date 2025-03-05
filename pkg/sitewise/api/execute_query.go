package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func ExecuteQuery(ctx context.Context, client client.ExecuteQueryClient, query models.ExecuteQuery) (*framer.QueryResults, error) {
	backend.Logger.FromContext(ctx).Debug("Running ExecuteQuery", "query", query.RawSQL)
	input := &iotsitewise.ExecuteQueryInput{
		QueryStatement: aws.String(query.RawSQL),
		MaxResults:     aws.Int64(2000),
	}

	backend.Logger.FromContext(ctx).Debug("Beginning the query loop")

	resp, err := client.ExecuteQueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return &framer.QueryResults{
		Rows:      resp.Rows,
		Columns:   resp.Columns,
		NextToken: resp.NextToken,
	}, nil
}
