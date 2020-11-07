package sitewise

import (
	"context"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func getBestResolution(interval time.Duration) string {
	if interval < time.Minute*30 {
		return "1m"
	}
	if interval < time.Hour*20 {
		return "1h"
	}
	return "1d" // largest interval
}

func aggregateQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyAggregatesInput {

	resolution := query.Resolution
	if resolution == "AUTO" {
		resolution = getBestResolution(query.Interval)
	}

	var (
		propertyId     *string
		assetId        *string
		nextToken      *string
		aggregateTypes = aws.StringSlice(query.AggregateTypes)
		qualities      []*string
	)

	assetId = getAssetId(query.BaseQuery)
	propertyId = getPropertyId(query.BaseQuery)
	nextToken = getNextToken(query.BaseQuery)

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	return &iotsitewise.GetAssetPropertyAggregatesInput{
		AggregateTypes: aggregateTypes,
		AssetId:        assetId,
		EndDate:        to,
		MaxResults:     aws.Int64(250),
		NextToken:      nextToken,
		PropertyId:     propertyId,
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
		AggregatedValues: resp.AggregatedValues,
		NextToken:        resp.NextToken,
	}, nil
}
