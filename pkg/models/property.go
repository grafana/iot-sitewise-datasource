package models

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// AssetPropertyValueQuery encapsulates params for all 3 'Get' data APIs in Sitewise.
// Each API handler will simply ignore the unneeeded props.
// NOTES: We have decided to not support propertyAlias targets, as there is no good way to go from propertyAlias -> assetId/propertyId.
// This is done simply due to lack of solid generics support in golang.
type AssetPropertyValueQuery struct {
	BaseQuery
	AssetId        string   `json:"assetId"`
	PropertyId     string   `json:"propertyId"`
	NextToken      string   `json:"nextToken,omitempty"`
	Qualities      []string `json:"qualities,omitempty"`
	AggregateTypes []string `json:"aggregateTypes,omitempty"`
	Resolution     string   `json:"resolution,omitempty"`
	// Not from JSON
	// TODO: move to embedded struct?
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
	QueryType     string            `json:"-"`
}

func GetAssetPropertyValueQuery(dq *backend.DataQuery) (*AssetPropertyValueQuery, error) {

	query := &AssetPropertyValueQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.TimeRange = dq.TimeRange
	query.Interval = dq.Interval
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}

// AggregateFieldHelper is a struct used by both the meta provider and the row data.
// FieldName is used to populate the data frame field name.
// ValueGetter is a helper function for fetching the aggregation value from the Sitewise response.
type AggregateFieldHelper struct {
	FieldName   string
	ValueGetter func(value *iotsitewise.Aggregates) *float64
}

// AggregateOrder is the expected field order for aggregation queries
var AggregateOrder = []string{
	AggregateAvg,
	AggregateMin,
	AggregateMax,
	AggregateSum,
	AggregateCount,
	AggregateStdDev,
}

// AggregateFields assists with an ordering contract between the meta provider and row data.
var AggregateFields = map[string]AggregateFieldHelper{
	AggregateAvg: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.Average
		},
		FieldName: "avg",
	},
	AggregateMin: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.Minimum
		},
		FieldName: "min",
	},
	AggregateMax: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.Maximum
		},
		FieldName: "max",
	},
	AggregateSum: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.Sum
		},
		FieldName: "sum",
	},
	AggregateCount: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.Count
		},
		FieldName: "count",
	},
	AggregateStdDev: {
		ValueGetter: func(value *iotsitewise.Aggregates) *float64 {
			return value.StandardDeviation
		},
		FieldName: "std. dev.",
	},
}
