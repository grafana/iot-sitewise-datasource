package framer

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

type PropertyValueQueryFramer struct {
	framer.FrameData
	framer.MetaProvider
	Request models.AssetPropertyValueQuery
}

func (q PropertyValueQueryFramer) Frames(ctx context.Context) (data.Frames, error) {

	md, err := q.MetaProvider.Provide(ctx)

	if err != nil {
		return nil, err
	}

	fields, err := md.Fields()

	if err != nil {
		return nil, err
	}

	frame := data.NewFrame(
		md.FrameName(),
		fields...,
	)

	for _, v := range q.FrameData.Rows() {
		frame.AppendRow(v...)
	}

	return data.Frames{frame}, nil
}
