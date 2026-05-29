package pipefy

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"core-api/internal/providers/config"
)

type Provider interface {
	CreateCard(ctx context.Context, dto *CreateCardDto) (*Card, error)
	UpdateCardFields(ctx context.Context, dto *UpdateCardDto) error
}

type PipefyService struct{}

func NewPipefyService() *PipefyService {
	return &PipefyService{}
}

type graphQLRequest struct {
	Query string `json:"query"`
}

func (s *PipefyService) CreateCard(ctx context.Context, dto *CreateCardDto) (*Card, error) {
	cfg := config.Load()
	mutation := BuildCreateCardMutation(dto)

	slog.DebugContext(ctx, "pipefy create card mutation", "mutation", mutation)

	body, err := json.Marshal(graphQLRequest{Query: mutation})
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "pipefy create card", "url", cfg.PipefyApiUrl, "pipe_id", dto.PipeId, "title", dto.Title)

	resp, err := http.Post(cfg.PipefyApiUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.ErrorContext(ctx, "pipefy create card request failed", "url", cfg.PipefyApiUrl, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	slog.InfoContext(ctx, "pipefy create card response", "status", resp.StatusCode)
	return nil, nil
}

func (s *PipefyService) UpdateCardFields(ctx context.Context, dto *UpdateCardDto) error {
	cfg := config.Load()
	mutation := BuildUpdateCardFieldsMutation(dto.NodeId, dto.FieldsAttributes)

	slog.DebugContext(ctx, "pipefy update card fields mutation", "mutation", mutation)

	body, err := json.Marshal(graphQLRequest{Query: mutation})
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "pipefy update card fields", "url", cfg.PipefyApiUrl, "card_id", dto.NodeId)

	resp, err := http.Post(cfg.PipefyApiUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.ErrorContext(ctx, "pipefy update card fields request failed", "url", cfg.PipefyApiUrl, "card_id", dto.NodeId, "error", err)
		return err
	}
	defer resp.Body.Close()

	slog.InfoContext(ctx, "pipefy update card fields response", "status", resp.StatusCode)
	return nil
}
