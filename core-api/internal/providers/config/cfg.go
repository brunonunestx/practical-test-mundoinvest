package config

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DatabaseURL    string
	PipefyApiUrl   string
	PipefyApiToken string
	PipeId         int
}

func Load() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pipefyApiUrl := os.Getenv("PIPEFY_API_URL")
	if pipefyApiUrl == "" {
		log.Fatal("PIPEFY_API_URL is required")
	}

	pipefyApiToken := os.Getenv("PIPEFY_API_TOKEN")
	if pipefyApiToken == "" {
		log.Fatal("PIPEFY_API_TOKEN is required")
	}

	pipeIdStr := os.Getenv("PIPE_ID")
	if pipeIdStr == "" {
		log.Fatal("PIPE_ID is required")
	}
	pipeId, err := strconv.Atoi(pipeIdStr)
	if err != nil {
		log.Fatalf("PIPE_ID must be a valid integer: %v", err)
	}

	return &Config{
		DatabaseURL:    databaseURL,
		PipefyApiUrl:   pipefyApiUrl,
		PipefyApiToken: pipefyApiToken,
		PipeId:         pipeId,
	}
}
