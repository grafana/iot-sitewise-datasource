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
	Request  iotsitewise.BatchGetAssetPropertyAggregatesInput
	Response iotsitewise.BatchGetAssetPropertyAggregatesOutput
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
	frames := data.Frames{}

	properties, err := resources.Properties(ctx)
	if err != nil {
		return nil, err
	}

	for i, e := range resp.SuccessEntries {
		property := properties[*e.EntryId]
		frame, err := a.Frame(ctx, property, e.AggregatedValues)
		if err != nil {
			return nil, err
		}
		frame.Meta = &data.FrameMeta{
			Custom: models.SitewiseCustomMeta{
				NextToken:  aws.StringValue(resp.NextToken),
				Resolution: aws.StringValue(a.Request.Entries[i].Resolution),
				Aggregates: aws.StringValueSlice(a.Request.Entries[i].AggregateTypes),
			},
		}

		for _, e := range resp.ErrorEntries {
			property := properties[*e.EntryId]
			frame := data.NewFrame(*property.AssetName)
			if e.ErrorMessage != nil {
				frame.Meta = &data.FrameMeta{
					Notices: []data.Notice{{Severity: data.NoticeSeverityError, Text: *e.ErrorMessage}},
				}
			}
			frames = append(frames, frame)
		}

		frames = append(frames, frame)
	}

	return frames, nil
}

func (a AssetPropertyAggregates) Frame(ctx context.Context, property *iotsitewise.DescribeAssetPropertyOutput, v []*iotsitewise.AggregatedValue) (*data.Frame, error) {

	length := len(v)
	if length < 1 {
		return &data.Frame{}, nil
	}

	timeField := fields.TimeField(length)
	// this will enforce ordering
	aggregateTypes, aggregateFields := getAggregationFields(length, v[0].Value)

	for i, v := range v {
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

	return frame, nil

}
