package framer

import (
	"context"
	"encoding/json"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetModels iotsitewise.ListAssetModelsOutput

func getErrorDescription(details *iotsitewise.ErrorDetails) (*string, error) {
	jb, err := json.Marshal(*details)
	if err != nil {
		return nil, err
	}
	return aws.String(string(jb)), nil
}

func (a AssetModels) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.AssetModelSummaries)

	fName := newFieldWithName(fields.Name, data.FieldTypeString, length)
	fArn := newFieldWithName(fields.Arn, data.FieldTypeString, length)
	fDescription := newFieldWithName(fields.Description, data.FieldTypeString, length)
	fID := newFieldWithName(fields.Id, data.FieldTypeString, length)
	fCreationDate := newFieldWithName(fields.CreationDate, data.FieldTypeTime, length)
	fLastUpdate := newFieldWithName(fields.LastUpdate, data.FieldTypeTime, length)
	fStatusError := newFieldWithName("error", data.FieldTypeNullableString, length)
	fStatusState := newFieldWithName("state", data.FieldTypeString, length)

	for i, asset := range a.AssetModelSummaries {

		fName.Set(i, *asset.Name)
		fArn.Set(i, *asset.Arn)
		fDescription.Set(i, *asset.Description)
		fID.Set(i, *asset.Id)
		fCreationDate.Set(i, *asset.CreationDate)
		fLastUpdate.Set(i, *asset.LastUpdateDate)

		if asset.Status != nil {
			if asset.Status.Error != nil {
				val, err := getErrorDescription(asset.Status.Error)
				if err != nil {
					fStatusError.Set(i, val)
				}
			}
			fStatusState.Set(i, *asset.Status.State)
		}
	}

	frame := data.NewFrame("",
		fName, fDescription, fID, fArn, fStatusError, fStatusState, fCreationDate, fLastUpdate,
	)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
