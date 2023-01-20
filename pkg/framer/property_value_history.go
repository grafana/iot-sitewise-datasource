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

type AssetPropertyValueHistory struct {
	*iotsitewise.BatchGetAssetPropertyValueHistoryOutput
	Query models.AssetPropertyValueQuery
}

func (p AssetPropertyValueHistory) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	frames := make(data.Frames, 0, len(p.SuccessEntries))

	for _, h := range p.SuccessEntries {
		frame, err := p.Frame(ctx, resources, h.AssetPropertyValueHistory)
		if err != nil {
			return nil, err
		}
		frames = append(frames, frame)
	}

	return frames, nil
}

func (p AssetPropertyValueHistory) Frame(ctx context.Context, resources resource.ResourceProvider, h []*iotsitewise.AssetPropertyValue) (*data.Frame, error) {

	length := len(h)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: make this work with the API instead of ad-hoc dataType inference
	// https://github.com/grafana/iot-sitewise-datasource/issues/98#issuecomment-892947756
	if *property.AssetProperty.DataType == *aws.String("?") {
		property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(h[0].Value))
	}

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueFieldForQuery(p.Query, property, length)
	qualityField := fields.QualityField(length)

	frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(p.NextToken),
			Resolution: "RAW", //circular dep
		},
	}

	for i, v := range h {
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, getPropertyVariantValue(v.Value))
		qualityField.Set(i, *v.Quality)
	}

	return frame, nil
}
