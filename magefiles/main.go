package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// init performs some sanity checks before running anything
func init() {
	mustBeInRoot()
}

// Test perform the whole project's unit tests
func Test() error {
	return sh.Run(
		"go", "run", "-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "-repeat=5", "--fail-on-pending", "--race",
		"--trace", "--junit-report=test-report.xml", "./...")
}

// E2E performs an end-to-end test
func E2E(env, instance string) error {
	if err := Graph.Gen(Graph{}); err != nil {
		return err
	}
	if err := Infra.Build(Infra{}); err != nil {
		return err
	}
	if err := Infra.Deploy(Infra{}, env, instance); err != nil {
		return err
	}

	os.Chdir("..") // back to root after infra command
	if err := Test(); err != nil {
		return err
	}
	return nil
}

// mustBeInRoot checks that the command is run in the project root
func mustBeInRoot() {
	if _, err := os.Stat("go.mod"); err != nil {
		panic("must be in root, couldn't stat go.mod file: " + err.Error())
	}
}

// profileFromEnv determines the AWS credentials profile from the env argument
func profileFromEnv(env string) string {
	if env != "prod" && env != "stag" && env != "dev" {
		panic("invalid env: '" + env + "', supports: 'prod', 'stag' or 'dev'")
	}
	return fmt.Sprintf("cl-%s", env)
}
