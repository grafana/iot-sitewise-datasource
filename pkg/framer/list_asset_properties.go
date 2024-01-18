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

type AssetProperties iotsitewise.ListAssetPropertiesOutput

func (a AssetProperties) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.AssetPropertySummaries)

	fAlias := fields.AliasField(length)
	fId := fields.IdField(length)
	fUnit := fields.UnitField(length)

	for i, assetProperty := range a.AssetPropertySummaries {
		fAlias.Set(i, *assetProperty.Alias)
		fId.Set(i, assetProperty.Id)
		fUnit.Set(i, assetProperty.Unit)
	}

	frame := data.NewFrame("", fAlias, fId)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}