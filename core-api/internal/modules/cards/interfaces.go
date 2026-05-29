package cards

import (
	"context"

	"core-api/internal/providers/pipefy"
	db "core-api/services/database/generated"
)

type cardRepositoryInterface interface {
	GetEventsByClientEmail(ctx context.Context, email string) ([]db.Event, error)
	RegisterEvent(ctx context.Context, dto *CardUpdateDto) (db.Event, error)
	GetClientByEmail(ctx context.Context, email string) (db.Client, error)
	UpdateClientStatus(ctx context.Context, email string, status string) (db.Client, error)
}

type cardPipefyInterface interface {
	UpdateCardFields(ctx context.Context, dto *pipefy.UpdateCardDto) error
}
