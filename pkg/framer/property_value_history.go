package framer

import (
	"context"

	resource2 "github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type AssetPropertyValueHistory iotsitewise.GetAssetPropertyValueHistoryOutput

func (p AssetPropertyValueHistory) Frames(ctx context.Context, resources resource2.ResourceProvider) (data.Frames, error) {

	length := len(p.AssetPropertyValueHistory)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := data.NewFieldFromFieldType(data.FieldTypeTime, length)
	timeField.Name = "time"

	valueField := data.NewFieldFromFieldType(fieldTypeForPropertyValue(property), length)
	valueField.Name = *property.AssetProperty.Name

	frame := data.NewFrame(*property.AssetName, timeField, valueField)

	for i, v := range p.AssetPropertyValueHistory {
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, getPropertyVariantValue(v.Value))
	}

	return data.Frames{frame}, nil
}
