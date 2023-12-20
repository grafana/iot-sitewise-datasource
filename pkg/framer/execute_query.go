package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type Rows []*iotsitewise.Row

type QueryResults iotsitewise.ExecuteQueryOutput

func (a QueryResults) Frames(ctx context.Context, res resource.ResourceProvider) (data.Frames, error) {
	length := len(a.Rows)
	f := make([]*data.Field, 0)

	for _, col := range a.Columns {
		f = append(f, fields.DatumField(*col, length))
	}

	for i, row := range a.Rows {
		for j, datum := range row.Data {
			f[j].Set(i, datum)
		}
	}

	frame := data.NewFrame("", f...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
