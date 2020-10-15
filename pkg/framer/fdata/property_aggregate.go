package fdata

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type AssetPropertyAggregates iotsitewise.GetAssetPropertyAggregatesOutput

func (a AssetPropertyAggregates) Rows() [][]interface{} {
	var rows [][]interface{}

	for _, v := range a.AggregatedValues {
		row := []interface{}{v.Timestamp.Unix()}
		row = append(row, aggregateValues(v.Value)...)
		rows = append(rows, row)
	}

	return rows
}

func aggregateValues(value *iotsitewise.Aggregates) []interface{} {
	var vals []interface{}

	for _, k := range models.AggregateOrder {
		if agg, ok := models.AggregateFields[k]; ok {
			if val := agg.ValueGetter(value); val != nil {
				vals = append(vals, val)
			}
		}
	}

	return vals
}
