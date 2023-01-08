# gqlgen-appsync

Plugin for gqlgen that allows it to be used with AWS AppSync

## Usage

Since it is a plugin you'll need to run gqlgen with your own entrypoints. This process is
documented more in-depth over here: https://gqlgen.com/reference/plugins/

For this plugin, your entrypoint will look something like the code below. It MUST output
in the same directory and package as the "exec" base plugin of gqlgen. This is because the
plugin uses private methods that base gqlgen generates.

```Go
import (
	"fmt"
	"path/filepath"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/crewlinker/gqlgen-appsync/genappsync"
)

func main(){
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err = api.Generate(cfg,
        // add the plugin here to use it
		api.AddPlugin(genappsync.New("<output_dir>", "<output_pkg>")),
	); err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	return nil
}
```

After running `gqlgen` the plugin should have produced an new go file that holds the code
that is required to implement the lambda resolver for AWS AppSync. It doesn't require any
request or response mapping, this is also called ["direct lambda resolving"](https://docs.aws.amazon.com/appsync/latest/devguide/direct-lambda-reference.html).

Your direct lambda resolver can use the generated code and will look something like this:

```Go
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
	// Input describes the input of a direct BATCH call from AWS AppSync that
	// we need for the generated resolve call to work.
	Input = []struct {
        // Argumens should be presented as a map
		Arguments map[string]any  `json:"arguments"`

        // NOTE source will be decoded inside the generated code. Used for non-root resolvers
        // to present the "parent" object to the resolver at the time of calling.
		Source    json.RawMessage `json:"source"`
		Info      struct {

            // From the info, we just need the "fieldName" and the "parentTypeName"
			FieldName           string `json:"fieldName"`
			ParentTypeName      string `json:"parentTypeName"`
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

        // NOTE: Here we use the generated code to turn lambda input into calls on
        // our ResolverRoot
		data, err := graph.AppSyncResolve(ctx, h.rr,
			call.Info.ParentTypeName, call.Info.FieldName, call.Arguments, call.Source)
		if err != nil {
			return nil, err
		}

        // NOTE: for batch calls we need to present the data on a "data" field. For
        // non-batch calls this works differently. Check the AWS docs for this.
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
```

## Backlog

- [x] MUST fix bug with uncompilable source for schema without resolver arguments
- [ ] MUST add ci for e2e testing of generated resolver code
- [ ] SHOULD add test case for uncompilable source for schema without resolver arguments (require tests with multiple schemas)
- [ ] SHOULD test all the AWS scalars, and describe as feature
- [ ] SHOULD test with funky type, arg and field names (identifiers), camelCase, snake_case, etc
- [ ] SHOULD come with an input event struct that can be used (because aws-lambda-go/events doesn't provide it)
- [ ] SHOULD test that it works with all AWS scalar values
- [ ] SHOULD test that it works with AWS appsync subscription directives
- [ ] SHOULD test that it works with auth directives
- [ ] SHOULD setup e2e test with deployed graphql api and our example resolver implementation
- [ ] SHOULD test that it works if the model is generated in a different directory
- [ ] SHOULD test that it works with errors inside of the batch, instead of returning from lambda
- [ ] SHOULD test Relay setup with connections and cursors and node endpoint
- [ ] SHOULD test an example that makes use of the "source" input being serialized, what guarantees?
- [ ] SHOULD check if we can build the plugin in such a way that it adjust the model configuration to always
      emit local resolvers for all fields with arguments. BUT: what if an embedded field without arguments
      also needs a resolver? That is a more exceptional exception.

## Document Backlog

- SHOULD document that this plugin only works if it outputs to the same directory as the "exec" output
