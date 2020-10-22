package fields

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func NewFieldWithName(name string, fieldType data.FieldType, length int) *data.Field {
	field := data.NewFieldFromFieldType(fieldType, length)
	field.Name = name
	return field
}

func NameField(length int) *data.Field {
	return NewFieldWithName(Name, data.FieldTypeString, length)
}

func IdField(length int) *data.Field {
	return NewFieldWithName(Id, data.FieldTypeString, length)
}

func ArnField(length int) *data.Field {
	return NewFieldWithName(Arn, data.FieldTypeString, length)
}

func ModelIdField(length int) *data.Field {
	return NewFieldWithName(ModelId, data.FieldTypeString, length)
}

func DescriptionField(length int) *data.Field {
	return NewFieldWithName(Description, data.FieldTypeString, length)
}

func StatusErrorField(length int) *data.Field {
	return NewFieldWithName(StatusError, data.FieldTypeNullableString, length)
}

func StatusStateField(length int) *data.Field {
	return NewFieldWithName(StatusState, data.FieldTypeString, length)
}

func HierarchiesField(length int) *data.Field {
	return NewFieldWithName(Hierarchies, data.FieldTypeString, length)
}

func CreationDateField(length int) *data.Field {
	return NewFieldWithName(CreationDate, data.FieldTypeTime, length)
}

func LastUpdateField(length int) *data.Field {
	return NewFieldWithName(LastUpdate, data.FieldTypeTime, length)
}

func TimeField(length int) *data.Field {
	return NewFieldWithName(Time, data.FieldTypeTime, length)
}

func QualityField(length int) *data.Field {
	return NewFieldWithName(Quality, data.FieldTypeString, length)
}

func PropertyValueField(property *iotsitewise.DescribeAssetPropertyOutput, length int) *data.Field {
	valueField := NewFieldWithName(*property.AssetProperty.Name, FieldTypeForPropertyValue(property), length)
	valueField.Config = &data.FieldConfig{
		Unit: ToGrafanaUnit(property.AssetProperty.Unit),
	}
	return valueField
}

func AggregationField(length int, name string) *data.Field {
	return NewFieldWithName(name, data.FieldTypeFloat64, length)
}
