package models

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const (
	QueryTypePropertyValueHistory = "PropertyValueHistory"
	QueryTypePropertyValue        = "PropertyValue"
	QueryTypePropertyAggregate    = "PropertyAggregate"
	QueryTypePropertyInterpolated = "PropertyInterpolated"
	QueryTypeListAssetModels      = "ListAssetModels"
	QueryTypeListAssets           = "ListAssets"
	QueryTypeListAssociatedAssets = "ListAssociatedAssets"
	QueryTypeDescribeAsset        = "DescribeAsset"
	QueryTypeDescribeAssetModel   = "DescribeAssetModel"
	QueryTypeListAssetProperties  = "ListAssetProperties"
	QueryTypeListTimeSeries       = "ListTimeSeries"
	QueryTypeExecuteQuery         = "ExecuteQuery"
)

const (
	AggregateMin    = "MINIMUM"
	AggregateMax    = "MAXIMUM"
	AggregateCount  = "COUNT"
	AggregateAvg    = "AVERAGE"
	AggregateStdDev = "STANDARD_DEVIATION"
	AggregateSum    = "SUM"
)

type BaseQuery struct {
	// General
	AwsRegion string `json:"region,omitempty"`

	QueryType string `json:"-"`

	// Sitewise specific
	// Deprecated: use AssetIds
	AssetId  string   `json:"assetId,omitempty"`
	AssetIds []string `json:"assetIds,omitempty"`
	// Deprecated: use PropertyIds
	PropertyId  string   `json:"propertyId,omitempty"`
	PropertyIds []string `json:"propertyIds,omitempty"`
	// Deprecated: use PropertyAliases
	PropertyAlias        string               `json:"propertyAlias,omitempty"`
	PropertyAliases      []string             `json:"propertyAliases,omitempty"`
	AssetPropertyEntries []AssetPropertyEntry `json:"assetPropertyEntries,omitempty"`
	NextToken            string               `json:"nextToken,omitempty"`
	NextTokens           map[string]string    `json:"nextTokens,omitempty"`
	MaxPageAggregations  int                  `json:"maxPageAggregations,omitempty"`
	ResponseFormat       string               `json:"responseFormat,omitempty"`

	// Also provided by sqlutil.Query. Migrate to that
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int32             `json:"-"`
}

// MigrateAssetProperty handles AssetId, PropertyId, PropertyAlias --> AssetIds, PropertyIds, PropertyAliases backward compatibility.
// This is needed for compatibility for queries saved before the Batch API changes were introduced in 2.1.0
func (query *BaseQuery) MigrateAssetProperty() {
	if query.AssetId != "" {
		query.AssetIds = []string{query.AssetId}
	}

	if query.PropertyId != "" {
		query.PropertyIds = []string{query.PropertyId}
	}

	if query.PropertyAlias != "" {
		query.PropertyAliases = []string{query.PropertyAlias}
	}
}
