package testutil

import (
	"context"
	"fmt"
	"reference-service-go/internal/outgoing/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer holds a running testcontainer and its connection URL.
type PostgresContainer struct {
	Container testcontainers.Container
	URL       string
}

// StartPostgres creates a PostgreSQL testcontainer for integration tests.
func StartPostgres(ctx context.Context) (*PostgresContainer, error) {
	container, err := tcpostgres.Run(ctx,
		"postgres:17-alpine",
		tcpostgres.WithDatabase("pokemon_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("starting postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("getting connection string: %w", err)
	}

	return &PostgresContainer{
		Container: container,
		URL:       connStr,
	}, nil
}

// SetupDatabase connects to the testcontainer and runs migrations.
func SetupDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, *postgres.Queries, error) {
	pool, err := postgres.Connect(ctx, databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to test database: %w", err)
	}

	queries := postgres.New(pool)

	return pool, queries, nil
}

// TruncateTables resets all application tables between tests.
func TruncateTables(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(ctx, "TRUNCATE TABLE pokemon, imports")
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}
}
