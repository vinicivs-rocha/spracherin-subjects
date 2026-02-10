package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/vinicivs-rocha/spracherin-subjects/internal/data"
	grpcserver "github.com/vinicivs-rocha/spracherin-subjects/internal/transport/grpc"
)

type noopMessager struct{}

func (noopMessager) MessageDetectedChanges(ctx context.Context, changes data.SubjectChanges) error {
	return nil
}

func main() {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	maxOpen := envInt("DB_MAX_OPEN", 10)
	maxIdle := envInt("DB_MAX_IDLE", 5)

	repo, err := data.NewMySQLSubjectRepositoryFromEnv(data.MySQLPoolConfig{
		MaxOpen: maxOpen,
		MaxIdle: maxIdle,
	})
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	if err := grpcserver.ListenAndServe(addr, repo, noopMessager{}); err != nil {
		log.Fatalf("grpc server stopped: %v", err)
	}
}

func envInt(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	return value
}
