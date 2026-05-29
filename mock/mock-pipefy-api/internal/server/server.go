package server

import (
	"log/slog"
	"net/http"
	"time"

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

	mux := http.NewServeMux()
	mux.Handle("/query", requestLogMiddleware(srv))
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	slog.Info("mock pipefy api listening", "addr", ":8001")
	if err := http.ListenAndServe(":8001", mux); err != nil {
		return err
	}

	return nil
}

func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
