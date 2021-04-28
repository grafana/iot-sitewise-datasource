package models

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const (
	QueryTypePropertyValuesForTimeRange = "PropertyValuesForTimeRange"
	QueryTypePropertyValueHistory       = "PropertyValueHistory"
	QueryTypePropertyValue              = "PropertyValue"
	QueryTypePropertyAggregate          = "PropertyAggregate"
	QueryTypeListAssetModels            = "ListAssetModels"
	QueryTypeListAssets                 = "ListAssets"
	QueryTypeListAssociatedAssets       = "ListAssociatedAssets"
	QueryTypeDescribeAsset              = "DescribeAsset"
	QueryTypeDescribeAssetModel         = "DescribeAssetModel"
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
	AwsRegion           string `json:"region,omitempty"`
	AssetId             string `json:"assetId,omitempty"`
	PropertyId          string `json:"propertyId,omitempty"`
	NextToken           string `json:"nextToken,omitempty"`
	MaxPageAggregations int    `json:"maxPageAggregations,omitempty"`

	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
	QueryType     string            `json:"-"`
}
