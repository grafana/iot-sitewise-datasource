package sitewise

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

type AssetPropertyAggregates iotsitewise.GetAssetPropertyAggregatesOutput

func (a AssetPropertyAggregates) Rows() [][]interface{} {
	var rows [][]interface{}

	for _, v := range a.AggregatedValues {
		row := []interface{}{v.Timestamp.Unix()}
		row = append(row, aggregateValues(v.Value)...)
		rows = append(rows, row)
	}

	return rows
}

func aggregateValues(value *iotsitewise.Aggregates) []interface{} {
	var vals []interface{}

	for _, agg := range models.AggregateFields {
		if val := agg.ValueGetter(value); val != nil {
			vals = append(vals, val)
		}
	}

	return vals
}

func aggregateQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyAggregatesInput {

	var (
		propertyId     *string
		assetId        *string
		nextToken      *string
		aggregateTypes = aws.StringSlice(query.AggregateTypes)
		qualities      []*string
		resolution     = aws.String(query.Resolution)
	)

	if query.AssetId != "" {
		assetId = aws.String(query.AssetId)
	}

	if query.PropertyId != "" {
		propertyId = aws.String(query.PropertyId)
	}

	if query.NextToken != "" {
		nextToken = aws.String(query.NextToken)
	}

	if len(query.Qualities) > 0 {
		qualities = aws.StringSlice(query.Qualities)
	}

	if query.NextToken != "" {
		nextToken = aws.String(query.NextToken)
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

func GetAssetPropertyAggregates(ctx context.Context, client client.Client, query models.AssetPropertyValueQuery) (*AssetPropertyAggregates, error) {

	awsReq := aggregateQueryToInput(query)

	// NOTE: there is a paginated API if we want to push pagination requests down to the server
	// See: https://docs.aws.amazon.com/sdk-for-go/api/service/iotsitewise/#IoTSiteWise.GetAssetPropertyAggregatesPagesWithContext
	resp, err := client.GetAssetPropertyAggregatesWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &AssetPropertyAggregates{
		AggregatedValues: resp.AggregatedValues,
		NextToken:        resp.NextToken,
	}, nil
}
