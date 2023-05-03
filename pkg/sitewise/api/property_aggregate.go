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

func aggregateQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.BatchGetAssetPropertyAggregatesInput {

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

	if query.Quality != "" && query.Quality != "ANY" {
		qualities = aws.StringSlice([]string{query.Quality})
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	if query.TimeOrdering != "" {
		timeOrdering = aws.String(query.TimeOrdering)
	}

	if query.MaxDataPoints < 1 || query.MaxDataPoints > 250 {
		query.MaxDataPoints = 250
	}

	entries := make([]*iotsitewise.BatchGetAssetPropertyAggregatesEntry, 0)

	if query.PropertyAlias != "" {
		id := getAssetId(query.BaseQuery)
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyAggregatesEntry{
			AggregateTypes: aggregateTypes,
			EndDate:        to,
			EntryId:        id,
			PropertyAlias:  getPropertyAlias(query.BaseQuery),
			Qualities:      qualities,
			Resolution:     aws.String(resolution),
			StartDate:      from,
			TimeOrdering:   timeOrdering,
		})

		return &iotsitewise.BatchGetAssetPropertyAggregatesInput{
			Entries:    entries,
			MaxResults: aws.Int64(query.MaxDataPoints),
			NextToken:  getNextToken(query.BaseQuery),
		}
	}

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

	return &iotsitewise.BatchGetAssetPropertyAggregatesInput{
		Entries:    entries,
		MaxResults: aws.Int64(query.MaxDataPoints),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func GetAssetPropertyAggregates(ctx context.Context, client client.SitewiseClient, query *models.AssetPropertyValueQuery) (*framer.AssetPropertyAggregates, error) {

	var (
		maxDps = int(query.MaxDataPoints)
	)

	err := getAndSetAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return nil, err
	}

	awsReq := aggregateQueryToInput(*query)

	resp, err := client.BatchGetAssetPropertyAggregatesPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyAggregates{
		Request: *awsReq,
		Response: iotsitewise.BatchGetAssetPropertyAggregatesOutput{
			SuccessEntries: resp.SuccessEntries,
			SkippedEntries: resp.SkippedEntries,
			ErrorEntries:   resp.ErrorEntries,
			NextToken:      resp.NextToken,
		},
	}, nil
}
