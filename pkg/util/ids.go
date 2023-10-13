package util

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func GetAssetId(query models.BaseQuery) *string {
	if len(query.AssetIds) == 0 {
		return nil
	}
	return aws.String(query.AssetIds[0])
}

func GetPropertyId(query models.BaseQuery) *string {
	if query.PropertyId == "" {
		return nil
	}
	return aws.String(query.PropertyId)
}

func GetPropertyAlias(query models.BaseQuery) *string {
	if query.PropertyAlias == "" {
		return nil
	}
	return aws.String(query.PropertyAlias)
}

func GetEntryId(query models.BaseQuery) *string {
	if query.PropertyAlias != "" && len(query.AssetIds) == 0 {
		return aws.String(strings.ReplaceAll(query.PropertyAlias, "/", "_"))
	}
	return GetAssetId(query)
}
