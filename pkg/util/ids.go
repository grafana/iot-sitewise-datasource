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
	if len(query.PropertyIds) == 0 {
		return nil
	}
	return aws.String(query.PropertyIds[0])
}

func GetEntryIdFromAssetPropertyEntry(entry models.AssetPropertyEntry) *string {
	if entry.AssetId != "" && entry.PropertyId != "" {
		return GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
	} else {
		return GetEntryIdFromPropertyAlias(entry.PropertyAlias)
	}
}

func GetEntryIdFromPropertyAlias(propertyAlias string) *string {
	return aws.String(EncodeEntryId(propertyAlias))
}

func GetEntryIdFromAssetProperty(assetId string, propertyId string) *string {
	// Encode the assetId and propertyId to create a unique entryId
	return aws.String(EncodeEntryId(assetId + "-" + propertyId))
}

func GetEntryId(query models.BaseQuery) *string {
	// if stream is unassociated, use the hashed property alias as the entry id
	if len(query.PropertyAliases) > 0 && len(query.AssetIds) == 0 && len(query.PropertyIds) == 0 {
		// Use the first property alias as the entry id
		return GetEntryIdFromPropertyAlias(query.PropertyAliases[0])
	}
	assetId := GetAssetId(query)
	propertyId := GetPropertyId(query)
	return GetEntryIdFromAssetProperty(*assetId, *propertyId)
}

// API constraint: EntryId cannot be more than 64 characters long, so we're encoding it
func EncodeEntryId(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
