package clients

import (
	"context"

	"core-api/internal/providers/pipefy"
)

type mockClientRepository struct {
	createClient func(ctx context.Context, client *Client) (*Client, error)
}

func (m *mockClientRepository) CreateClient(ctx context.Context, client *Client) (*Client, error) {
	return m.createClient(ctx, client)
}

type mockClientPipefy struct {
	createCard func(ctx context.Context, dto *pipefy.CreateCardDto) (*pipefy.Card, error)
}

func (m *mockClientPipefy) CreateCard(ctx context.Context, dto *pipefy.CreateCardDto) (*pipefy.Card, error) {
	return m.createCard(ctx, dto)
}
