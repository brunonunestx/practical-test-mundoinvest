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
	repository *ClientRepository
	pipefy     *pipefy.PipefyService
}

func NewClientService(pool *pgxpool.Pool, pipefy *pipefy.PipefyService) *ClientService {
	return &ClientService{
		repository: NewClientRepository(pool),
		pipefy:     pipefy,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, dto *CreateClientDto) (*Client, error) {
	priority := "LOW"

	if dto.Value > 200000 {
		priority = "HIGH"
	}

	createdUser, err := s.repository.CreateClient(ctx, &Client{
		name:          dto.Name,
		email:         dto.Email,
		clientType:    dto.ClientType,
		priority:      priority,
		heritageValue: pkg.DoubleToCents(dto.Value),
	})
	if err != nil {
		return nil, err
	}

	cfg := config.Load()

	_, err = s.pipefy.CreateCard(ctx, &pipefy.CreateCardDto{
		PipeId: cfg.PipeId,
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
		return nil, err
	}

	fmt.Println(createdUser)

	return createdUser, nil
}
