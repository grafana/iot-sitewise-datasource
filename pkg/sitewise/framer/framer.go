package framer

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

// Framer is an interface that allows any type to be treated as a data frame
type Framer interface {
	Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error)
}
