package framer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

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
			return parseTimestamp(s)
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

	if value != nil {
		field.Set(index, value)
	}

	return nil
}

func parseTimestamp(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported timestamp format: %q", s)
}
