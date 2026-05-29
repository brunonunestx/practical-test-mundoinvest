package helpers

import (
	"context"

	"core-api/internal/providers/pipefy"
)

type MockPipefyService struct {
	CreateCardFn       func(ctx context.Context, dto *pipefy.CreateCardDto) (*pipefy.Card, error)
	UpdateCardFieldsFn func(ctx context.Context, dto *pipefy.UpdateCardDto) error
}

func NewMockPipefyService() *MockPipefyService {
	return &MockPipefyService{
		CreateCardFn: func(_ context.Context, _ *pipefy.CreateCardDto) (*pipefy.Card, error) {
			return &pipefy.Card{}, nil
		},
		UpdateCardFieldsFn: func(_ context.Context, _ *pipefy.UpdateCardDto) error {
			return nil
		},
	}
}

func (m *MockPipefyService) CreateCard(ctx context.Context, dto *pipefy.CreateCardDto) (*pipefy.Card, error) {
	return m.CreateCardFn(ctx, dto)
}

func (m *MockPipefyService) UpdateCardFields(ctx context.Context, dto *pipefy.UpdateCardDto) error {
	return m.UpdateCardFieldsFn(ctx, dto)
}
