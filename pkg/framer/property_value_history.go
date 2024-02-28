package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

type AssetPropertyValueHistory struct {
	*iotsitewise.GetAssetPropertyValueHistoryOutput
	Query models.AssetPropertyValueQuery
}

func (p AssetPropertyValueHistory) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(p.AssetPropertyValueHistory)
	property, err := resources.Property(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: make this work with the API instead of ad-hoc dataType inference
	// https://github.com/grafana/iot-sitewise-datasource/issues/98#issuecomment-892947756
	if util.IsAssetProperty(property) && *property.AssetProperty.DataType == *aws.String("?") {
		property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(p.AssetPropertyValueHistory[0].Value))
	}

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueFieldForQuery(p.Query, property, length)
	qualityField := fields.QualityField(length)
	frame := data.NewFrame(getFrameName(property), timeField, valueField, qualityField)
	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(p.NextToken),
			Resolution: "RAW", //circular dep
		},
	}

	for i, v := range p.AssetPropertyValueHistory {
		timeField.Set(i, getTime(v.Timestamp))
		valueField.Set(i, getPropertyVariantValue(v.Value))
		qualityField.Set(i, *v.Quality)
	}

	return data.Frames{frame}, nil
}