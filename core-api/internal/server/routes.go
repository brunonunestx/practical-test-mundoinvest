package server

import (
	"log/slog"
	"net/http"
	"time"

	"core-api/internal/modules/cards"
	"core-api/internal/modules/clients"
	pkg "core-api/packages"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	clientsHandler := clients.NewClientHandler(s.pool, s.pipefy)
	cardsHandler := cards.NewHandler(s.pool, s.pipefy)

	r := mux.NewRouter()

	r.Use(s.corsMiddleware)
	r.Use(s.requestLogMiddleware)

	r.HandleFunc("/clients", clientsHandler.CreateClient).Methods("POST")
	r.HandleFunc("/webhooks/pipefy/card-updated", cardsHandler.UpdateCard).Methods("POST")

	return r
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = pkg.NewRequestID()
		}

		ctx := pkg.WithRequestID(r.Context(), requestID)
		w.Header().Set("X-Request-ID", requestID)

		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))

		slog.InfoContext(ctx, "request",
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
