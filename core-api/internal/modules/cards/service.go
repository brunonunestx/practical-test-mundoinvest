package cards

import (
	"context"
	"fmt"

	"core-api/internal/providers/pipefy"
	pkg "core-api/packages"
	db "core-api/services/database/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryInterface interface {
	GetEventsByClientEmail(ctx context.Context, email string) ([]db.Event, error)
	RegisterEvent(ctx context.Context, dto *CardUpdateDto) (db.Event, error)
	GetClientByEmail(ctx context.Context, email string) (db.Client, error)
	UpdateClientStatus(ctx context.Context, email string, status string) (db.Client, error)
}

type pipefyInterface interface {
	UpdateCardFields(ctx context.Context, dto *pipefy.UpdateCardDto) error
}

type Service struct {
	repository repositoryInterface
	pipefy     pipefyInterface
}

func NewService(pool *pgxpool.Pool, pipefySvc *pipefy.PipefyService) *Service {
	return &Service{
		repository: NewRepository(pool),
		pipefy:     pipefySvc,
	}
}

func (s *Service) UpdateCard(ctx context.Context, dto *CardUpdateDto) error {
	alreadyRegisteredEvents, err := s.repository.GetEventsByClientEmail(ctx, dto.ClientMail)
	if err != nil {
		return err
	}

	for _, event := range alreadyRegisteredEvents {
		if event.EventID == dto.EventID {
			fmt.Printf("event with ID %s already registered\n", dto.EventID)
			return nil
		}
	}

	registeredEvent, err := s.repository.RegisterEvent(ctx, dto)
	if err != nil {
		return err
	}

	fmt.Printf("registered event: %+v\n", registeredEvent)

	client, err := s.repository.GetClientByEmail(ctx, dto.ClientMail)
	if err != nil {
		return err
	}

	priority := "LOW"
	if pkg.CentsToDouble(client.Amount) >= 200000 {
		priority = "HIGH"
	}

	if err := s.pipefy.UpdateCardFields(ctx, &pipefy.UpdateCardDto{
		NodeId: dto.CardID,
		FieldsAttributes: []pipefy.FieldAttribute{
			{FieldId: "prioridade", Value: priority},
			{FieldId: "status", Value: "Processado"},
		},
	}); err != nil {
		return err
	}

	updatedClient, err := s.repository.UpdateClientStatus(ctx, dto.ClientMail, "PROCESSED")
	if err != nil {
		return err
	}

	fmt.Printf("updated client status to PROCESSED: %+v\n", updatedClient)
	fmt.Printf("retrieved client: %+v\n", client)

	return nil
}
