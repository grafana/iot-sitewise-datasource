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

// `query.MaxDataPoints` is ignored and it always requests with the maximum number of data points the SiteWise API can support
func aggregateBatchQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyAggregatesInput {

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

	entries := make([]*iotsitewise.BatchGetAssetPropertyAggregatesEntry, 0)

	switch {
	case query.PropertyAlias != "":
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyAggregatesEntry{
			AggregateTypes: aggregateTypes,
			EndDate:        to,
			EntryId:        util.GetEntryId(query.BaseQuery),
			PropertyAlias:  util.GetPropertyAlias(query.BaseQuery),
			Qualities:      qualities,
			Resolution:     aws.String(resolution),
			StartDate:      from,
			TimeOrdering:   timeOrdering,
		})
	default:
		for _, assetId := range query.AssetIds {
			var id *string
			if assetId != "" {
				id = aws.String(assetId)
			}
			entries = append(entries, &iotsitewise.BatchGetAssetPropertyAggregatesEntry{
				AggregateTypes: aggregateTypes,
				EndDate:        to,
				EntryId:        id,
				AssetId:        id,
				PropertyId:     aws.String(query.PropertyId),
				Qualities:      qualities,
				Resolution:     aws.String(resolution),
				StartDate:      from,
				TimeOrdering:   timeOrdering,
			})
		}
	}

	return &iotsitewise.BatchGetAssetPropertyAggregatesInput{
		Entries: entries,
		// performance: hardcoded to fetch the maximum number of data points
		MaxResults: aws.Int64(BatchGetAssetPropertyAggregatesMaxResults),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyAggregates(ctx context.Context, client client.SitewiseClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyAggregatesBatch, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReq := aggregateBatchQueryToInput(modifiedQuery)

	resp, err := client.BatchGetAssetPropertyAggregatesPageAggregation(ctx, awsReq, modifiedQuery.MaxPageAggregations, maxDps)

	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	return modifiedQuery,
		&framer.AssetPropertyAggregatesBatch{
			Request: *awsReq,
			Response: iotsitewise.BatchGetAssetPropertyAggregatesOutput{
				SuccessEntries: resp.SuccessEntries,
				SkippedEntries: resp.SkippedEntries,
				ErrorEntries:   resp.ErrorEntries,
				NextToken:      resp.NextToken,
			},
		}, nil
}
