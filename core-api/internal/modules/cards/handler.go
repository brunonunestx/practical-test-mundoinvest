package cards

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"core-api/internal/providers/pipefy"
	httputil "core-api/packages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cardServiceInterface interface {
	UpdateCard(ctx context.Context, dto *CardUpdateDto) error
}

type Handler struct {
	service cardServiceInterface
}

func NewHandler(pool *pgxpool.Pool, pipefy pipefy.Provider) *Handler {
	return &Handler{service: NewService(pool, pipefy)}
}

func (h *Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var dto CardUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if err := httputil.NewValidator().Struct(dto); err != nil {
		fields, ok := httputil.ValidationErrors(err)
		if !ok {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"errors": fields})
		return
	}

	if err := h.service.UpdateCard(r.Context(), &dto); err != nil {
		fmt.Printf("Error updating card: %v\n", err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update card"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "card updated successfully"})
}
