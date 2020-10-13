package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func main() {
	backend.SetupPluginEnvironment("timestream-datasource")
	//
	//err := datasource.Serve(timestream.NewDatasource())
	//
	//// Log any error if we could start the plugin.
	//if err != nil {
	//	backend.Logger.Error(err.Error())
	//	os.Exit(1)
	//}
}
