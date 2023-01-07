# gqlgen-appsync

Plugin for gqlgen that allows it to be used with AWS AppSync

## Backlog

- SHOULD setup e2e test with deployed graphql api and our example resolver implementation
- SHOULD test that it works if the model is generated in a different directory
- SHOULD test that it works with errors inside of the batch, instead of returning from lambda
- SHOULD test Relay setup with connections and cursors and node endpoint
- SHOULD test an example that makes use of the "source" input being serialized, what guarantees?

## Document Backlog

- SHOULD document that this plugin only works if it outputs to the same directory as the "exec" output
