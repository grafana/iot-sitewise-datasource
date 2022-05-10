package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type InterpolatedAssetPropertyValue struct {
	*iotsitewise.GetInterpolatedAssetPropertyValuesOutput
	Query models.AssetPropertyValueQuery
}

func (p InterpolatedAssetPropertyValue) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	length := len(p.InterpolatedAssetPropertyValues)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: make this work with the API instead of ad-hoc dataType inference
	// https://github.com/grafana/iot-sitewise-datasource/issues/98#issuecomment-892947756
	if *property.AssetProperty.DataType == *aws.String("?") {
		property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(p.InterpolatedAssetPropertyValues[0].Value))
	}

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueFieldForQuery(p.Query, property, length)

	frame := data.NewFrame(*property.AssetName, timeField, valueField)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(p.NextToken),
			Resolution: "RAW", //circular dep
		},
	}

	for i, v := range p.InterpolatedAssetPropertyValues {
		value := getPropertyVariantValue(v.Value)
		if value == nil {
			continue
		}
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, value)
	}

	return data.Frames{frame}, nil
}
