package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
	"golang.org/x/sync/errgroup"
)

var (
	LOCF_INTERPOLATION   string = "LOCF_INTERPOLATION"
	LINEAR_INTERPOLATION string = "LINEAR_INTERPOLATION"
)

type responseWrapper struct {
	DataResponse *iotsitewise.GetInterpolatedAssetPropertyValuesOutput
	EntryId      string
}

func interpolatedQueryToInputs(query models.AssetPropertyValueQuery) []*iotsitewise.GetInterpolatedAssetPropertyValuesInput {

	from, to := util.TimeRangeToUnix(query.TimeRange)
	startTimeInSeconds := from.Unix()
	endTimeInSeconds := to.Unix()

	quality := query.Quality
	if quality == "" || quality == "ANY" {
		quality = "GOOD"
	}

	interpolationType := LINEAR_INTERPOLATION

	intervalInSeconds := int64(propvals.ResolutionToDuration(propvals.InterpolatedResolution(query)).Seconds())
	if query.Resolution != "AUTO" && query.Resolution != "" {
		intervalInSeconds = int64(propvals.ResolutionToDuration(query.Resolution).Seconds())
	}

	if intervalInSeconds > (endTimeInSeconds - startTimeInSeconds) {
		intervalInSeconds = endTimeInSeconds - startTimeInSeconds
	}

	if intervalInSeconds < 1 {
		intervalInSeconds = 1
	}

	awsReqs := make([]*iotsitewise.GetInterpolatedAssetPropertyValuesInput, 0)

	// All unique properties are collected in AssetPropertyEntries and used in
	// separate GetInterpolatedAssetPropertyValues requests
	for _, entry := range query.AssetPropertyEntries {
		interpolatedInput := iotsitewise.GetInterpolatedAssetPropertyValuesInput{
			StartTimeInSeconds: &startTimeInSeconds,
			EndTimeInSeconds:   &endTimeInSeconds,
			IntervalInSeconds:  aws.Int64(intervalInSeconds),
			MaxResults:         aws.Int64(10),
			Quality:            &quality,
			Type:               &interpolationType,
		}
		var entryId *string
		if entry.AssetId != "" && entry.PropertyId != "" {
			interpolatedInput.AssetId = aws.String(entry.AssetId)
			interpolatedInput.PropertyId = aws.String(entry.PropertyId)
			entryId = util.GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
		} else {
			// If there is no assetId or propertyId, then we use the propertyAlias
			interpolatedInput.PropertyAlias = aws.String(entry.PropertyAlias)
			entryId = util.GetEntryIdFromPropertyAlias(entry.PropertyAlias)
		}
		var nextToken *string
		token, ok := query.NextTokens[*entryId]
		if ok {
			nextToken = aws.String(token)
		}
		interpolatedInput.NextToken = nextToken
		awsReqs = append(awsReqs, &interpolatedInput)
	}

	return awsReqs
}

func GetInterpolatedAssetPropertyValues(ctx context.Context, client client.SitewiseClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.InterpolatedAssetPropertyValue, error) {
	maxDps := int(query.MaxDataPoints)

	modifiedQuery, err := getAssetIdAndPropertyId(query, client, ctx)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	awsReqs := interpolatedQueryToInputs(modifiedQuery)

	resultChan := make(chan *responseWrapper, len(awsReqs))
	eg, ectx := errgroup.WithContext(ctx)
	for _, req := range awsReqs {
		awsReq := req
		eg.Go(func() error {
			resp, err := client.GetInterpolatedAssetPropertyValuesPageAggregation(ectx, awsReq, query.MaxPageAggregations, maxDps)
			if err != nil {
				return err
			}
			entryId := ""
			if awsReq.AssetId != nil && awsReq.PropertyId != nil {
				entryId = *util.GetEntryIdFromAssetProperty(*awsReq.AssetId, *awsReq.PropertyId)
			} else {
				entryId = *util.GetEntryIdFromPropertyAlias(*awsReq.PropertyAlias)
			}
			resultChan <- &responseWrapper{
				DataResponse: resp,
				EntryId:      entryId,
			}

			return nil
		})
	}

	err = eg.Wait()
	close(resultChan)
	if err != nil {
		return models.AssetPropertyValueQuery{}, nil, err
	}

	responses := make(map[string]*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, len(awsReqs))
	for result := range resultChan {
		responses[result.EntryId] = result.DataResponse
	}

	return modifiedQuery,
		&framer.InterpolatedAssetPropertyValue{
			Responses: responses,
			Query:     modifiedQuery,
		}, nil
}
