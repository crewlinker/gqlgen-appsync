package main

import (
	"fmt"
	"path/filepath"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/crewlinker/gqlgen-appsync/genappsync"
	"github.com/magefile/mage/mg"
)

// Graph namespace holds automation for graphql
type Graph mg.Namespace

// Gen generate code from our graphql schema
func (Graph) Gen() error {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err = api.Generate(cfg,
		// add our appsync plugin to generate resolver code for AWS AppSync
		api.AddPlugin(genappsync.New(filepath.Join("graph", "appsync_gen.go"), "graph")),
	); err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	return nil
}
