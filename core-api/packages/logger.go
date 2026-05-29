package pkg

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromCtx(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func Logger(ctx context.Context) *slog.Logger {
	if id := RequestIDFromCtx(ctx); id != "" {
		return slog.With("request_id", id)
	}
	return slog.Default()
}

func NewRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
