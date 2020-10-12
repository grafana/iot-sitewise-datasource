package models

type DescribeAssetQuery struct {
	AssetId string `json:"assetId"`
}

type DescribeAssetPropertyQuery struct {
	AssetId    string `json:"assetId"`
	PropertyId string `json:"propertyId"`
}
