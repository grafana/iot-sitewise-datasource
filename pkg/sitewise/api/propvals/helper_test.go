package propvals

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"testing"
)

type scenario struct {
	query          models.BaseQuery
	expectedResult string
}

var scenarios []scenario = []scenario{
	{
		query: models.BaseQuery{
			TimeRange:     backend.TimeRange{},
			MaxDataPoints: 720,
		},
		expectedResult: ResolutionRaw,
	},
}

func TestResolution(t *testing.T) {

}
