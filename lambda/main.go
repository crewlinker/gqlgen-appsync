package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/crewlinker/gqlgen-appsync/graph"
	"github.com/crewlinker/gqlgen-appsync/resolver"
)

type (
	// Input describes the input of a direct batch call from AWS AppSync that
	// we need for the generated resolve call to work
	Input = []struct {
		Arguments map[string]any  `json:"arguments"`
		Source    json.RawMessage `json:"source"`
		Info      struct {
			FieldName           string `json:"fieldName"`
			ParentTypeName      string `json:"parentTypeName"`
			SelectionSetGraphQL string `json:"selectionSetGraphQL"`
		} `json:"info"`
	}

	// Output for a direct batch call from AWS AppSync
	Output = []map[string]any
)

// Handler handles lambda inputs
type Handler struct {
	rr graph.ResolverRoot
}

// Handle direct lambda resolving from aws AppSync
func (h Handler) Handle(ctx context.Context, in Input) (out Output, err error) {
	log.Printf("Input: %+v", in)

	for _, call := range in {
		data, err := graph.AppSyncResolve(ctx, h.rr,
			call.Info.ParentTypeName, call.Info.FieldName, call.Arguments, call.Source)
		if err != nil {
			return nil, err
		}

		out = append(out, map[string]any{
			"data": data,
		})
	}

	log.Printf("Output: %+v", out)
	return
}

// lambda entry point
func main() {
	lambda.Start((Handler{rr: resolver.New()}).Handle)
}
