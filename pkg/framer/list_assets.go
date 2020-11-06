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

type assetSummaryFields struct {
	Name         *data.Field
	Id           *data.Field
	Arn          *data.Field
	ModelId      *data.Field
	StatusError  *data.Field
	StatusState  *data.Field
	Hierarchies  *data.Field
	CreationDate *data.Field
	LastUpdate   *data.Field
}

func (f *assetSummaryFields) fields() data.Fields {
	return data.Fields{
		f.Name,
		f.Id,
		f.ModelId,
		f.Arn,
		f.CreationDate,
		f.LastUpdate,
		f.StatusState,
		f.StatusError,
		f.Hierarchies,
	}
}

func newAssetSummaryFields(length int) *assetSummaryFields {
	return &assetSummaryFields{
		Name:         fields.NameField(length),
		Id:           fields.IdField(length),
		Arn:          fields.ArnField(length),
		ModelId:      fields.ModelIdField(length),
		StatusError:  fields.StatusErrorField(length),
		StatusState:  fields.StatusStateField(length),
		Hierarchies:  fields.HierarchiesField(length),
		CreationDate: fields.CreationDateField(length),
		LastUpdate:   fields.LastUpdateField(length),
	}
}

func (a Assets) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {

	length := len(a.AssetSummaries)

	assetFields := newAssetSummaryFields(length)

	for i, asset := range a.AssetSummaries {
		assetFields.Name.Set(i, *asset.Name)
		assetFields.Id.Set(i, *asset.Id)
		assetFields.Arn.Set(i, *asset.Arn)
		assetFields.ModelId.Set(i, *asset.AssetModelId)
		assetFields.StatusState.Set(i, *asset.Status.State)
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
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
