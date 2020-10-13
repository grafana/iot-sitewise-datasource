//+build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Default configures the default target.
func Default() {
	_ = MockGen()
	build.BuildAll()
}

// MockGen generates mocks.
// this can be changed to look at the directives when more mocks are needed.
func MockGen() error {
	//return sh.RunV("mockery", "--dir=pkg/sitewise/client/", "--name=Client")
	return sh.RunV("go", "generate", "./pkg/...")

}
