//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"text/template"
	"time"

	"github.com/monkescience/testastic"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const databaseEnvVar = "REFERENCE_SERVICE_TEST_DATABASE_URL"

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
	configPath, err := writeConfig("http://unused:0", 8080)
	if err != nil {
		return fmt.Errorf("writing migration config: %w", err)
	}

	defer os.Remove(configPath)

	cmd := exec.Command(
		"go", "run",
		"-trimpath",
		"-ldflags", "-X reference-service-go/internal/build.version=test",
		"./cmd/reference-service-go",
		"-config", configPath, "-migrate-only",
	)
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), databaseEnvVar+"="+databaseURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}

func writeConfig(pokeapiURL string, port int) (string, error) {
	tmpl, err := template.ParseFiles("testdata/config/config.yaml.tmpl")
	if err != nil {
		return "", fmt.Errorf("parse config template: %w", err)
	}

	f, err := os.CreateTemp("", "e2e-config-*.yaml")
	if err != nil {
		return "", fmt.Errorf("creating temp config: %w", err)
	}

	err = tmpl.Execute(f, struct {
		DatabaseEnvVar string
		PokeAPIURL     string
		Port           int
	}{
		DatabaseEnvVar: databaseEnvVar,
		PokeAPIURL:     pokeapiURL,
		Port:           port,
	})
	if err != nil {
		f.Close()

		return "", fmt.Errorf("rendering config: %w", err)
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
	t.Setenv(databaseEnvVar, postgresURL)

	configPath, err := writeConfig(pokeapiURL, port)
	if err != nil {
		t.Fatalf("writing config: %v", err)
	}

	t.Cleanup(func() { os.Remove(configPath) })

	return testastic.StartProcess(t.Context(), t,
		"reference-service-go/cmd/reference-service-go",
		testastic.HTTPCheck(port, "/health/ready"),
		testastic.WithPort(port),
		testastic.WithArgs("-config", configPath),
		testastic.WithBuildArgs("-trimpath", "-ldflags", "-X reference-service-go/internal/build.version=test"),
		testastic.WithReadyTimeout(10*time.Second),
	)
}

func truncateTables(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE catches, pokemon, imports")
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}
}
