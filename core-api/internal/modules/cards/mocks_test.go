package cards

import (
	"context"

	"core-api/internal/providers/pipefy"
	db "core-api/services/database/generated"
)

type mockRepository struct {
	getEventsByClientEmail func(ctx context.Context, email string) ([]db.Event, error)
	registerEvent          func(ctx context.Context, dto *CardUpdateDto) (db.Event, error)
	getClientByEmail       func(ctx context.Context, email string) (db.Client, error)
	updateClientStatus     func(ctx context.Context, email string, status string) (db.Client, error)
}

func (m *mockRepository) GetEventsByClientEmail(ctx context.Context, email string) ([]db.Event, error) {
	return m.getEventsByClientEmail(ctx, email)
}

func (m *mockRepository) RegisterEvent(ctx context.Context, dto *CardUpdateDto) (db.Event, error) {
	return m.registerEvent(ctx, dto)
}

func (m *mockRepository) GetClientByEmail(ctx context.Context, email string) (db.Client, error) {
	return m.getClientByEmail(ctx, email)
}

func (m *mockRepository) UpdateClientStatus(ctx context.Context, email string, status string) (db.Client, error) {
	return m.updateClientStatus(ctx, email, status)
}

type mockPipefy struct {
	updateCardFields func(ctx context.Context, dto *pipefy.UpdateCardDto) error
}

func (m *mockPipefy) UpdateCardFields(ctx context.Context, dto *pipefy.UpdateCardDto) error {
	return m.updateCardFields(ctx, dto)
}
