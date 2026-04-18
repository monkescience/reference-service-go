//go:build integration

package integration_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monkescience/testastic"
)

var (
	postgresURL string        //nolint:gochecknoglobals // Shared test state via TestMain.
	testPool    *pgxpool.Pool //nolint:gochecknoglobals // Shared test state via TestMain.
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(0)
	}

	ctx := context.Background()
	exitCode := 1

	pg, err := startPostgres(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "starting postgres: %v\n", err)
		os.Exit(1)
	}

	cleanup := func() {
		if testPool != nil {
			testPool.Close()
		}

		if termErr := pg.Terminate(ctx); termErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "terminating postgres: %v\n", termErr)
		}
	}

	postgresURL, err = pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "getting postgres connection string: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	err = runMigrations(postgresURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running migrations: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	testPool, err = pgxpool.New(ctx, postgresURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "connecting to postgres: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	exitCode = testastic.CollectProcessCoverage(m, filepath.Join("..", "bin", "coverage.out"))
	cleanup()

	os.Exit(exitCode)
}
