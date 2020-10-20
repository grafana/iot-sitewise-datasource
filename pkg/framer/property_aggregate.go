package framer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	"github.com/pkg/errors"
)

type AssetPropertyAggregates iotsitewise.GetAssetPropertyAggregatesOutput

func newAggregationField(length int, name string) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableFloat64, length)
	field.Name = name
	return field
}

// getAggregationFields enforces ordering of aggregate fields
// Golang maps return a random order during iteration
func getAggregationFields(length int, aggs *iotsitewise.Aggregates) ([]string, map[string]*data.Field) {

	aggregateTypes := []string{}
	aggregateFields := map[string]*data.Field{}

	if val := aggs.Average; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateAvg)
		aggregateFields[models.AggregateAvg] = newAggregationField(length, "avg")
	}

	if val := aggs.Minimum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateMin)
		aggregateFields[models.AggregateMin] = newAggregationField(length, "min")
	}

	if val := aggs.Maximum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateMax)
		aggregateFields[models.AggregateMax] = newAggregationField(length, "max")
	}

	if val := aggs.Sum; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateSum)
		aggregateFields[models.AggregateSum] = newAggregationField(length, "sum")
	}

	if val := aggs.Count; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateCount)
		aggregateFields[models.AggregateAvg] = newAggregationField(length, "count")
	}

	if val := aggs.StandardDeviation; val != nil {
		aggregateTypes = append(aggregateTypes, models.AggregateStdDev)
		aggregateFields[models.AggregateStdDev] = newAggregationField(length, "std. dev.")
	}

	return aggregateTypes, aggregateFields
}

func addAggregateFieldValues(idx int, fields map[string]*data.Field, aggs *iotsitewise.Aggregates) {

	if val := aggs.Average; val != nil {
		fields[models.AggregateAvg].Set(idx, aggs.Average)
	}

	if val := aggs.Minimum; val != nil {
		fields[models.AggregateMin].Set(idx, aggs.Minimum)
	}

	if val := aggs.Maximum; val != nil {
		fields[models.AggregateMax].Set(idx, aggs.Maximum)
	}

	if val := aggs.Sum; val != nil {
		fields[models.AggregateSum].Set(idx, aggs.Sum)
	}

	if val := aggs.Count; val != nil {
		fields[models.AggregateCount].Set(idx, aggs.Count)
	}

	if val := aggs.StandardDeviation; val != nil {
		fields[models.AggregateStdDev].Set(idx, aggs.StandardDeviation)
	}

}

func (a AssetPropertyAggregates) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(a.AggregatedValues)

	if length < 1 {
		return nil, errors.New("no aggregation values found for query")
	}

	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, length)
	timeField.Name = "time"
	// this will enforce ordering
	aggregateTypes, aggregateFields := getAggregationFields(length, a.AggregatedValues[0].Value)

	for i, v := range a.AggregatedValues {
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

	return data.Frames{frame}, nil
}
