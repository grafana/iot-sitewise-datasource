package resource

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

// SitewiseResourceProvider is domain specific an interface which returns asset/property/model descriptions
type ResourceProvider interface {
	Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error)
	Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error)
	AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error)
}
