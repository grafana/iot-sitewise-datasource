package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValue iotsitewise.GetAssetPropertyValueOutput

func (p AssetPropertyValue) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := 1

	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, length)
	timeField.Name = "time"

	valueField := data.NewFieldFromFieldType(fieldTypeForPropertyValue(property), length)
	valueField.Name = *property.AssetProperty.Name

	qualityField := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)
	qualityField.Name = "Quality"

	frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

	timeField.Set(0, getTime(p.PropertyValue.Timestamp))
	valueField.Set(0, getPropertyVariantValue(p.PropertyValue.Value))
	qualityField.Set(0, p.PropertyValue.Quality)

	return data.Frames{frame}, nil
}
