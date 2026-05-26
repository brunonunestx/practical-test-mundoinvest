package server

import (
	"net/http"

	"core-api/internal/modules/clients"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	clientsHandler := clients.NewClientHandler()

	r := mux.NewRouter()

	r.Use(s.corsMiddleware)

	r.HandleFunc("/clients", clientsHandler.CreateClient).Methods("POST")

	return r
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // When deploy to production, change this to the actual domain
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
