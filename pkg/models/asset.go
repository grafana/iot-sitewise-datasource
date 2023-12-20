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

type ListAssociatedAssetsQuery struct {
	BaseQuery
	HierarchyId     string `json:"hierarchyId,omitempty"`
	LoadAllChildren bool   `json:"loadAllChildren,omitempty"`
}

type ExecuteQuery struct {
	BaseQuery
	QueryStatement string `json:"queryStatement"`
}

func GetDescribeAssetQuery(dq *backend.DataQuery) (*DescribeAssetQuery, error) {
	query := &DescribeAssetQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// AssetId <--> AssetIds backward compatibility
	query.MigrateAssetId()

	// add on the DataQuery params
	query.QueryType = dq.QueryType
	return query, nil
}

func GetExecuteQuery(dq *backend.DataQuery) (*ExecuteQuery, error) {
	query := &ExecuteQuery{}
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

	// AssetId <--> AssetIds backward compatibility
	query.MigrateAssetId()

	// add on the DataQuery params
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType
	return query, nil
}

func GetListAssociatedAssetsQuery(dq *backend.DataQuery) (*ListAssociatedAssetsQuery, error) {
	query := &ListAssociatedAssetsQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// AssetId <--> AssetIds backward compatibility
	query.MigrateAssetId()

	// add on the DataQuery params
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType
	return query, nil
}
