package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func aggregateQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyAggregatesInput {

	resolution := query.Resolution
	if resolution == "AUTO" {
		resolution = propvals.Resolution(query.BaseQuery)
		if resolution == propvals.ResolutionSecond {
			// override with 1m until 1s resolution is supported
			resolution = propvals.ResolutionMinute
		}
	}

	var (
		aggregateTypes = aws.StringSlice(query.AggregateTypes)
		qualities      []*string
	)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	return &iotsitewise.GetAssetPropertyAggregatesInput{
		AggregateTypes: aggregateTypes,
		EndDate:        to,
		MaxResults:     aws.Int64(250),
		NextToken:      getNextToken(query.BaseQuery),
		AssetId:        getAssetId(query.BaseQuery),
		PropertyId:     getPropertyId(query.BaseQuery),
		PropertyAlias:  getPropertyAlias(query.BaseQuery),
		Qualities:      qualities,
		Resolution:     aws.String(resolution),
		StartDate:      from,
	}
}

func GetAssetPropertyAggregates(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyAggregates, error) {

	var (
		maxDps = int(query.MaxDataPoints)
	)

	awsReq := aggregateQueryToInput(query)

	resp, err := client.GetAssetPropertyAggregatesPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyAggregates{
		Request: *awsReq,
		Response: iotsitewise.GetAssetPropertyAggregatesOutput{
			AggregatedValues: resp.AggregatedValues,
			NextToken:        resp.NextToken,
		},
	}, nil
}
