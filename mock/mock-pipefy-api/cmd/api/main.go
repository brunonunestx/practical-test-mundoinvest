package main

import (
	"log/slog"
	"os"

	"mock-pipefy-api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	slog.Info("mock pipefy api starting")
	defer slog.Info("mock pipefy api stopped")

	if err := server.Init().Start(); err != nil {
		slog.Error("mock pipefy api failed", "error", err)
	}
}
