package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/iot-sitewise-datasource/pkg/plugin"
)

func main() {
	// if err := datasource.Manage("sitewise-datasource", server.NewServerInstance, datasource.ManageOpts{}); err != nil {

	if err := datasource.Manage("sitewise-datasource", plugin.NewSitewiseDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
