package framer

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValue iotsitewise.BatchGetAssetPropertyValueOutput

func (p AssetPropertyValue) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(p.SuccessEntries)
	var frame *data.Frame

	for i, e := range p.SuccessEntries {
		property, err := resources.Property(ctx)
		if err != nil {
			return nil, err
		}

		timeField := fields.TimeField(length)
		valueField := fields.PropertyValueField(property, length)
		qualityField := fields.QualityField(length)

		frame = data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

		if e.AssetPropertyValue != nil {
			timeField.Set(i, getTime(e.AssetPropertyValue.Timestamp))
			valueField.Set(i, getPropertyVariantValue(e.AssetPropertyValue.Value))
			qualityField.Set(i, *e.AssetPropertyValue.Quality)
		}
	}

	return data.Frames{frame}, nil
}
