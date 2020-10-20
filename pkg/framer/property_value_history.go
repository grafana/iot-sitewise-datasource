package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValueHistory iotsitewise.GetAssetPropertyValueHistoryOutput

func (p AssetPropertyValueHistory) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(p.AssetPropertyValueHistory)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, length)
	timeField.Name = "time"

	valueField := newPropertyValueField(property, length)

	qualityField := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)
	qualityField.Name = "Quality"

	frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

	for i, v := range p.AssetPropertyValueHistory {
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, getPropertyVariantValue(v.Value))
		qualityField.Set(i, v.Quality)
	}

	return data.Frames{frame}, nil
}
