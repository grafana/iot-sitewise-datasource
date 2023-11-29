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

type InterpolatedAssetPropertyValue struct {
	Responses map[string]*iotsitewise.GetInterpolatedAssetPropertyValuesOutput
	Query     models.AssetPropertyValueQuery
}

func (p InterpolatedAssetPropertyValue) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	properties, err := resources.Properties(ctx)

	if err != nil {
		return nil, err
	}

	frames := data.Frames{}

	for entryId, res := range p.Responses {
		property := properties[entryId]
		if property == nil {
			property = properties[*util.GetEntryId(p.Query.BaseQuery)]
		}
		frame, err := p.Frame(ctx, property, res.InterpolatedAssetPropertyValues)
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame)
	}

	return frames, nil
}

func (p InterpolatedAssetPropertyValue) Frame(ctx context.Context, property *iotsitewise.DescribeAssetPropertyOutput, v []*iotsitewise.InterpolatedAssetPropertyValue) (*data.Frame, error) {
	// TODO: make this work with the API instead of ad-hoc dataType inference
	// https://github.com/grafana/iot-sitewise-datasource/issues/98#issuecomment-892947756
	if *property.AssetProperty.DataType == *aws.String("?") {
		property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(v[0].Value))
	}

	timeField := fields.TimeField(0)
	valueField := fields.PropertyValueFieldForQuery(p.Query, property, 0)
	name := *property.AssetName
	if name == "" {
		name = *property.AssetProperty.Name
	}
	frame := data.NewFrame(name, timeField, valueField)

	entryId := ""
	if property.AssetId != nil {
		entryId = *property.AssetId
	} else {
		entryId = *property.AssetProperty.Name
	}
	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken:  aws.StringValue(p.Responses[entryId].NextToken),
			EntryId:    entryId,
			Resolution: p.Query.Resolution,
		},
	}

	for _, v := range v {
		value := getPropertyVariantValue(v.Value)
		if value == nil {
			continue
		}
		timeField.Append(getTime(v.Timestamp))
		valueField.Append(value)
	}

	return frame, nil
}
