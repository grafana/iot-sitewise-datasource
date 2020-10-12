package framer

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type MetaProvider interface {
	Provide(cxt context.Context) (Metadata, error)
}

// Metadata is an interface that manages retrieving metadata about the query from the associated models and request
// Ideally this should contain non-domain specific methods for creating a frame
// TODO: have a 'FrameMeta' for getting standard meta from the models. For response specific meta, we may need a second method for FrameData
type Metadata interface {
	// FrameName returns the name used to initialize the data frame.
	FrameName() string
	// Fields returns all fields expected to be present in the data frame.
	Fields() ([]*data.Field, error)
	// FrameMeta
	// FrameMeta() (*data.FrameMeta, error)
}
