package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetModelDescription iotsitewise.DescribeAssetModelOutput

type describeAssetModelFields struct {
	Name            *data.Field
	Arn             *data.Field
	Id              *data.Field
	Properties      *data.Field
	CompositeModels *data.Field
	Hierarchies     *data.Field
	CreationDate    *data.Field
	Description     *data.Field
	LastUpdate      *data.Field
	StatusError     *data.Field
	StatusState     *data.Field
}

func (f *describeAssetModelFields) fields() data.Fields {
	return data.Fields{
		f.Name,
		f.Id,
		f.Arn,
		f.Description,
		f.StatusState,
		f.StatusError,
		f.CreationDate,
		f.LastUpdate,
		f.Hierarchies,
		f.Properties,
		f.CompositeModels,
	}
}

func newDescribeAssetModelFields() *describeAssetModelFields {
	return &describeAssetModelFields{
		Name:            fields.NameField(1),
		Arn:             fields.ArnField(1),
		Id:              fields.IdField(1),
		Description:     fields.DescriptionField(1),
		Properties:      fields.PropertiesField(1),
		CompositeModels: fields.CompositeModelsField(1),
		Hierarchies:     fields.HierarchiesField(1),
		CreationDate:    fields.CreationDateField(1),
		LastUpdate:      fields.LastUpdateField(1),
		StatusError:     fields.StatusErrorField(1),
		StatusState:     fields.StatusStateField(1),
	}
}

func (a AssetModelDescription) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	assetModelFields := newDescribeAssetModelFields()

	assetModelFields.Name.Set(0, *a.AssetModelName)
	assetModelFields.Arn.Set(0, *a.AssetModelArn)
	assetModelFields.Description.Set(0, *a.AssetModelDescription)
	assetModelFields.Id.Set(0, *a.AssetModelId)
	assetModelFields.CreationDate.Set(0, *a.AssetModelCreationDate)
	assetModelFields.LastUpdate.Set(0, *a.AssetModelLastUpdateDate)
	assetModelFields.StatusState.Set(0, *a.AssetModelStatus.State)

	if a.AssetModelStatus.Error != nil {
		statusErr, err := getErrorDescription(a.AssetModelStatus.Error)
		if err != nil {
			return nil, err
		}
		assetModelFields.StatusError.Set(0, statusErr)
	}

	hierarchies, err := serialize(a.AssetModelHierarchies)
	if err != nil {
		return nil, err
	}
	assetModelFields.Hierarchies.Set(0, hierarchies)

	properties, err := serialize(a.AssetModelProperties)
	if err != nil {
		return nil, err
	}
	assetModelFields.Properties.Set(0, properties)

	compositeModels, err := serialize(a.AssetModelCompositeModels)
	if err != nil {
		return nil, err
	}
	assetModelFields.CompositeModels.Set(0, compositeModels)

	frame := data.NewFrame("", assetModelFields.fields()...)

	return data.Frames{frame}, nil
}
