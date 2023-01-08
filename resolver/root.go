package resolver

import (
	"context"

	"github.com/crewlinker/gqlgen-appsync/graph"
)

// Root resolver provides all the other resolvers
type Root struct {
	query graph.QueryResolver
}

// New inits the root resolver
func New() graph.ResolverRoot {
	return &Root{query: &query{}}
}

// Query returns the Query resolver
func (r Root) Query() graph.QueryResolver { return r.query }

// query implements the query resolver
type query struct{}

func (query) Version(ctx context.Context) (string, error) { return "v0.1.2", nil }
