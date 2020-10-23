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

type AssetModels iotsitewise.ListAssetModelsOutput

func (a AssetModels) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.AssetModelSummaries)

	fName := fields.NameField(length)
	fArn := fields.ArnField(length)
	fDescription := fields.DescriptionField(length)
	fId := fields.IdField(length)
	fCreationDate := fields.CreationDateField(length)
	fLastUpdate := fields.LastUpdateField(length)
	fStatusError := fields.StatusErrorField(length)
	fStatusState := fields.StatusStateField(length)

	for i, asset := range a.AssetModelSummaries {
		fName.Set(i, *asset.Name)
		fArn.Set(i, *asset.Arn)
		fDescription.Set(i, *asset.Description)
		fId.Set(i, *asset.Id)
		fCreationDate.Set(i, *asset.CreationDate)
		fLastUpdate.Set(i, *asset.LastUpdateDate)

		if asset.Status.Error != nil {
			val, err := getErrorDescription(asset.Status.Error)
			if err != nil {
				fStatusError.Set(i, val)
			}
		}
		fStatusState.Set(i, *asset.Status.State)
	}

	frame := data.NewFrame("",
		fName, fDescription, fId, fArn, fStatusError, fStatusState, fCreationDate, fLastUpdate,
	)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
