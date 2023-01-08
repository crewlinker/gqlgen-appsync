package tests_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tests")
}

var _ = Describe("graphql", func() {
	var cl *graphql.Client
	BeforeEach(func() {
		out := ReadOutputs()
		cl = graphql.NewClient(out.ClGASyncAppMain.GraphGraphQlUrl, http.DefaultClient)
		cl = cl.WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-api-key", out.ClGASyncAppMain.GraphGraphQlApiKey)
		})
	})

	DescribeTable("queries", func(ctx context.Context, query any, exp string) {
		Expect(cl.Query(ctx, query, nil)).To(Succeed())
		Expect(json.Marshal(query)).To(MatchJSON(exp))
	},
		Entry("string", &struct{ Version string }{}, `{"Version":"v0.1.2"}`),
	)
})

func ReadOutputs() (out struct {
	ClGASyncAppMain struct {
		GraphGraphQlUrl    string `json:"GraphGraphQlUrlEEC0AE8C"`
		GraphGraphQlApiKey string `json:"GraphGraphQlApiKey50B87373"`
	}
}) {
	data, err := os.ReadFile(filepath.Join("..", "infra", "cdk.outputs.json"))
	Expect(err).ToNot(HaveOccurred())
	Expect(json.Unmarshal(data, &out)).To(Succeed())
	return
}
