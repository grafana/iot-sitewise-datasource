package util

import (
	"crypto/sha256"
	"encoding/hex"

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
	// if stream is unassociated, use the hashed property alias as the entry id
	if query.PropertyAlias != "" && len(query.AssetIds) == 0 && query.PropertyId == "" {
		// API constraint: EntryId cannot be more than 64 characters long, so we're encoding it 
		return aws.String(EncodeEntryId(query.PropertyAlias))
	}
	return GetAssetId(query)
}

func EncodeEntryId(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
