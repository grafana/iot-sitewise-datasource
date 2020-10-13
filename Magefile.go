//+build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Default configures the default target.
func Default() {
	build.BuildAll()
}

// MockGen generates mocks.
// this can be changed to look at the directives when more mocks are needed.
func MockGen() error {

	if err := sh.RunV("docker", "pull", "vektra/mockery"); err != nil {
		return err
	}

	return sh.RunV("go", "generate", "./pkg/...")
}
