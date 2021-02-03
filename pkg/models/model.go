package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ListAssetModelsQuery struct {
	BaseQuery
}

type DescribeAssetModelQuery struct {
	BaseQuery
	AssetModelId string `json:"assetModelId"`
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

func GetDescribeAssetModelQuery(dq *backend.DataQuery) (*DescribeAssetModelQuery, error) {
	query := &DescribeAssetModelQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.QueryType = dq.QueryType

	return query, nil
}
