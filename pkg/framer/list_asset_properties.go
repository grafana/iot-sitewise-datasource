package framer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetProperties iotsitewise.ListAssetPropertiesOutput

type assetPropertySummaryFields struct {
	Id   *data.Field
	Name *data.Field
}

func (f *assetPropertySummaryFields) fields() data.Fields {
	return data.Fields{
		f.Id,
		f.Name,
	}
}

func newAssetPropertySummaryFields(length int) *assetPropertySummaryFields {
	return &assetPropertySummaryFields{
		Id:   fields.IdField(length),
		Name: fields.NameField(length),
	}
}

func (a AssetProperties) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.AssetPropertySummaries)

	assetPropertyFields := newAssetPropertySummaryFields(length)

	for i, assetProperty := range a.AssetPropertySummaries {
		assetPropertyFields.Id.Set(i, *assetProperty.Id)
		assetPropertyFields.Name.Set(i, *assetProperty.Path[len(assetProperty.Path)-1].Name)
	}

	frame := data.NewFrame("", assetPropertyFields.fields()...)

	nextToken := ""
	if a.NextToken != nil {
		nextToken = *a.NextToken
	}
	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: nextToken,
		},
	}

	return data.Frames{frame}, nil
}
