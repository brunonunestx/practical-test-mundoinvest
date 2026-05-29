//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"core-api/internal/server"
	"core-api/tests/helpers"
	db "core-api/services/database/generated"
)

func TestMain(m *testing.M) {
	// config.Load() é chamado dentro de NewClientService; essas vars não são usadas
	// de fato pois o banco e o pipefy são substituídos nos testes.
	os.Setenv("DATABASE_URL", "unused-overridden-by-testcontainer")
	os.Setenv("PIPEFY_API_URL", "http://localhost")
	os.Setenv("PIPEFY_API_TOKEN", "test-token")
	os.Setenv("PIPE_ID", "1")
	os.Exit(m.Run())
}

func TestCreateClient_Integration(t *testing.T) {
	cases := []struct {
		name             string
		body             map[string]any
		wantStatus       int
		wantPriority     db.PriorityEnum
		wantAmountCents  int32
	}{
		{
			name: "high priority when value above 200k",
			body: map[string]any{
				"cliente_nome":     "João Silva",
				"cliente_email":    "joao@example.com",
				"tipo_solicitacao": "PF",
				"valor_patrimonio": 300000.0,
			},
			wantStatus:      http.StatusCreated,
			wantPriority:    db.PriorityEnumHIGH,
			wantAmountCents: 30000000,
		},
		{
			name: "low priority when value below 200k",
			body: map[string]any{
				"cliente_nome":     "Maria Souza",
				"cliente_email":    "maria@example.com",
				"tipo_solicitacao": "PJ",
				"valor_patrimonio": 100000.0,
			},
			wantStatus:      http.StatusCreated,
			wantPriority:    db.PriorityEnumLOW,
			wantAmountCents: 10000000,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := helpers.NewTestDB(t)
			mock := helpers.NewMockPipefyService()

			s := server.NewServerWithDeps(pool, mock)
			handler := s.RegisterRoutes()

			rawBody, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader(rawBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("HTTP status: want %d, got %d — body: %s", tc.wantStatus, rec.Code, rec.Body.String())
			}

			// Verifica a resposta HTTP
			var resp map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if resp["email"] != tc.body["cliente_email"] {
				t.Errorf("response email: want %q, got %q", tc.body["cliente_email"], resp["email"])
			}

			// Verifica o estado real no banco
			queries := db.New(pool)
			email := tc.body["cliente_email"].(string)
			client, err := queries.GetClientByEmail(context.Background(), email)
			if err != nil {
				t.Fatalf("client not found in database: %v", err)
			}

			if client.Name != tc.body["cliente_nome"] {
				t.Errorf("name: want %q, got %q", tc.body["cliente_nome"], client.Name)
			}
			if client.Email != email {
				t.Errorf("email: want %q, got %q", email, client.Email)
			}
			if client.RequestType != tc.body["tipo_solicitacao"] {
				t.Errorf("request_type: want %q, got %q", tc.body["tipo_solicitacao"], client.RequestType)
			}
			if client.Status != db.RequestStatusEnumPENDINGANALYSIS {
				t.Errorf("status: want PENDING_ANALYSIS, got %q", client.Status)
			}
			if client.Priority != tc.wantPriority {
				t.Errorf("priority: want %q, got %q", tc.wantPriority, client.Priority)
			}
			if client.Amount != tc.wantAmountCents {
				t.Errorf("amount: want %d cents, got %d", tc.wantAmountCents, client.Amount)
			}
		})
	}
}
