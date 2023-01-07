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

## Backlog

- [ ] MUST add ci for e2e testing of generated resolver code
- [ ] SHOULD test that it works with all AWS scalar values
- [ ] SHOULD test that it works with AWS appsync subscription directives
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
