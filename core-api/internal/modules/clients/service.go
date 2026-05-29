package clients

import (
	"context"
	"fmt"

	"core-api/internal/providers/config"
	"core-api/internal/providers/pipefy"
	pkg "core-api/packages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientService struct {
	repository clientRepositoryInterface
	pipefy     clientPipefyInterface
	pipeId     int
}

func NewClientService(pool *pgxpool.Pool, pipefy pipefy.Provider) *ClientService {
	cfg := config.Load()
	return &ClientService{
		repository: NewClientRepository(pool),
		pipefy:     pipefy,
		pipeId:     cfg.PipeId,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, dto *CreateClientDto) (*Client, error) {
	logger := pkg.Logger(ctx)

	priority := "LOW"
	if dto.Value > 200000 {
		priority = "HIGH"
	}
	logger.Info("priority assigned", "email", dto.Email, "priority", priority, "value", dto.Value)

	createdUser, err := s.repository.CreateClient(ctx, &Client{
		name:          dto.Name,
		email:         dto.Email,
		clientType:    dto.ClientType,
		priority:      priority,
		heritageValue: pkg.DoubleToCents(dto.Value),
	})
	if err != nil {
		logger.Error("save client failed", "email", dto.Email, "error", err)
		return nil, err
	}
	logger.Info("client saved", "email", dto.Email)

	_, err = s.pipefy.CreateCard(ctx, &pipefy.CreateCardDto{
		PipeId: s.pipeId,
		Title:  dto.Email,
		FieldsAttributes: []pipefy.FieldAttribute{
			{FieldId: "nome", Value: dto.Name},
			{FieldId: "email", Value: dto.Email},
			{FieldId: "tipo_cliente", Value: dto.ClientType},
			{FieldId: "valor_patrimonio", Value: fmt.Sprintf("%g", dto.Value)},
			{FieldId: "prioridade", Value: priority},
			{FieldId: "status", Value: "Aguardando Análise"},
		},
	})
	if err != nil {
		logger.Error("pipefy create card failed", "email", dto.Email, "error", err)
		return nil, err
	}
	logger.Info("pipefy card created", "email", dto.Email)

	return createdUser, nil
}
