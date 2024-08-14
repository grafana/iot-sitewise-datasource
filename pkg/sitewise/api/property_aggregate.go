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
		timeOrdering   = aws.String("ASCENDING")
	)

	quality := query.Quality

	if quality == "" || quality == "ANY" {
		qualities = aws.StringSlice([]string{"GOOD"})
	} else {
		qualities = aws.StringSlice([]string{quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.TimeOrdering != "" {
		timeOrdering = aws.String(query.TimeOrdering)
	}

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 250 {
		query.MaxDataPoints = 250
	}

	return &iotsitewise.GetAssetPropertyAggregatesInput{
		AggregateTypes: aggregateTypes,
		EndDate:        to,
		MaxResults:     aws.Int64(query.MaxDataPoints),
		NextToken:      getNextToken(query.BaseQuery),
		AssetId:        getAssetId(query.BaseQuery),
		PropertyId:     getPropertyId(query.BaseQuery),
		PropertyAlias:  getPropertyAlias(query.BaseQuery),
		Qualities:      qualities,
		Resolution:     aws.String(resolution),
		StartDate:      from,
		TimeOrdering:   timeOrdering,
	}
}

func GetAssetPropertyAggregates(ctx context.Context, client client.SitewiseClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyAggregates, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReq := aggregateQueryToInput(modifiedQuery)

	resp, err := client.GetAssetPropertyAggregatesPageAggregation(ctx, awsReq, modifiedQuery.MaxPageAggregations, maxDps)

	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	return modifiedQuery,
		&framer.AssetPropertyAggregates{
			Request: *awsReq,
			Response: iotsitewise.GetAssetPropertyAggregatesOutput{
				AggregatedValues: resp.AggregatedValues,
				NextToken:        resp.NextToken,
			},
		}, nil
}
