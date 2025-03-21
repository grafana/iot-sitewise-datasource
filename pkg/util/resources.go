package util

import (
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
)

func IsAssetProperty(property *iotsitewise.DescribeAssetPropertyOutput) bool {
	return property.AssetProperty != nil
}

func IsComponentProperty(property *iotsitewise.DescribeAssetPropertyOutput) bool {
	return property.CompositeModel != nil && property.CompositeModel.AssetProperty != nil
}

func GetPropertyName(property *iotsitewise.DescribeAssetPropertyOutput) string {
	if IsAssetProperty(property) {
		return *property.AssetProperty.Name
	} else if IsComponentProperty(property) {
		return *property.CompositeModel.AssetProperty.Name
	}

	return ""
}

func GetPropertyUnit(property *iotsitewise.DescribeAssetPropertyOutput) string {
	if IsAssetProperty(property) && property.AssetProperty.Unit != nil {
		return *property.AssetProperty.Unit
	} else if IsComponentProperty(property) && property.CompositeModel.AssetProperty.Unit != nil {
		return *property.CompositeModel.AssetProperty.Unit
	}

	return ""
}

func GetPropertyDataType(property *iotsitewise.DescribeAssetPropertyOutput) iotsitewisetypes.PropertyDataType {
	if IsAssetProperty(property) {
		return property.AssetProperty.DataType
	} else if IsComponentProperty(property) {
		return property.CompositeModel.AssetProperty.DataType
	}

	return ""
}
