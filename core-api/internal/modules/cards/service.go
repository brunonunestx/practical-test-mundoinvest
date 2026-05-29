package cards

import (
	"context"
	"log/slog"

	"core-api/internal/providers/pipefy"
	pkg "core-api/packages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repository cardRepositoryInterface
	pipefy     cardPipefyInterface
}

func NewService(pool *pgxpool.Pool, pipefySvc pipefy.Provider) *Service {
	return &Service{
		repository: NewRepository(pool),
		pipefy:     pipefySvc,
	}
}

func (s *Service) UpdateCard(ctx context.Context, dto *CardUpdateDto) error {
	logger := pkg.Logger(ctx)

	alreadyRegisteredEvents, err := s.repository.GetEventsByClientEmail(ctx, dto.ClientMail)
	if err != nil {
		logger.Error("fetch events failed", "client_email", dto.ClientMail, "error", err)
		return err
	}

	for _, event := range alreadyRegisteredEvents {
		if event.EventID == dto.EventID {
			slog.WarnContext(ctx, "duplicate event skipped", "event_id", dto.EventID, "client_email", dto.ClientMail)
			return nil
		}
	}

	registeredEvent, err := s.repository.RegisterEvent(ctx, dto)
	if err != nil {
		logger.Error("register event failed", "event_id", dto.EventID, "client_email", dto.ClientMail, "error", err)
		return err
	}
	logger.Info("event registered", "event_id", registeredEvent.EventID, "card_id", registeredEvent.CardID)

	client, err := s.repository.GetClientByEmail(ctx, dto.ClientMail)
	if err != nil {
		logger.Error("fetch client failed", "client_email", dto.ClientMail, "error", err)
		return err
	}

	priority := "LOW"
	if pkg.CentsToDouble(client.Amount) >= 200000 {
		priority = "HIGH"
	}
	logger.Info("priority resolved", "client_email", dto.ClientMail, "priority", priority)

	if err := s.pipefy.UpdateCardFields(ctx, &pipefy.UpdateCardDto{
		NodeId: dto.CardID,
		FieldsAttributes: []pipefy.FieldAttribute{
			{FieldId: "prioridade", Value: priority},
			{FieldId: "status", Value: "Processado"},
		},
	}); err != nil {
		logger.Error("pipefy update card failed", "card_id", dto.CardID, "error", err)
		return err
	}

	if _, err := s.repository.UpdateClientStatus(ctx, dto.ClientMail, "PROCESSED"); err != nil {
		logger.Error("update client status failed", "client_email", dto.ClientMail, "error", err)
		return err
	}

	logger.Info("client status updated", "client_email", dto.ClientMail, "status", "PROCESSED")

	return nil
}
