package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
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

	qualities := make([]iotsitewisetypes.Quality, 1)

	if query.Quality == "" || query.Quality == "ANY" {
		qualities[0] = iotsitewisetypes.QualityGood
	} else {
		qualities[0] = query.Quality
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	timeOrdering := iotsitewisetypes.TimeOrderingAscending
	if query.TimeOrdering != "" {
		timeOrdering = query.TimeOrdering
	}

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 250 {
		query.MaxDataPoints = 250
	}

	return &iotsitewise.GetAssetPropertyAggregatesInput{
		AggregateTypes: query.AggregateTypes,
		EndDate:        to,
		MaxResults:     aws.Int32(query.MaxDataPoints),
		NextToken:      getNextToken(query.BaseQuery),
		AssetId:        getFirstAssetId(query.BaseQuery),
		PropertyId:     getFirstPropertyId(query.BaseQuery),
		PropertyAlias:  getFirstPropertyAlias(query.BaseQuery),
		Qualities:      qualities,
		Resolution:     aws.String(resolution),
		StartDate:      from,
		TimeOrdering:   timeOrdering,
	}
}

func GetAssetPropertyAggregates(ctx context.Context, sw client.SitewiseAPIClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyAggregates, error) {

	modifiedQuery, err := getAssetIdAndPropertyId(query, sw, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReq := aggregateQueryToInput(modifiedQuery)

	resp, err := sw.GetAssetPropertyAggregatesPageAggregation(ctx, awsReq, modifiedQuery.MaxPageAggregations, int(query.MaxDataPoints))

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
