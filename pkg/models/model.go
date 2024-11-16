package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

type ListAssetModelsQuery struct {
	BaseQuery
}

type DescribeAssetModelQuery struct {
	BaseQuery
	AssetModelId string `json:"assetModelId"`
}

type ExecuteQuery struct {
	BaseQuery
	sqlutil.Query
}

func GetListAssetModelsQuery(dq *backend.DataQuery) (*ListAssetModelsQuery, error) {

	query := &ListAssetModelsQuery{}
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

func GetDescribeAssetModelQuery(dq *backend.DataQuery) (*DescribeAssetModelQuery, error) {
	query := &DescribeAssetModelQuery{}
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
	backend.Logger.Debug("Running GetExecuteQuery", "JSON", dq.JSON)
	query := &ExecuteQuery{}
	if err := json.Unmarshal(dq.JSON, &query); err != nil {
		return nil, err
	}

	query.QueryType = dq.QueryType
	return query, nil
}

func GetQuery(eq *ExecuteQuery) (*sqlutil.Query, error) {
	query := &sqlutil.Query{}

	query.RawSQL = eq.RawSQL

	query.Interval = eq.Query.Interval
	query.TimeRange = eq.Query.TimeRange
	query.MaxDataPoints = eq.Query.MaxDataPoints

	return query, nil
}
