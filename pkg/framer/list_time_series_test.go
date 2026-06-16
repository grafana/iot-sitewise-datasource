package framer

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeSeriesFrames_SetsAllFieldTypes(t *testing.T) {
	creationTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	updateTime := time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC)

	ts := TimeSeries(iotsitewise.ListTimeSeriesOutput{
		TimeSeriesSummaries: []iotsitewisetypes.TimeSeriesSummary{
			{
				Alias:                    aws.String("my-alias"),
				AssetId:                  aws.String("asset-123"),
				DataType:                 iotsitewisetypes.PropertyDataTypeDouble,
				DataTypeSpec:             aws.String("DOUBLE"),
				PropertyId:               aws.String("prop-456"),
				TimeSeriesArn:            aws.String("arn:aws:iotsitewise:us-east-1:123:time-series/ts-789"),
				TimeSeriesId:             aws.String("ts-789"),
				TimeSeriesCreationDate:   &creationTime,
				TimeSeriesLastUpdateDate: &updateTime,
			},
		},
		NextToken: aws.String("next-page"),
	})

	frames, err := ts.Frames(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, frames, 1)

	frame := frames[0]
	require.Len(t, frame.Fields, 9, "frame should have all 9 time series summary fields")

	expected := []struct {
		name  string
		value interface{}
	}{
		{fields.Alias, "my-alias"},
		{fields.AssetId, "asset-123"},
		{fields.DataType, "DOUBLE"},
		{fields.DataTypeSpec, "DOUBLE"},
		{fields.PropertyId, "prop-456"},
		{fields.TimeSeriesArn, "arn:aws:iotsitewise:us-east-1:123:time-series/ts-789"},
		{fields.TimeSeriesId, "ts-789"},
		{fields.TimeSeriesCreationDate, creationTime},
		{fields.TimeSeriesLastUpdateDate, updateTime},
	}
	for i, exp := range expected {
		assert.Equal(t, exp.name, frame.Fields[i].Name)
		assert.Equal(t, exp.value, frame.Fields[i].At(0))
	}

	// Meta
	require.NotNil(t, frame.Meta)
	require.NotNil(t, frame.Meta.Custom)
	custom, ok := frame.Meta.Custom.(models.SitewiseCustomMeta)
	require.True(t, ok)
	assert.Equal(t, "next-page", custom.NextToken)
}
