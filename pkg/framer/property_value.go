package framer

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValue iotsitewise.GetAssetPropertyValueOutput

func (p AssetPropertyValue) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	length := 0
	if p.PropertyValue != nil {
		length = 1
	}

	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueField(property, length)
	qualityField := fields.QualityField(length)

	frame := data.NewFrame(getFrameName(property), timeField, valueField, qualityField)

	if p.PropertyValue != nil && getPropertyVariantValue(p.PropertyValue.Value) != nil {
		timeField.Set(0, getTime(p.PropertyValue.Timestamp))
		valueField.Set(0, getPropertyVariantValue(p.PropertyValue.Value))
		qualityField.Set(0, *p.PropertyValue.Quality)
	}

	return data.Frames{frame}, nil
}
