package cards

import (
	"context"
	"fmt"

	"core-api/internal/providers/pipefy"
	pkg "core-api/packages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repository *Repository
	pipefy     *pipefy.PipefyService
}

func NewService(pool *pgxpool.Pool, pipefy *pipefy.PipefyService) *Service {
	return &Service{
		repository: NewRepository(pool),
		pipefy:     pipefy,
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
