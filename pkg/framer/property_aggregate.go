package framer

import (
	"context"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyAggregates struct {
	Request  iotsitewise.GetAssetPropertyAggregatesInput
	Response iotsitewise.GetAssetPropertyAggregatesOutput
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
		getFrameName(property),
		fields...,
	)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  util.Dereference(resp.NextToken),
			Resolution: util.Dereference(a.Request.Resolution),
			Aggregates: aggregateTypesToStrings(a.Request.AggregateTypes),
		},
	}

	return data.Frames{frame}, nil
}
