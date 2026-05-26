package server

import (
	"net/http"

	generated "mock-pipefy-api/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func Init() *Server {
	return &Server{}
}

type Server struct{}

func (s *Server) Start() error {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &generated.Resolver{},
			},
		),
	)

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	if err := http.ListenAndServe(":8001", nil); err != nil {
		return err
	}

	return nil
}
