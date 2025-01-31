package framer

import (
	"context"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssociatedAssets iotsitewise.ListAssociatedAssetsOutput

func (a AssociatedAssets) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {

	length := len(a.AssetSummaries)

	assetFields := newAssetSummaryFields(length)

	for i, asset := range a.AssetSummaries {
		assetFields.Name.Set(i, *asset.Name)
		assetFields.Id.Set(i, *asset.Id)
		assetFields.Arn.Set(i, *asset.Arn)
		assetFields.ModelId.Set(i, *asset.AssetModelId)
		assetFields.StatusState.Set(i, string(asset.Status.State))
		assetFields.CreationDate.Set(i, *asset.CreationDate)
		assetFields.LastUpdate.Set(i, *asset.LastUpdateDate)

		statusErr, err := getErrorDescription(asset.Status.Error)
		if err != nil {
			return nil, err
		}
		assetFields.StatusError.Set(i, statusErr)

		hierarchies, err := serialize(asset.Hierarchies)
		if err != nil {
			return nil, err
		}
		assetFields.Hierarchies.Set(i, hierarchies)
	}

	frame := data.NewFrame("", assetFields.fields()...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: util.Dereference(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
