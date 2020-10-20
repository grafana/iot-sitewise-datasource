package framer

import (
	"context"
	"encoding/json"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

const assetModelsFrameName = "Asset Models"

type AssetModels iotsitewise.ListAssetModelsOutput

func getAssetModelDescription(asset *iotsitewise.AssetModelSummary) (*string, error) {
	jb, err := json.Marshal(*asset.Status)
	if err != nil {
		return nil, err
	}
	return aws.String(string(jb)), nil
}

func (a AssetModels) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.AssetModelSummaries)

	fName := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)

	fArn := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)
	fDescription := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)
	fId := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)
	fCreationDate := data.NewFieldFromFieldType(data.FieldTypeNullableTime, length)
	fLastUpdate := data.NewFieldFromFieldType(data.FieldTypeNullableTime, length)
	fStatus := data.NewFieldFromFieldType(data.FieldTypeNullableString, length)

	for i, asset := range a.AssetModelSummaries {

		fName.Set(i, asset.Name)
		fArn.Set(i, asset.Arn)
		fDescription.Set(i, asset.Description)
		fId.Set(i, asset.Id)
		fCreationDate.Set(i, asset.CreationDate)
		fLastUpdate.Set(i, asset.LastUpdateDate)

		summary, err := getAssetModelDescription(asset)
		if err != nil {
			return nil, err
		}

		fStatus.Set(i, summary)
	}

	frame := data.NewFrame(assetModelsFrameName,
		fName, fDescription, fId, fArn, fStatus, fCreationDate, fLastUpdate,
	)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
			HasSeries: true,
		},
	}

	return data.Frames{frame}, nil
}
