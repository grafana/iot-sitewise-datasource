package fdata

import "github.com/aws/aws-sdk-go/service/iotsitewise"

type AssetDescription iotsitewise.DescribeAssetOutput

func (a AssetDescription) Rows() [][]interface{} {
	panic("implement me")
}
