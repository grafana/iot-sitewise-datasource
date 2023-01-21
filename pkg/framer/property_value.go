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
	frames := data.Frames{}
	length := len(p.SuccessEntries)

	properties, err := resources.Properties(ctx)
	if err != nil {
		return nil, err
	}

	for i, e := range p.SuccessEntries {
		property := properties[*e.EntryId]
		timeField := fields.TimeField(length)
		valueField := fields.PropertyValueField(property, length)
		qualityField := fields.QualityField(length)

		frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

		if e.AssetPropertyValue != nil {
			timeField.Set(i, getTime(e.AssetPropertyValue.Timestamp))
			valueField.Set(i, getPropertyVariantValue(e.AssetPropertyValue.Value))
			qualityField.Set(i, *e.AssetPropertyValue.Quality)
		}
		frames = append(frames, frame)
	}

	for _, e := range p.ErrorEntries {
		property := properties[*e.EntryId]
		frame := data.NewFrame(*property.AssetName)
		if e.ErrorMessage != nil {
			frame.Meta = &data.FrameMeta{
				Notices: []data.Notice{{Severity: data.NoticeSeverityError, Text: *e.ErrorMessage}},
			}
		}
		frames = append(frames, frame)
	}

	return frames, nil
}
