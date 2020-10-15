package fdata

import "github.com/aws/aws-sdk-go/service/iotsitewise"

type AssetProperty iotsitewise.DescribeAssetPropertyOutput

func (ap AssetProperty) Rows() [][]interface{} {
	panic("implement me!!!")
}
