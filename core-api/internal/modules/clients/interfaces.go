package clients

import (
	"context"

	"core-api/internal/providers/pipefy"
)

type clientRepositoryInterface interface {
	CreateClient(ctx context.Context, client *Client) (*Client, error)
}

type clientPipefyInterface interface {
	CreateCard(ctx context.Context, dto *pipefy.CreateCardDto) (*pipefy.Card, error)
}
