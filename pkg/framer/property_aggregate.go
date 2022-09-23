package framer

import (
	"context"
	"fmt"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyAggregates struct {
	Request  iotsitewise.GetAssetPropertyAggregatesInput
	Response iotsitewise.GetAssetPropertyAggregatesOutput
}

// getAggregationFields enforces ordering of aggregate fields
// Golang maps return a random order during iteration
func getAggregationFields(length int, aggs *iotsitewise.Aggregates) ([]string, map[string]*data.Field) {

	aggregateTypes := []string{}
	aggregateFields := map[string]*data.Field{}

	if val := aggs.Average; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateAvg)
		aggregateFields[models.AggregateAvg] = fields.AggregationField(length, "avg")
	}

	if val := aggs.Minimum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateMin)
		aggregateFields[models.AggregateMin] = fields.AggregationField(length, "min")
	}

	if val := aggs.Maximum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateMax)
		aggregateFields[models.AggregateMax] = fields.AggregationField(length, "max")
	}

	if val := aggs.Sum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateSum)
		aggregateFields[models.AggregateSum] = fields.AggregationField(length, "sum")
	}

	if val := aggs.Count; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateCount)
		aggregateFields[models.AggregateCount] = fields.AggregationField(length, "count")
	}

	if val := aggs.StandardDeviation; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateStdDev)
		aggregateFields[models.AggregateStdDev] = fields.AggregationField(length, "stddev")
	}

	return aggregateTypes, aggregateFields
}

func addAggregateFieldValues(idx int, fields map[string]*data.Field, aggs *iotsitewise.Aggregates) {

	if val := aggs.Average; val != nil {
		fields[models.AggregateAvg].Set(idx, *aggs.Average)
	}

	if val := aggs.Minimum; val != nil {
		fields[models.AggregateMin].Set(idx, *aggs.Minimum)
	}

	if val := aggs.Maximum; val != nil {
		fields[models.AggregateMax].Set(idx, *aggs.Maximum)
	}

	if val := aggs.Sum; val != nil {
		fields[models.AggregateSum].Set(idx, *aggs.Sum)
	}

	if val := aggs.Count; val != nil {
		fields[models.AggregateCount].Set(idx, *aggs.Count)
	}

	if val := aggs.StandardDeviation; val != nil {
		fields[models.AggregateStdDev].Set(idx, *aggs.StandardDeviation)
	}

}

func (a AssetPropertyAggregates) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	resp := a.Response

	length := len(resp.AggregatedValues)

	if length < 1 {
		return data.Frames{}, nil
	}

	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := fields.TimeField(length)
	// this will enforce ordering
	aggregateTypes, aggregateFields := getAggregationFields(length, resp.AggregatedValues[0].Value)

	for i, v := range resp.AggregatedValues {
		timeField.Set(i, *v.Timestamp)
		addAggregateFieldValues(i, aggregateFields, v.Value)
	}

	fields := []*data.Field{timeField}

	for _, aggType := range aggregateTypes {
		fields = append(fields, aggregateFields[aggType])
	}

	frame := data.NewFrame(
		fmt.Sprintf("%s %s", *property.AssetName, *property.AssetProperty.Name),
		fields...,
	)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(resp.NextToken),
			Resolution: aws.StringValue(a.Request.Resolution),
			Aggregates: aws.StringValueSlice(a.Request.AggregateTypes),
		},
	}

	return data.Frames{frame}, nil
}
