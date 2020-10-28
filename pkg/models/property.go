package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// AssetPropertyValueQuery encapsulates params for all 3 'Get' data APIs in Sitewise.
// Each API handler will simply ignore the unneeeded props.
// NOTES: We have decided to not support propertyAlias targets, as there is no good way to go from propertyAlias -> assetId/propertyId.
// This is done simply due to lack of solid generics support in golang.
type AssetPropertyValueQuery struct {
	BaseQuery
	Qualities      []string `json:"qualities,omitempty"`
	AggregateTypes []string `json:"aggregates,omitempty"`
	Resolution     string   `json:"resolution,omitempty"`
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
