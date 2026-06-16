package resource

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

type queryResourceProvider struct {
	resources *cachingResourceProvider
	baseQuery models.BaseQuery
}

func NewQueryResourceProvider(cachingProvider *cachingResourceProvider, query models.BaseQuery) *queryResourceProvider {
	return &queryResourceProvider{
		resources: cachingProvider,
		baseQuery: query,
	}
}

func (rp *queryResourceProvider) Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error) {
	assetId := ""

	// use the first asset id if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	return rp.resources.Asset(ctx, assetId)
}

func (rp *queryResourceProvider) Assets(ctx context.Context) (map[string]*iotsitewise.DescribeAssetOutput, error) {
	assets := map[string]*iotsitewise.DescribeAssetOutput{}
	for _, id := range rp.baseQuery.AssetIds {
		asset, err := rp.resources.Asset(ctx, id)
		if err != nil {
			return nil, err
		}
		assets[id] = asset
	}
	return assets, nil
}

func (rp *queryResourceProvider) Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	assetId := ""
	propertyId := ""
	propertyAlias := ""

	// ok to only use the first item in the list since this function is called for non-batch queries

	// use the first assetId if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	// use the first propertyId if there are multiple
	if len(rp.baseQuery.PropertyIds) > 0 {
		propertyId = rp.baseQuery.PropertyIds[0]
	}

	// use the first propertyAlias if there are multiple
	if len(rp.baseQuery.PropertyAliases) > 0 {
		propertyAlias = rp.baseQuery.PropertyAliases[0]
	}

	return rp.resources.Property(ctx, assetId, propertyId, propertyAlias)
}

func (rp *queryResourceProvider) Properties(ctx context.Context) (map[string]*iotsitewise.DescribeAssetPropertyOutput, error) {
	properties := map[string]*iotsitewise.DescribeAssetPropertyOutput{}
	// if the query for a PropertyAlias doesn't have an assetId or propertyId, it means it's a disassociated stream
	// in that case, we call Property() with empty values, which will set AssetProperty.Name to the alias
	// and will set the EntryId to the hashed alias (to access values in results)
	for _, entry := range rp.baseQuery.AssetPropertyEntries {
		prop, err := rp.resources.Property(ctx, entry.AssetId, entry.PropertyId, entry.PropertyAlias)
		if err != nil {
			return nil, err
		}
		entryId := util.GetEntryIdFromAssetPropertyEntry(entry)
		properties[*entryId] = prop
	}

	return properties, nil
}

func (rp *queryResourceProvider) AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error) {
	assetId := ""

	// use the first asset id if there are multiple
	if len(rp.baseQuery.AssetIds) > 0 {
		assetId = rp.baseQuery.AssetIds[0]
	}

	asset, err := rp.resources.Asset(ctx, assetId)

	if err != nil {
		return nil, err
	}

	return rp.resources.AssetModel(ctx, *asset.AssetModelId)
}
