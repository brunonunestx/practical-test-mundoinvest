package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"core-api/internal/providers/pipefy"
	httputil "core-api/packages"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type clientServiceInterface interface {
	CreateClient(ctx context.Context, dto *CreateClientDto) (*Client, error)
}

type ClientHandler struct {
	service   clientServiceInterface
	validator *validator.Validate
}

func NewClientHandler(pool *pgxpool.Pool, pipefyService *pipefy.PipefyService) *ClientHandler {
	return &ClientHandler{
		service:   NewClientService(pool, pipefyService),
		validator: httputil.NewValidator(),
	}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto CreateClientDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if err := h.validator.Struct(dto); err != nil {
		fields, ok := httputil.ValidationErrors(err)
		if !ok {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"errors": fields})
		return
	}

	client, err := h.service.CreateClient(r.Context(), &dto)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create client"})
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, client.toResponse())
}
