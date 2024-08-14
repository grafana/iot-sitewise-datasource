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
	//if propertyAlias is set make sure to set the assetId and propertyId to nil
	if query.PropertyAlias != "" {
		query.PropertyId = ""
		query.AssetIds = []string{}
		// nolint:staticcheck
		query.AssetId = ""
	}

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

	awsReqsLen := len(query.AssetIds)
	if query.PropertyAlias != "" {
		awsReqsLen = 1
	}
	awsReqs := make([]*iotsitewise.GetInterpolatedAssetPropertyValuesInput, awsReqsLen)
	if query.PropertyAlias != "" {
		var nextToken *string
		token, ok := query.NextTokens[query.PropertyAlias]
		if ok {
			nextToken = aws.String(token)
		}
		awsReqs[0] = &iotsitewise.GetInterpolatedAssetPropertyValuesInput{
			StartTimeInSeconds: &startTimeInSeconds,
			EndTimeInSeconds:   &endTimeInSeconds,
			IntervalInSeconds:  aws.Int64(intervalInSeconds),
			MaxResults:         aws.Int64(10),
			NextToken:          nextToken,
			AssetId:            util.GetAssetId(query.BaseQuery),
			PropertyId:         util.GetPropertyId(query.BaseQuery),
			PropertyAlias:      util.GetPropertyAlias(query.BaseQuery),
			Quality:            &quality,
			Type:               &interpolationType,
		}
		return awsReqs
	} else {
		for idx, assetId := range query.AssetIds {
			var nextToken *string
			token, ok := query.NextTokens[assetId]
			if ok {
				nextToken = aws.String(token)
			}
			awsReqs[idx] = &iotsitewise.GetInterpolatedAssetPropertyValuesInput{
				StartTimeInSeconds: &startTimeInSeconds,
				EndTimeInSeconds:   &endTimeInSeconds,
				IntervalInSeconds:  aws.Int64(intervalInSeconds),
				MaxResults:         aws.Int64(10),
				NextToken:          nextToken,
				AssetId:            aws.String(assetId),
				PropertyId:         util.GetPropertyId(query.BaseQuery),
				PropertyAlias:      util.GetPropertyAlias(query.BaseQuery),
				Quality:            &quality,
				Type:               &interpolationType,
			}
		}
	}

	return awsReqs
}

func GetInterpolatedAssetPropertyValues(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.InterpolatedAssetPropertyValue, error) {
	maxDps := int(query.MaxDataPoints)
	awsReqs := interpolatedQueryToInputs(query)

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
			if awsReq.AssetId != nil {
				entryId = *awsReq.AssetId
			} else {
				entryId = *awsReq.PropertyAlias
			}
			resultChan <- &responseWrapper{
				DataResponse: resp,
				EntryId:      entryId,
			}

			return nil
		})
	}

	err := eg.Wait()
	close(resultChan)
	if err != nil {
		return nil, err
	}

	responses := make(map[string]*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, len(awsReqs))
	for result := range resultChan {
		responses[result.EntryId] = result.DataResponse
	}

	return &framer.InterpolatedAssetPropertyValue{
		Responses: responses,
		Query:     query,
	}, nil
}
