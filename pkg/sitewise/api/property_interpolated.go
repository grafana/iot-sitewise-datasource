package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

var (
	LOCF_INTERPOLATION   string = "LOCF_INTERPOLATION"
	LINEAR_INTERPOLATION string = "LINEAR_INTERPOLATION"
)

func interpolatedQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetInterpolatedAssetPropertyValuesInput {
	//if propertyAlias is set make sure to set the assetId and propertyId to nil
	if query.PropertyAlias != "" {
		query.PropertyId = ""
		query.AssetId = ""
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)
	startTimeInSeconds := from.Unix()
	endTimeInSeconds := to.Unix()

	quality := query.Quality
	if quality == "" {
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

	log.DefaultLogger.Error("Querying interpolated asset property values", "intervalInSeconds", intervalInSeconds, "startTimeInSeconds", startTimeInSeconds, "endTimeInSeconds", endTimeInSeconds, "quality", quality, "interpolationType", interpolationType)

	return &iotsitewise.GetInterpolatedAssetPropertyValuesInput{
		StartTimeInSeconds: &startTimeInSeconds,
		EndTimeInSeconds:   &endTimeInSeconds,
		IntervalInSeconds:  aws.Int64(intervalInSeconds),
		MaxResults:         aws.Int64(10),
		NextToken:          getNextToken(query.BaseQuery),
		AssetId:            getAssetId(query.BaseQuery),
		PropertyId:         getPropertyId(query.BaseQuery),
		PropertyAlias:      getPropertyAlias(query.BaseQuery),
		Quality:            &quality,
		Type:               &interpolationType,
	}
}

func GetInterpolatedAssetPropertyValues(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.InterpolatedAssetPropertyValue, error) {
	maxDps := int(query.MaxDataPoints)
	awsReq := interpolatedQueryToInput(query)

	resp, err := client.GetInterpolatedAssetPropertyValuesPageAggregation(ctx, awsReq, query.MaxPageAggregations, maxDps)
	if err != nil {
		return nil, err
	}

	return &framer.InterpolatedAssetPropertyValue{
		GetInterpolatedAssetPropertyValuesOutput: resp,
		Query:                                    query,
	}, nil
}
