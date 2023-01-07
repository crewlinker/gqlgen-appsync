package appstack

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/crewlinker/gqlgen-appsync/graph"
)

// WithResources builds the resources for the instanced app stack
func WithResources(s constructs.Construct) {
	WithAppSync(s, awslogs.RetentionDays_ONE_DAY)
}

func WithAppSync(s constructs.Construct, logRetention awslogs.RetentionDays) {
	s = constructs.NewConstruct(s, jsii.String("Graph"))

	api := awsappsync.NewCfnGraphQLApi(s, jsii.String("GraphApi"), &awsappsync.CfnGraphQLApiProps{
		AuthenticationType: jsii.String("API_KEY"),
		Name:               jsii.String(*awscdk.Stack_Of(s).StackName() + "Graph"),
	})

	def, err := os.ReadFile(filepath.Join("..", "example.graphql"))
	if err != nil {
		panic("failed to load graphql definition: " + err.Error())
	}

	lambda := awslambda.NewFunction(s, jsii.String("GraphHandler"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(jsii.String(filepath.Join("..", "lambda", "pkg.zip")), nil),
		Handler:      jsii.String("bootstrap"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		LogRetention: logRetention,
		Tracing:      awslambda.Tracing_ACTIVE,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(50)),
	})

	role := awsiam.NewRole(s, jsii.String("GraphServiceRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("appsync.amazonaws.com"), nil),
	})

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: &[]*string{lambda.FunctionArn()},
		Actions:   jsii.Strings("lambda:InvokeFunction"),
	}))

	schema := awsappsync.NewCfnGraphQLSchema(s, jsii.String("GraphSchema"), &awsappsync.CfnGraphQLSchemaProps{
		ApiId:      api.AttrApiId(),
		Definition: jsii.String(string(def)),
	})

	awsappsync.NewCfnApiKey(s, jsii.String("GraphApiKey"), &awsappsync.CfnApiKeyProps{
		ApiId:       api.AttrApiId(),
		Description: jsii.String("Main API Key"),
		ApiKeyId:    jsii.String("MainApiKey"),
	})

	ds := awsappsync.NewCfnDataSource(s, jsii.String("GraphLambdaSource"), &awsappsync.CfnDataSourceProps{
		ApiId:          api.AttrApiId(),
		Name:           jsii.String("MainLambda"),
		Type:           jsii.String("AWS_LAMBDA"),
		ServiceRoleArn: role.RoleArn(),
		LambdaConfig: awsappsync.CfnDataSource_LambdaConfigProperty{
			LambdaFunctionArn: lambda.FunctionArn(),
		},
	})

	for _, typfield := range graph.ResolverFields {
		typ, field, _ := strings.Cut(typfield, ".")
		awsappsync.NewCfnResolver(s, jsii.String(typ+field+"GraphResolver"), &awsappsync.CfnResolverProps{
			ApiId:          api.AttrApiId(),
			TypeName:       jsii.String(typ),
			FieldName:      jsii.String(field),
			DataSourceName: ds.AttrName(),
			MaxBatchSize:   jsii.Number(10), // enable batching for direct lambda
		}).AddDependency(schema)
	}
}
