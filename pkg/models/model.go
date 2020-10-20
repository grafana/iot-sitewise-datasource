package models

type ListAssetModelsQuery struct {
	BaseQuery
	NextToken string `json:"nextToken,omitempty"`
}
