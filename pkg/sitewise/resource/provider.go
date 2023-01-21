package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

// ResourceProvider is SiteWise domain specific an interface which returns asset/property/model descriptions
type ResourceProvider interface {
	Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error)
	Assets(ctx context.Context) (map[string]*iotsitewise.DescribeAssetOutput, error)
	Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error)
	Properties(ctx context.Context) (map[string]*iotsitewise.DescribeAssetPropertyOutput, error)
	AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error)
}
