package framer

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValueHistory iotsitewise.GetAssetPropertyValueHistoryOutput

func (p AssetPropertyValueHistory) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(p.AssetPropertyValueHistory)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueField(property, length)
	qualityField := fields.QualityField(length)

	frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(p.NextToken),
			Resolution: propvals.ResolutionRaw,
		},
	}

	for i, v := range p.AssetPropertyValueHistory {
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, getPropertyVariantValue(v.Value))
		qualityField.Set(i, *v.Quality)
	}

	return data.Frames{frame}, nil
}
