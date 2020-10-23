package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type DescribeAssetQuery struct {
	BaseQuery
}

type DescribeAssetPropertyQuery struct {
	BaseQuery
}

type ListAssetsQuery struct {
	BaseQuery
	ModelId string `json:"modelId,omitempty"`
	Filter  string `json:"filter,omitempty"`
}

func GetDescribeAssetQuery(dq *backend.DataQuery) (*DescribeAssetQuery, error) {
	query := &DescribeAssetQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.QueryType = dq.QueryType
	return query, nil
}

func GetListAssetsQuery(dq *backend.DataQuery) (*ListAssetsQuery, error) {
	query := &ListAssetsQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType
	return query, nil
}
