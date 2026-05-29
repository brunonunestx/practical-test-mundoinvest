package pipefy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"core-api/internal/providers/config"
)

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

	fmt.Printf("Sending mutation to Pipefy: %s\n", mutation)
	fmt.Printf("Using API URL: %s\n", cfg.PipefyApiUrl)

	body, err := json.Marshal(graphQLRequest{Query: mutation})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(cfg.PipefyApiUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return nil, nil
}
