//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"core-api/internal/providers/pipefy"
	"core-api/internal/server"
	db "core-api/services/database/generated"
	"core-api/tests/helpers"
)

func TestWebhookPriority_Integration(t *testing.T) {
	cases := []struct {
		name         string
		clientBody   map[string]any
		wantPriority string
	}{
		{
			name: "HIGH priority when patrimônio >= 200k",
			clientBody: map[string]any{
				"cliente_nome":     "Carlos Alto",
				"cliente_email":    "carlos@example.com",
				"tipo_solicitacao": "PF",
				"valor_patrimonio": 300000.0,
			},
			wantPriority: "HIGH",
		},
		{
			name: "LOW priority when patrimônio < 200k",
			clientBody: map[string]any{
				"cliente_nome":     "Ana Baixo",
				"cliente_email":    "ana@example.com",
				"tipo_solicitacao": "PJ",
				"valor_patrimonio": 100000.0,
			},
			wantPriority: "LOW",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := helpers.NewTestDB(t)

			var capturedPriority string
			mock := helpers.NewMockPipefyService()
			mock.UpdateCardFieldsFn = func(_ context.Context, dto *pipefy.UpdateCardDto) error {
				for _, f := range dto.FieldsAttributes {
					if f.FieldId == "prioridade" {
						capturedPriority = f.Value
					}
				}
				return nil
			}

			s := server.NewServerWithDeps(pool, mock)
			handler := s.RegisterRoutes()

			rawClient, _ := json.Marshal(tc.clientBody)
			req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader(rawClient))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusCreated {
				t.Fatalf("create client: want 201, got %d — %s", rec.Code, rec.Body.String())
			}

			webhookBody := map[string]any{
				"event_id":      "evt-priority-001",
				"card_id":       "card-001",
				"cliente_email": tc.clientBody["cliente_email"],
				"timestamp":     "2024-01-01T00:00:00Z",
			}
			rawWebhook, _ := json.Marshal(webhookBody)
			req2 := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(rawWebhook))
			req2.Header.Set("Content-Type", "application/json")
			rec2 := httptest.NewRecorder()
			handler.ServeHTTP(rec2, req2)
			if rec2.Code != http.StatusOK {
				t.Fatalf("webhook: want 200, got %d — %s", rec2.Code, rec2.Body.String())
			}

			if capturedPriority != tc.wantPriority {
				t.Errorf("priority sent to pipefy: want %q, got %q", tc.wantPriority, capturedPriority)
			}

			queries := db.New(pool)
			email := tc.clientBody["cliente_email"].(string)
			client, err := queries.GetClientByEmail(context.Background(), email)
			if err != nil {
				t.Fatalf("get client: %v", err)
			}
			if client.Status != db.RequestStatusEnumPROCESSED {
				t.Errorf("client status: want PROCESSED, got %q", client.Status)
			}
		})
	}
}

func TestWebhookDuplicateEvent_Integration(t *testing.T) {
	pool := helpers.NewTestDB(t)

	updateCallCount := 0
	mock := helpers.NewMockPipefyService()
	mock.UpdateCardFieldsFn = func(_ context.Context, _ *pipefy.UpdateCardDto) error {
		updateCallCount++
		return nil
	}

	s := server.NewServerWithDeps(pool, mock)
	handler := s.RegisterRoutes()

	clientBody := map[string]any{
		"cliente_nome":     "Evento Duplicado",
		"cliente_email":    "dup@example.com",
		"tipo_solicitacao": "PF",
		"valor_patrimonio": 500000.0,
	}
	rawClient, _ := json.Marshal(clientBody)
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader(rawClient))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create client: want 201, got %d — %s", rec.Code, rec.Body.String())
	}

	webhookBody := map[string]any{
		"event_id":      "evt-dup-001",
		"card_id":       "card-001",
		"cliente_email": "dup@example.com",
		"timestamp":     "2024-01-01T00:00:00Z",
	}

	// First webhook — should process normally
	rawWebhook, _ := json.Marshal(webhookBody)
	req1 := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(rawWebhook))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first webhook: want 200, got %d — %s", rec1.Code, rec1.Body.String())
	}
	if updateCallCount != 1 {
		t.Fatalf("after first webhook: pipefy call count: want 1, got %d", updateCallCount)
	}

	// Second webhook with same event_id — should be skipped silently
	rawWebhook2, _ := json.Marshal(webhookBody)
	req2 := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewReader(rawWebhook2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("duplicate webhook: want 200, got %d — %s", rec2.Code, rec2.Body.String())
	}
	if updateCallCount != 1 {
		t.Errorf("after duplicate webhook: pipefy call count: want still 1, got %d", updateCallCount)
	}

	queries := db.New(pool)
	events, err := queries.GetEventsByClientEmail(context.Background(), "dup@example.com")
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("events in DB: want 1, got %d", len(events))
	}
}
