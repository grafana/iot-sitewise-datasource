package framer

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type Assets iotsitewise.ListAssetsOutput

func getErrorDescription(details *iotsitewise.ErrorDetails) (*string, error) {

	if details == nil {
		return nil, nil
	}

	jb, err := json.Marshal(*details)
	if err != nil {
		return nil, err
	}
	return aws.String(string(jb)), nil
}

func getAssetSummaryHierarchies(asset *iotsitewise.AssetSummary) (string, error) {

	hvalues := []iotsitewise.AssetHierarchy{}

	for _, h := range asset.Hierarchies {
		hvalues = append(hvalues, *h)
	}

	heirarchies, err := json.Marshal(hvalues)
	if err != nil {
		return "", err
	}
	return string(heirarchies), nil
}

func (a Assets) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	length := len(a.AssetSummaries)

	fName := newFieldWithName(fields.Name, data.FieldTypeString, length)
	fId := newFieldWithName(fields.Id, data.FieldTypeString, length)
	fArn := newFieldWithName(fields.Arn, data.FieldTypeString, length)
	fModelId := newFieldWithName(fields.Id, data.FieldTypeString, length)
	fStatusError := newFieldWithName(fields.StatusError, data.FieldTypeNullableString, length)
	fStatusState := newFieldWithName(fields.StatusState, data.FieldTypeString, length)
	fHierarchies := newFieldWithName(fields.Id, data.FieldTypeString, length)
	fCreationDate := newFieldWithName(fields.CreationDate, data.FieldTypeTime, length)
	fLastUpdate := newFieldWithName(fields.LastUpdate, data.FieldTypeTime, length)

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

		heirarchies, err := getAssetSummaryHierarchies(asset)
		if err != nil {
			return nil, err
		}
		fHierarchies.Set(i, heirarchies)
	}

	allFields := fieldsSlice(fName, fId, fModelId, fArn, fCreationDate, fLastUpdate, fStatusState, fStatusError, fHierarchies)

	frame := data.NewFrame("", allFields...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
