package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ListAssetModelsQuery struct {
	BaseQuery
	NextToken string `json:"nextToken,omitempty"`
}

func GetListAssetModelsQuery(dq *backend.DataQuery) (*ListAssetModelsQuery, error) {

	query := &ListAssetModelsQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}
