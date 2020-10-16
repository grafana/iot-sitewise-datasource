package fdata

import "github.com/aws/aws-sdk-go/service/iotsitewise"

type AssetPropertyValueHistory iotsitewise.GetAssetPropertyValueHistoryOutput

func (p AssetPropertyValueHistory) Rows() [][]interface{} {
	var rows [][]interface{}

	for _, v := range p.AssetPropertyValueHistory {
		rows = append(rows, []interface{}{
			getTimeValue(v.Timestamp),
			getPropertyVariantValue(v.Value),
		})
	}
	return rows
}
