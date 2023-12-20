package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func ExecuteQuery(ctx context.Context, client client.SitewiseClient, query models.ExecuteQuery) (*framer.QueryResults, error) {
	awsReq := &iotsitewise.ExecuteQueryInput{QueryStatement: &query.QueryStatement}

	resp, err := client.ExecuteQueryWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.QueryResults{
		Columns:   resp.Columns,
		Rows:      resp.Rows,
		NextToken: resp.NextToken,
	}, nil
}
