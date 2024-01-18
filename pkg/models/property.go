package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ListAssetPropertiesQuery struct {
	BaseQuery
}

// AssetPropertyValueQuery encapsulates params for all 3 'Get' data APIs in Sitewise.
// Each API handler will simply ignore the unneeded props.
// This is done simply due to lack of solid generics support in golang.
type AssetPropertyValueQuery struct {
	BaseQuery
	AggregateTypes  []string `json:"aggregates,omitempty"` // Not used for the history query
	Quality         string   `json:"quality,omitempty"`
	Resolution      string   `json:"resolution,omitempty"`
	LastObservation bool     `json:"lastObservation,omitempty"`
	TimeOrdering    string   `json:"timeOrdering,omitempty"`
}

func GetAssetPropertyValueQuery(dq *backend.DataQuery) (*AssetPropertyValueQuery, error) {

	query := &AssetPropertyValueQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// AssetId <--> AssetIds backward compatibility
	query.MigrateAssetId()

	//if propertyAlias is set make sure to set the assetId and propertyId to nil
	if query.PropertyAlias != "" {
		query.AssetId = ""
		query.PropertyId = ""
	}

	if query.TimeOrdering == "" {
		query.TimeOrdering = "ASCENDING"
	}

	// default to 1 if unset
	if query.MaxPageAggregations < 1 {
		query.MaxPageAggregations = 1
	}

	// add on the DataQuery params
	query.TimeRange = dq.TimeRange
	query.Interval = dq.Interval
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}
