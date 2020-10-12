package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

// SitewiseResourceProvider is an interface which returns asset/property/model metadata for common query identifiers
type SitewiseResourceProvider interface {
	Asset(ctx context.Context, assetId string) (*iotsitewise.DescribeAssetOutput, error)
	Property(ctx context.Context, assetId string, propertyId string) (*iotsitewise.DescribeAssetPropertyOutput, error)
	AssetModel(ctx context.Context, modelId string) (*iotsitewise.DescribeAssetModelOutput, error)
}
