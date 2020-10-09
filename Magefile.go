//+build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Default configures the default target.
var Default = build.BuildAll

var Test = func() error {

	// generate mocks
	if err := sh.RunV("mockery", "--dir=pkg/sitewise/client/", "--name=Client"); err != nil {
		return err
	}

	return build.Test()
}
