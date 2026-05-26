package clients

import (
	"encoding/json"
	"net/http"
)

type ClientHandler struct {
	Service *ClientService
}

func NewClientHandler() *ClientHandler {
	service := NewClientService()
	return &ClientHandler{
		Service: service,
	}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client CreateClientDto
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdClient, err := h.Service.CreateClient(&client)
	if err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdClient)
}
