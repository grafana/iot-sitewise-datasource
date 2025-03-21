package framer

import (
	"context"
	"fmt"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"

	"github.com/pkg/errors"
)

type Rows []iotsitewisetypes.Row

type QueryResults iotsitewise.ExecuteQueryOutput

func (a QueryResults) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {
	length := len(a.Rows)
	f := make([]*data.Field, 0)

	for _, col := range a.Columns {
		f = append(f, fields.DatumField(length, col))
	}

	for i, row := range a.Rows {
		for j, datum := range row.Data {
			if datum.ScalarValue == nil {
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
			NextToken: util.Dereference(a.NextToken),
		},
	}

	return data.Frames{frame}, nil
}

func SetValue(col iotsitewisetypes.ColumnInfo, scalarValue string, field *data.Field, index int) error {
	typeConverter := map[iotsitewisetypes.ScalarType]func(string) (interface{}, error){
		iotsitewisetypes.ScalarTypeBoolean: func(s string) (interface{}, error) {
			return strconv.ParseBool(s)
		},
		iotsitewisetypes.ScalarTypeInt: func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 64)
		},
		iotsitewisetypes.ScalarTypeString: func(s string) (interface{}, error) {
			return s, nil
		},
		iotsitewisetypes.ScalarTypeDouble: func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		},
		iotsitewisetypes.ScalarTypeTimestamp: func(s string) (interface{}, error) {
			if t, err := strconv.ParseInt(s, 10, 64); err == nil {
				return time.Unix(0, t*int64(time.Second)), nil
			} else {
				return nil, err
			}
		},
	}

	converter, exists := typeConverter[col.Type.ScalarType]
	if !exists {
		return errors.New(fmt.Sprintf("Unsupported scalar type: %s", col.Type.ScalarType))
	}

	value, err := converter(scalarValue)
	if err != nil {
		return err
	}

	// Override event_timestamp columns to be time values
	if *col.Name == "event_timestamp" {
		intValue := value.(int64)
		// Detect if value is in seconds or nanoseconds
		if intValue < 10000000000 {
			intValue = intValue * 1000000000
		}
		value = time.Unix(0, intValue)
	}

	if value != nil {
		field.Set(index, value)
	}

	return nil
}
