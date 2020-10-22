package sitewise

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func getNextToken(query models.BaseQuery) *string {
	if query.NextToken == "" {
		return nil
	}
	return aws.String(query.NextToken)
}

func getAssetId(query models.BaseQuery) *string {
	if query.AssetId == "" {
		return nil
	}
	return aws.String(query.AssetId)
}

func getPropertyId(query models.BaseQuery) *string {
	if query.PropertyId == "" {
		return nil
	}
	return aws.String(query.PropertyId)
}
