//go:build integration

package tests_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monkescience/testastic"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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

	pg, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("starting postgres: %v", err)
	}

	defer func() {
		if termErr := pg.Terminate(ctx); termErr != nil {
			log.Printf("terminating postgres: %v", termErr)
		}
	}()

	postgresURL, err = pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("getting postgres connection string: %v", err)
	}

	err = runMigrations(postgresURL)
	if err != nil {
		log.Fatalf("running migrations: %v", err)
	}

	testPool, err = pgxpool.New(ctx, postgresURL)
	if err != nil {
		log.Fatalf("connecting to postgres: %v", err)
	}

	defer testPool.Close()

	os.Exit(testastic.CollectProcessCoverage(m, filepath.Join("..", "bin", "coverage.out")))
}

func startPostgres(ctx context.Context) (*tcpostgres.PostgresContainer, error) {
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

	return container, nil
}

func runMigrations(databaseURL string) error {
	configPath, err := writeConfig(databaseURL, "http://unused:0", 0)
	if err != nil {
		return fmt.Errorf("writing migration config: %w", err)
	}

	defer os.Remove(configPath)

	cmd := exec.Command(
		"go", "run",
		"-ldflags", "-X reference-service-go/internal/build.Version=test",
		"./cmd/reference-service-go",
		"-config", configPath, "-migrate-only",
	)
	cmd.Dir = ".."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}

func writeConfig(databaseURL string, pokeapiURL string, port int) (string, error) {
	content := fmt.Sprintf(`log_config:
  level: "debug"
  format: "text"
  add_source: false

server:
  port: %d

database:
  url: "%s"

pokeapi:
  base_url: "%s"
  timeout: "30s"
  concurrency: 5
`, port, databaseURL, pokeapiURL)

	f, err := os.CreateTemp("", "e2e-config-*.yaml")
	if err != nil {
		return "", fmt.Errorf("creating temp config: %w", err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		f.Close()

		return "", fmt.Errorf("writing config: %w", err)
	}

	err = f.Close()
	if err != nil {
		return "", fmt.Errorf("closing config: %w", err)
	}

	return f.Name(), nil
}

func findFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding free port: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	listener.Close()

	return port
}

func startService(t *testing.T, pokeapiURL string) *testastic.Process {
	t.Helper()

	port := findFreePort(t)

	configPath, err := writeConfig(postgresURL, pokeapiURL, port)
	if err != nil {
		t.Fatalf("writing config: %v", err)
	}

	t.Cleanup(func() { os.Remove(configPath) })

	return testastic.StartProcess(t.Context(), t,
		"reference-service-go/cmd/reference-service-go",
		testastic.HTTPCheck(port, "/pokemon"),
		testastic.WithPort(port),
		testastic.WithArgs("-config", configPath),
		testastic.WithBuildArgs("-ldflags", "-X reference-service-go/internal/build.Version=test"),
		testastic.WithReadyTimeout(10*time.Second),
	)
}

func truncateTables(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE pokemon, imports")
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}
}
