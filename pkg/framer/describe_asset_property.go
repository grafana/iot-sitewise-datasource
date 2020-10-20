package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetProperty iotsitewise.DescribeAssetPropertyOutput

func (ap AssetProperty) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	panic("implement me!!!")
}
