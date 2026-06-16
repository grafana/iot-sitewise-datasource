package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/iot-sitewise-datasource/pkg/server"
)

func main() {
	if err := datasource.Manage("grafana-iot-sitewise-datasource", server.NewServerInstance, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
