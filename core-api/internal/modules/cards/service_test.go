package cards

import (
	"context"
	"errors"
	"testing"

	"core-api/internal/providers/pipefy"
	db "core-api/services/database/generated"
)

func NewServiceTest(repo repositoryInterface, pip pipefyInterface) *Service {
	return &Service{repository: repo, pipefy: pip}
}

func TestUpdateCard_HappyPath(t *testing.T) {
	ctx := context.Background()

	dto := &CardUpdateDto{
		EventID:    "test-event-id",
		CardID:     "test-card-id",
		ClientMail: "teste@email.com",
		Timestamp:  "2024-06-01T12:00:00Z",
	}

	repo := &mockRepository{
		getEventsByClientEmail: func(_ context.Context, _ string) ([]db.Event, error) {
			return []db.Event{}, nil
		},
		registerEvent: func(_ context.Context, _ *CardUpdateDto) (db.Event, error) {
			return db.Event{EventID: dto.EventID}, nil
		},
		getClientByEmail: func(_ context.Context, _ string) (db.Client, error) {
			return db.Client{Email: dto.ClientMail, Amount: 100}, nil
		},
		updateClientStatus: func(_ context.Context, _ string, _ string) (db.Client, error) {
			return db.Client{Email: dto.ClientMail, Status: db.RequestStatusEnumPROCESSED}, nil
		},
	}

	pip := &mockPipefy{
		updateCardFields: func(_ context.Context, _ *pipefy.UpdateCardDto) error {
			return nil
		},
	}

	svc := NewServiceTest(repo, pip)

	if err := svc.UpdateCard(ctx, dto); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestUpdateCard_EventAlreadyRegistered(t *testing.T) {
	ctx := context.Background()

	dto := &CardUpdateDto{
		EventID:    "duplicate-event-id",
		CardID:     "test-card-id",
		ClientMail: "teste@email.com",
		Timestamp:  "2024-06-01T12:00:00Z",
	}

	repo := &mockRepository{
		getEventsByClientEmail: func(_ context.Context, _ string) ([]db.Event, error) {
			return []db.Event{{EventID: dto.EventID}}, nil
		},
	}

	svc := NewServiceTest(repo, &mockPipefy{})

	if err := svc.UpdateCard(ctx, dto); err != nil {
		t.Fatalf("expected no error on duplicate event, got: %v", err)
	}
}

func TestUpdateCard_RepositoryError(t *testing.T) {
	ctx := context.Background()

	dto := &CardUpdateDto{
		EventID:    "test-event-id",
		CardID:     "test-card-id",
		ClientMail: "teste@email.com",
		Timestamp:  "2024-06-01T12:00:00Z",
	}

	repo := &mockRepository{
		getEventsByClientEmail: func(_ context.Context, _ string) ([]db.Event, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewServiceTest(repo, &mockPipefy{})

	if err := svc.UpdateCard(ctx, dto); err == nil {
		t.Fatal("expected error, got nil")
	}
}
