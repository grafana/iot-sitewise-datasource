package framer

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValueBatch iotsitewise.BatchGetAssetPropertyValueOutput

func (p AssetPropertyValueBatch) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	frames := data.Frames{}

	properties, err := resources.Properties(ctx)
	if err != nil {
		return nil, err
	}

	for _, e := range p.SuccessEntries {
		property := properties[*e.EntryId]
		if *property.AssetProperty.DataType == *aws.String("?") && e.AssetPropertyValue != nil {
			property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(e.AssetPropertyValue.Value))
		}
		timeField := fields.TimeField(0)
		valueField := fields.PropertyValueField(property, 0)
		qualityField := fields.QualityField(0)

		frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

		if e.AssetPropertyValue != nil {
			timeField.Append(getTime(e.AssetPropertyValue.Timestamp))
			valueField.Append(getPropertyVariantValue(e.AssetPropertyValue.Value))
			qualityField.Append(*e.AssetPropertyValue.Quality)
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