package api

import (
	"context"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
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

	qualities := make([]iotsitewisetypes.Quality, 1)
	if query.Quality == "" || query.Quality == "ANY" {
		qualities[0] = iotsitewisetypes.QualityGood
	} else {
		qualities[0] = query.Quality
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	timeOrdering := iotsitewisetypes.TimeOrderingDescending
	if query.TimeOrdering != "" {
		timeOrdering = query.TimeOrdering
	}

	entries := make([]iotsitewisetypes.BatchGetAssetPropertyAggregatesEntry, 0)

	// All unique properties are collected in AssetPropertyEntries and assigned to
	// a BatchGetAssetPropertyAggregatesEntry
	for _, entry := range query.AssetPropertyEntries {
		aggregatesEntry := iotsitewisetypes.BatchGetAssetPropertyAggregatesEntry{
			AggregateTypes: query.AggregateTypes,
			EndDate:        to,
			Qualities:      qualities,
			Resolution:     aws.String(resolution),
			StartDate:      from,
			TimeOrdering:   timeOrdering,
		}
		if entry.AssetId != "" && entry.PropertyId != "" {
			aggregatesEntry.AssetId = aws.String(entry.AssetId)
			aggregatesEntry.PropertyId = aws.String(entry.PropertyId)
			aggregatesEntry.EntryId = util.GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
		} else {
			// If there is no assetId or propertyId, then we use the propertyAlias
			aggregatesEntry.PropertyAlias = aws.String(entry.PropertyAlias)
			aggregatesEntry.EntryId = util.GetEntryIdFromPropertyAlias(entry.PropertyAlias)
		}
		entries = append(entries, aggregatesEntry)
	}

	return &iotsitewise.BatchGetAssetPropertyAggregatesInput{
		Entries: entries,
		// performance: hardcoded to fetch the maximum number of data points
		MaxResults: aws.Int32(BatchGetAssetPropertyAggregatesMaxResults),
		NextToken:  getNextToken(query.BaseQuery),
	}
}

func BatchGetAssetPropertyAggregates(ctx context.Context, client client.SitewiseAPIClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyAggregatesBatch, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	batchedQueries := batchQueries(modifiedQuery, BatchGetAssetPropertyAggregatesMaxEntries)
	requests := []iotsitewise.BatchGetAssetPropertyAggregatesInput{}
	responses := []iotsitewise.BatchGetAssetPropertyAggregatesOutput{}
	for _, q := range batchedQueries {
		awsReq := aggregateBatchQueryToInput(q)
		requests = append(requests, *awsReq)
		resp, err := client.BatchGetAssetPropertyAggregatesPageAggregation(ctx, awsReq, modifiedQuery.MaxPageAggregations, maxDps)
		if err != nil {
			return models.AssetPropertyValueQuery{}, nil, err
		}
		responses = append(responses, *resp)
	}

	return modifiedQuery,
		&framer.AssetPropertyAggregatesBatch{
			Requests:  requests,
			Responses: responses,
		}, nil
}
