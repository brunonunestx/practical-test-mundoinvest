package server

import (
	"net/http"

	"core-api/internal/modules/cards"
	"core-api/internal/modules/clients"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	clientsHandler := clients.NewClientHandler(s.pool, s.pipefy)
	cardsHandler := cards.NewHandler(s.pool, s.pipefy)

	r := mux.NewRouter()

	r.Use(s.corsMiddleware)

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
