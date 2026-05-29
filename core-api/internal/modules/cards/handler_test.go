package cards

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httputil "core-api/packages"
)

type mockCardService struct {
	updateCard func(ctx context.Context, dto *CardUpdateDto) error
}

func (m *mockCardService) UpdateCard(ctx context.Context, dto *CardUpdateDto) error {
	return m.updateCard(ctx, dto)
}

func newCardHandlerTest(svc cardServiceInterface) *Handler {
	return &Handler{service: svc}
}

func TestUpdateCardHandler_Success(t *testing.T) {
	svc := &mockCardService{
		updateCard: func(_ context.Context, _ *CardUpdateDto) error {
			return nil
		},
	}

	body := `{"event_id":"evt-1","card_id":"card-1","cliente_email":"joao@email.com","timestamp":"2024-06-01T12:00:00Z"}`
	r := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(body))
	w := httptest.NewRecorder()

	newCardHandlerTest(svc).UpdateCard(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateCardHandler_InvalidJSON(t *testing.T) {
	svc := &mockCardService{}

	r := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(`{invalid`))
	w := httptest.NewRecorder()

	newCardHandlerTest(svc).UpdateCard(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateCardHandler_ValidationError(t *testing.T) {
	svc := &mockCardService{}

	cases := []struct {
		name string
		body string
	}{
		{"missing event_id", `{"card_id":"card-1","cliente_email":"joao@email.com","timestamp":"2024-06-01T12:00:00Z"}`},
		{"missing card_id", `{"event_id":"evt-1","cliente_email":"joao@email.com","timestamp":"2024-06-01T12:00:00Z"}`},
		{"invalid email", `{"event_id":"evt-1","card_id":"card-1","cliente_email":"not-an-email","timestamp":"2024-06-01T12:00:00Z"}`},
		{"missing timestamp", `{"event_id":"evt-1","card_id":"card-1","cliente_email":"joao@email.com"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			newCardHandlerTest(svc).UpdateCard(w, r)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestUpdateCardHandler_ServiceError(t *testing.T) {
	svc := &mockCardService{
		updateCard: func(_ context.Context, _ *CardUpdateDto) error {
			return errors.New("service error")
		},
	}

	body := `{"event_id":"evt-1","card_id":"card-1","cliente_email":"joao@email.com","timestamp":"2024-06-01T12:00:00Z"}`
	r := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(body))
	w := httptest.NewRecorder()

	_ = httputil.NewValidator()
	newCardHandlerTest(svc).UpdateCard(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
