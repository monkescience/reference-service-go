package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reference-service-go/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver registration for goose.
)

// Connect creates a pgx connection pool and runs migrations.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	err = runMigrations(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return pool, nil
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("opening migration connection: %w", err)
	}

	defer db.Close() //nolint:errcheck // Best-effort close after migrations.

	goose.SetBaseFS(migrations.FS)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("setting dialect: %w", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("running up migrations: %w", err)
	}

	return nil
}
