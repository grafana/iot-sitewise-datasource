package framer

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Framer is an interface that allows any type to be treated as a data frame
type Framer interface {
	Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error)
}

// FrameResponse creates a backend.DataResponse that contains the Framer's data.Frames
func FrameResponse(ctx context.Context, f Framer) backend.DataResponse {

	frames, err := f.Frames(ctx, nil)

	return backend.DataResponse{
		Frames: frames,
		Error:  err,
	}
}

// FrameResponseWithError creates a backend.DataResponse with the error's contents (if not nil), and the Framer's data.Frames
// This function is particularly useful if you have a function that returns `(Framer, error)`, which is a very common pattern
//func FrameResponseWithError(f Framer, ctx context.Context, err error) backend.DataResponse {
//	if err != nil {
//		return backend.DataResponse{
//			Error: err,
//		}
//	}
//
//	return FrameResponse(f, ctx)
//}
