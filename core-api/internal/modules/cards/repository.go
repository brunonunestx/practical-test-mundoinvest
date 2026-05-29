package cards

import (
	"context"
	"fmt"
	"time"

	db "core-api/services/database/generated"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	queries := db.New(pool)
	return &Repository{queries: queries}
}

func (r *Repository) RegisterEvent(ctx context.Context, dto *CardUpdateDto) (db.Event, error) {
	t, err := time.Parse(time.RFC3339, dto.Timestamp)
	if err != nil {
		return db.Event{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	createdEvent, err := r.queries.CreateEvent(ctx, db.CreateEventParams{
		EventID:     dto.EventID,
		CardID:      dto.CardID,
		ClientEmail: dto.ClientMail,
		Timestamp:   pgtype.Timestamptz{Time: t.UTC(), Valid: true},
	})

	return createdEvent, err
}

func (r *Repository) GetEventsByClientEmail(ctx context.Context, email string) ([]db.Event, error) {
	events, err := r.queries.GetEventsByClientEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repository) GetClientByEmail(ctx context.Context, email string) (db.Client, error) {
	client, err := r.queries.GetClientByEmail(ctx, email)
	if err != nil {
		return db.Client{}, err
	}

	return client, nil
}

func (r *Repository) UpdateClientStatus(ctx context.Context, email string, status string) (db.Client, error) {
	updatedClient, err := r.queries.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
		Email:  email,
		Status: db.RequestStatusEnum(status),
	})
	if err != nil {
		return db.Client{}, err
	}

	return updatedClient, nil
}
