package framer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetDescription iotsitewise.DescribeAssetOutput

type describeAssetFields struct {
	Name         *data.Field
	Arn          *data.Field
	Id           *data.Field
	ModelId      *data.Field
	Properties   *data.Field
	Hierarchies  *data.Field
	CreationDate *data.Field
	LastUpdate   *data.Field
	StatusError  *data.Field
	StatusState  *data.Field
}

func (f *describeAssetFields) fields() data.Fields {
	return data.Fields{
		f.Name,
		f.Id,
		f.Arn,
		f.ModelId,
		f.StatusState,
		f.StatusError,
		f.CreationDate,
		f.LastUpdate,
		f.Hierarchies,
		f.Properties,
	}
}

func newDescribeAssetFields() *describeAssetFields {
	return &describeAssetFields{
		Name:         fields.NameField(1),
		Arn:          fields.ArnField(1),
		Id:           fields.IdField(1),
		ModelId:      fields.ModelIdField(1),
		Properties:   fields.PropertiesField(1),
		Hierarchies:  fields.HierarchiesField(1),
		CreationDate: fields.CreationDateField(1),
		LastUpdate:   fields.LastUpdateField(1),
		StatusError:  fields.StatusErrorField(1),
		StatusState:  fields.StatusStateField(1),
	}
}

func (a AssetDescription) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	assetFields := newDescribeAssetFields()

	assetFields.Name.Set(0, *a.AssetName)
	assetFields.Id.Set(0, *a.AssetId)
	assetFields.ModelId.Set(0, *a.AssetModelId)
	assetFields.CreationDate.Set(0, *a.AssetCreationDate)
	assetFields.LastUpdate.Set(0, *a.AssetLastUpdateDate)
	assetFields.StatusState.Set(0, string(a.AssetStatus.State))

	if a.AssetStatus.Error != nil {
		statusErr, err := getErrorDescription(a.AssetStatus.Error)
		if err != nil {
			return nil, err
		}
		assetFields.StatusError.Set(0, statusErr)
	}

	hierarchies, err := serialize(a.AssetHierarchies)
	if err != nil {
		return nil, err
	}
	assetFields.Hierarchies.Set(0, hierarchies)

	properties, err := serialize(a.AssetProperties)
	if err != nil {
		return nil, err
	}
	assetFields.Properties.Set(0, properties)

	frame := data.NewFrame("", assetFields.fields()...)

	return data.Frames{frame}, nil
}
