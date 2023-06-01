package api

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/AnthonyThomasson/graph-ql/graph/generated"
	"github.com/AnthonyThomasson/graph-ql/graph/resolver"
	"gorm.io/gorm"
)

func getGqlHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{DB: db}}))
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}

func getPlaygroundHandler() http.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/query")
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
