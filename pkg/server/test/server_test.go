package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/server"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
)

var (
	timeRange = backend.TimeRange{
		From: time.Now().Add(time.Hour * -3),
		To:   time.Now(),
	}
)

type testServerScenarioFn func(t *testing.T) *testScenario

type testScenario struct {
	name           string
	queries        []backend.DataQuery
	mockSw         *mocks.Client
	goldenFileName string
	handlerFn      func(srvr *server.Server) backend.QueryDataHandlerFunc
	validationFn   func(t *testing.T, dr *backend.QueryDataResponse)
}

func (ts *testScenario) run(t *testing.T) {
	runTestScenario(t, ts)
}

// Golang's cwd is the executable file location.
// hack to find the test data directory
func testDataRelativePath(filename string) string {
	return "../../testdata/" + filename
}

func mockedDatasource(swmock *mocks.Client) server.Datasource {
	return &sitewise.Datasource{
		GetClient: func(_ backend.PluginContext, _ models.BaseQuery) (client client.Client, err error) {
			client = swmock
			return
		},
	}
}

func runTestScenario(t *testing.T, scenario *testScenario) {

	t.Run(scenario.name, func(t *testing.T) {

		ctx := context.Background()

		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{},
			Queries:       scenario.queries,
		}

		srvr := &server.Server{
			Datasource: mockedDatasource(scenario.mockSw),
		}

		qdr, err := scenario.handlerFn(srvr)(ctx, req)

		// this should always be nil, as the error is wrapped in the QueryDataResponse
		if err != nil {
			t.Fatal(err)
		}

		if scenario.validationFn != nil {
			scenario.validationFn(t, qdr)
		}

		// write out the golden for all data responses
		for i, dr := range qdr.Responses {
			fname := fmt.Sprintf(testDataRelativePath("%s-%s.golden.txt"), scenario.goldenFileName, i)

			// temporary fix for golden files https://github.com/grafana/grafana-plugin-sdk-go/issues/213
			for _, fr := range dr.Frames {
				if fr.Meta != nil {
					fr.Meta.Custom = nil
				}
			}

			if err := experimental.CheckGoldenDataResponse(fname, &dr, true); err != nil {
				if !strings.Contains(err.Error(), "no such file or directory") {
					t.Fatal(err)
				}
			}
		}
	})
}
