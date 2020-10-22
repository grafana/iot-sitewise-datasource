package sitewise

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func aggregateQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyAggregatesInput {

	var (
		propertyId     *string
		assetId        *string
		nextToken      *string
		aggregateTypes = aws.StringSlice(query.AggregateTypes)
		qualities      []*string
		resolution     = aws.String(query.Resolution)
	)

	assetId = getAssetId(query.BaseQuery)
	propertyId = getPropertyId(query.BaseQuery)
	nextToken = getNextToken(query.BaseQuery)

	if len(query.Qualities) > 0 {
		qualities = aws.StringSlice(query.Qualities)
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
		Resolution:     resolution,
		StartDate:      from,
	}
}

func GetAssetPropertyAggregates(ctx context.Context, client client.Client, query models.AssetPropertyValueQuery) (*framer.AssetPropertyAggregates, error) {

	awsReq := aggregateQueryToInput(query)

	// NOTE: there is a paginated API if we want to push pagination requests down to the server
	// See: https://docs.aws.amazon.com/sdk-for-go/api/service/iotsitewise/#IoTSiteWise.GetAssetPropertyAggregatesPagesWithContext
	resp, err := client.GetAssetPropertyAggregatesWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &framer.AssetPropertyAggregates{
		AggregatedValues: resp.AggregatedValues,
		NextToken:        resp.NextToken,
	}, nil
}
