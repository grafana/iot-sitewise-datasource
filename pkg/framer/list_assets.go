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

type Assets iotsitewise.ListAssetsOutput

func (a Assets) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {

	length := len(a.AssetSummaries)

	fName := fields.NameField(length)
	fId := fields.IdField(length)
	fArn := fields.ArnField(length)
	fModelId := fields.ModelIdField(length)
	fStatusError := fields.StatusErrorField(length)
	fStatusState := fields.StatusStateField(length)
	fHierarchies := fields.HierarchiesField(length)
	fCreationDate := fields.CreationDateField(length)
	fLastUpdate := fields.LastUpdateField(length)

	for i, asset := range a.AssetSummaries {
		fName.Set(i, *asset.Name)
		fId.Set(i, *asset.Id)
		fArn.Set(i, *asset.Arn)
		fModelId.Set(i, *asset.AssetModelId)
		fStatusState.Set(i, *asset.Status.State)
		fCreationDate.Set(i, *asset.CreationDate)
		fLastUpdate.Set(i, *asset.LastUpdateDate)

		statusErr, err := getErrorDescription(asset.Status.Error)
		if err != nil {
			return nil, err
		}
		fStatusError.Set(i, statusErr)

		hierarchies, err := getAssetHierarchies(asset.Hierarchies)
		if err != nil {
			return nil, err
		}
		fHierarchies.Set(i, hierarchies)
	}

	allFields := data.Fields{fName, fId, fModelId, fArn, fCreationDate, fLastUpdate, fStatusState, fStatusError, fHierarchies}

	frame := data.NewFrame("", allFields...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
