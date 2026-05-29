package clients

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httputil "core-api/packages"
)

type mockClientService struct {
	createClient func(ctx context.Context, dto *CreateClientDto) (*Client, error)
}

func (m *mockClientService) CreateClient(ctx context.Context, dto *CreateClientDto) (*Client, error) {
	return m.createClient(ctx, dto)
}

func newHandlerTest(svc clientServiceInterface) *ClientHandler {
	return &ClientHandler{service: svc, validator: httputil.NewValidator()}
}

func TestCreateClientHandler_Success(t *testing.T) {
	svc := &mockClientService{
		createClient: func(_ context.Context, _ *CreateClientDto) (*Client, error) {
			return &Client{name: "João Silva", email: "joao@email.com", priority: "LOW"}, nil
		},
	}

	body := `{"cliente_nome":"João Silva","cliente_email":"joao@email.com","tipo_solicitacao":"PF","valor_patrimonio":100000}`
	r := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(body))
	w := httptest.NewRecorder()

	newHandlerTest(svc).CreateClient(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCreateClientHandler_InvalidJSON(t *testing.T) {
	svc := &mockClientService{}

	r := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(`{invalid`))
	w := httptest.NewRecorder()

	newHandlerTest(svc).CreateClient(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateClientHandler_ValidationError(t *testing.T) {
	svc := &mockClientService{}

	cases := []struct {
		name string
		body string
	}{
		{"missing name", `{"cliente_email":"joao@email.com","tipo_solicitacao":"PF","valor_patrimonio":100000}`},
		{"invalid email", `{"cliente_nome":"João","cliente_email":"not-an-email","tipo_solicitacao":"PF","valor_patrimonio":100000}`},
		{"email without domain", `{"cliente_nome":"João","cliente_email":"joao@","tipo_solicitacao":"PF","valor_patrimonio":100000}`},
		{"email without @", `{"cliente_nome":"João","cliente_email":"joaoemail.com","tipo_solicitacao":"PF","valor_patrimonio":100000}`},
		{"zero value", `{"cliente_nome":"João","cliente_email":"joao@email.com","tipo_solicitacao":"PF","valor_patrimonio":0}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			newHandlerTest(svc).CreateClient(w, r)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestCreateClientHandler_ServiceError(t *testing.T) {
	svc := &mockClientService{
		createClient: func(_ context.Context, _ *CreateClientDto) (*Client, error) {
			return nil, errors.New("db error")
		},
	}

	body := `{"cliente_nome":"João Silva","cliente_email":"joao@email.com","tipo_solicitacao":"PF","valor_patrimonio":100000}`
	r := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(body))
	w := httptest.NewRecorder()

	newHandlerTest(svc).CreateClient(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
