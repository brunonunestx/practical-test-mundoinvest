package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"core-api/internal/providers/config"
	"core-api/internal/providers/pipefy"
)

type Server struct {
	port   int
	pool   *pgxpool.Pool
	pipefy pipefy.Provider
}

func NewServerWithDeps(pool *pgxpool.Pool, pipefySvc pipefy.Provider) *Server {
	return &Server{
		pool:   pool,
		pipefy: pipefySvc,
	}
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	newServer := &Server{
		port:   port,
		pool:   pool,
		pipefy: pipefy.NewPipefyService(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
