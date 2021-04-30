package propvals

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

type scenario struct {
	name     string
	query    models.BaseQuery
	expected string
}

var scenarios = []scenario{
	{
		// dps = 300, pages = 2
		name: "selects '1s' resolution",
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
			MaxDataPoints: 720,
		},
		expected: ResolutionSecond,
	},
	{
		// dps = 120, pages = 1
		name: "selects '1m' resolution",
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{From: testdata.TwoHours, To: testdata.Now},
			MaxDataPoints: 720,
		},
		expected: ResolutionMinute,
	},
	{
		// dps = 24, pages = 1
		name: "selects '1h' resolution",
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{From: testdata.OneDay, To: testdata.Now},
			MaxDataPoints: 720,
		},
		expected: ResolutionHour,
	},
	{
		// dps = 31, pages = 1
		name: "selects '1d' resolution",
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{From: testdata.OneMonth, To: testdata.Now},
			MaxDataPoints: 720,
		},
		expected: ResolutionDay,
	},
	{
		// dps = 300, pages = 2
		name: "elevates '1s' to '1m' when MaxDataPoints is less than total data points",
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{From: testdata.FiveMinutes, To: testdata.Now},
			MaxDataPoints: 299,
		},
		expected: ResolutionMinute,
	},
}

func TestResolution(t *testing.T) {

	for _, scene := range scenarios {
		t.Run(scene.name, func(t *testing.T) {
			actual := Resolution(scene.query)
			assert.Equal(t, scene.expected, actual)
		})
	}
}
