package framer

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type Rows []*iotsitewise.Row

type QueryResults iotsitewise.ExecuteQueryOutput

func (a QueryResults) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.Rows)
	f := make([]*data.Field, 0)

	for _, col := range a.Columns {
		f = append(f, fields.DatumField(*col, length))
	}

	for i, row := range a.Rows {
		for j, datum := range row.Data {
			if datum.ScalarValue == nil {
				backend.Logger.Debug("nil datum")
				continue
			}

			err := SetValue(a.Columns[j], *datum.ScalarValue, f[j], i)
			if err != nil {
				backend.Logger.Debug("Error setting value", "error", err)
			}
		}
	}

	frame := data.NewFrame("", f...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			// Not adding the NextToken, since it leads to streaming
			// NextToken: aws.StringValue(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}

func SetValue(col *iotsitewise.ColumnInfo, scalarValue string, field *data.Field, index int) error {
	typeConverter := map[string]func(string) (interface{}, error){
		"BOOLEAN": func(s string) (interface{}, error) {
			return strconv.ParseBool(s)
		},
		"INTEGER": func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 64)
		},
		"STRING": func(s string) (interface{}, error) {
			if col.Name != nil && *col.Name == "event_timestamp" {
				if t, err := strconv.ParseInt(s, 10, 64); err == nil {
					return time.Unix(0, t*int64(time.Nanosecond)), nil
				}
			}
			return s, nil
		},
		"DOUBLE": func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		},
	}

	converter, exists := typeConverter[*col.Type.ScalarType]
	if !exists {
		return nil // or return an error if you want to handle unsupported types
	}

	value, err := converter(scalarValue)
	if err != nil {
		return err
	}

	field.Set(index, value)
	return nil
}
