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

// MockGen generates mocks from directives
func MockGen() error {
	return sh.RunV("go", "generate", "./pkg/...")
}
