package clients

import (
	"context"
	"errors"
	"testing"

	"core-api/internal/providers/pipefy"
)

func newClientServiceTest(repo clientRepositoryInterface, pip clientPipefyInterface) *ClientService {
	return &ClientService{repository: repo, pipefy: pip, pipeId: 123456}
}

func TestCreateClient_LowPriority(t *testing.T) {
	ctx := context.Background()

	dto := &CreateClientDto{
		Name:       "João Silva",
		Email:      "joao@email.com",
		ClientType: "PF",
		Value:      100000,
	}

	var capturedClient *Client
	repo := &mockClientRepository{
		createClient: func(_ context.Context, c *Client) (*Client, error) {
			capturedClient = c
			return c, nil
		},
	}

	pip := &mockClientPipefy{
		createCard: func(_ context.Context, _ *pipefy.CreateCardDto) (*pipefy.Card, error) {
			return nil, nil
		},
	}

	svc := newClientServiceTest(repo, pip)

	if _, err := svc.CreateClient(ctx, dto); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if capturedClient.priority != "LOW" {
		t.Errorf("expected priority LOW, got %s", capturedClient.priority)
	}
}

func TestCreateClient_HighPriority(t *testing.T) {
	ctx := context.Background()

	dto := &CreateClientDto{
		Name:       "Maria Oliveira",
		Email:      "maria@email.com",
		ClientType: "PJ",
		Value:      250000,
	}

	var capturedCard *pipefy.CreateCardDto
	repo := &mockClientRepository{
		createClient: func(_ context.Context, c *Client) (*Client, error) {
			return c, nil
		},
	}

	pip := &mockClientPipefy{
		createCard: func(_ context.Context, d *pipefy.CreateCardDto) (*pipefy.Card, error) {
			capturedCard = d
			return nil, nil
		},
	}

	svc := newClientServiceTest(repo, pip)

	if _, err := svc.CreateClient(ctx, dto); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var priority string
	for _, f := range capturedCard.FieldsAttributes {
		if f.FieldId == "prioridade" {
			priority = f.Value
		}
	}

	if priority != "HIGH" {
		t.Errorf("expected priority HIGH in pipefy fields, got %s", priority)
	}
}

func TestCreateClient_RepositoryError(t *testing.T) {
	ctx := context.Background()

	dto := &CreateClientDto{
		Name:       "João Silva",
		Email:      "joao@email.com",
		ClientType: "PF",
		Value:      100000,
	}

	pipefyCalled := false
	repo := &mockClientRepository{
		createClient: func(_ context.Context, _ *Client) (*Client, error) {
			return nil, errors.New("db error")
		},
	}

	pip := &mockClientPipefy{
		createCard: func(_ context.Context, _ *pipefy.CreateCardDto) (*pipefy.Card, error) {
			pipefyCalled = true
			return nil, nil
		},
	}

	svc := newClientServiceTest(repo, pip)

	if _, err := svc.CreateClient(ctx, dto); err == nil {
		t.Fatal("expected error, got nil")
	}

	if pipefyCalled {
		t.Error("pipefy should not be called when repository fails")
	}
}

func TestCreateClient_PipefyError(t *testing.T) {
	ctx := context.Background()

	dto := &CreateClientDto{
		Name:       "João Silva",
		Email:      "joao@email.com",
		ClientType: "PF",
		Value:      100000,
	}

	repo := &mockClientRepository{
		createClient: func(_ context.Context, c *Client) (*Client, error) {
			return c, nil
		},
	}

	pip := &mockClientPipefy{
		createCard: func(_ context.Context, _ *pipefy.CreateCardDto) (*pipefy.Card, error) {
			return nil, errors.New("pipefy unavailable")
		},
	}

	svc := newClientServiceTest(repo, pip)

	if _, err := svc.CreateClient(ctx, dto); err == nil {
		t.Fatal("expected error when pipefy fails, got nil")
	}
}
