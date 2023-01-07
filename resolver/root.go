package resolver

import (
	"context"
	"fmt"

	"github.com/crewlinker/gqlgen-appsync/graph"
)

// Root resolver provides all the other resolvers
type Root struct {
	query    graph.QueryResolver
	mutation graph.MutationResolver
	sink     graph.SinkResolver
}

// New inits the root resolver
func New() *Root {
	return &Root{query: &Query{}, mutation: &Mutation{}, sink: &Sink{}}
}

func (r Root) Query() graph.QueryResolver       { return r.query }
func (r Root) Mutation() graph.MutationResolver { return r.mutation }
func (r Root) Sink() graph.SinkResolver         { return r.sink }

type Sink struct{}

func (Sink) Other(ctx context.Context, obj *graph.Sink, code graph.Bar) (*graph.Sink, error) {
	name := "other name"
	return &graph.Sink{Name: &name}, nil
}

type Mutation struct{}

func (Mutation) AddSink(ctx context.Context, name string) (*graph.Sink, error) {
	return &graph.Sink{Name: &name}, nil
}

type Query struct{}

func (Query) Version(ctx context.Context) (string, error) { return "AAA", nil }
func (Query) Profile(ctx context.Context) (*graph.Profile, error) {
	email := "foo@foo.com"
	email2 := "other_node@foo.com"
	return &graph.Profile{ID: 124, Email: &email, OtherNode: &graph.Profile{ID: 666, Email: &email2}}, nil
}
func (Query) AnotherProfile(ctx context.Context, name string) (*graph.Profile, error) {
	email := fmt.Sprintf("%v_foo@foo.com", name)
	email2 := fmt.Sprintf("%v_nested@other.com", name)
	return &graph.Profile{Email: &email, OtherNode: &graph.Profile{ID: 666, Email: &email2}}, nil
}
func (Query) KitchenSink(ctx context.Context, foo string, bar graph.Bar) (*graph.Sink, error) {
	name := fmt.Sprintf("%s.%v", foo, bar.Nr)

	return &graph.Sink{Name: &name}, nil
}
func (Query) Node(ctx context.Context, id int) (graph.Node, error) {
	return &graph.Profile{ID: id}, nil
}

func (Query) SinksOrProfiles(ctx context.Context) ([]graph.ProfileOrSink, error) {
	name := "sink1"
	return []graph.ProfileOrSink{
		&graph.Profile{ID: 888}, &graph.Sink{Name: &name},
	}, nil
}
