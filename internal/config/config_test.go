package config_test

import (
	"os"
	"reference-service-go/internal/config"
	"testing"

	"github.com/monkescience/testastic"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("loads valid config", func(t *testing.T) {
		t.Parallel()

		// when: loading a valid config file
		cfg, err := config.Load("../../config/config.yaml")

		// then: it returns the parsed config
		testastic.NoError(t, err)
		testastic.Equal(t, "info", cfg.LogConfig.Level)
		testastic.Equal(t, "text", cfg.LogConfig.Format)
		testastic.False(t, cfg.LogConfig.AddSource)
		testastic.Equal(t, 8080, cfg.Server.Port)
		testastic.Equal(t, "10s", cfg.Server.ReadTimeout.String())
		testastic.Equal(t, "10s", cfg.Server.WriteTimeout.String())
		testastic.Equal(t, "2m0s", cfg.Server.IdleTimeout.String())
		testastic.Equal(t, "20s", cfg.Server.ShutdownTimeout.String())
		testastic.Equal(t, "DATABASE_URL", cfg.Database.URLEnv)
		testastic.False(t, cfg.OTel.Enabled)
		testastic.Equal(t, "localhost:4317", cfg.OTel.Endpoint)
	})

	t.Run("decodes yaml config values", func(t *testing.T) {
		t.Parallel()

		// given: a temporary YAML config file
		configFile, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		testastic.NoError(t, err)

		_, err = configFile.WriteString(`log_config:
  level: "debug"
  format: "json"
  add_source: true

server:
  port: 9000
  read_timeout: "5s"
  write_timeout: "6s"
  idle_timeout: "7s"
  shutdown_timeout: "8s"

database:
  url_env: "TEST_DATABASE_URL"

pokeapi:
  base_url: "https://pokeapi.example/api/v2"
  timeout: "15s"
  concurrency: 7

otel:
  enabled: true
  endpoint: "otel.example:4317"
`)
		testastic.NoError(t, err)
		testastic.NoError(t, configFile.Close())

		// when: loading the YAML config file
		cfg, err := config.Load(configFile.Name())

		// then: it decodes the typed config values
		testastic.NoError(t, err)
		testastic.Equal(t, "debug", cfg.LogConfig.Level)
		testastic.Equal(t, "json", cfg.LogConfig.Format)
		testastic.True(t, cfg.LogConfig.AddSource)
		testastic.Equal(t, 9000, cfg.Server.Port)
		testastic.Equal(t, "5s", cfg.Server.ReadTimeout.String())
		testastic.Equal(t, "6s", cfg.Server.WriteTimeout.String())
		testastic.Equal(t, "7s", cfg.Server.IdleTimeout.String())
		testastic.Equal(t, "8s", cfg.Server.ShutdownTimeout.String())
		testastic.Equal(t, "TEST_DATABASE_URL", cfg.Database.URLEnv)
		testastic.Equal(t, "https://pokeapi.example/api/v2", cfg.PokeAPI.BaseURL)
		testastic.Equal(t, 7, cfg.PokeAPI.Concurrency)
		testastic.Equal(t, "15s", cfg.PokeAPI.Timeout.String())
		testastic.True(t, cfg.OTel.Enabled)
		testastic.Equal(t, "otel.example:4317", cfg.OTel.Endpoint)
	})

	t.Run("returns validation errors for missing required fields", func(t *testing.T) {
		t.Parallel()

		configFile, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
		testastic.NoError(t, err)

		_, err = configFile.WriteString(`log_config:
  level: ""
  format: "text"
  add_source: false

server:
  port: 0
  read_timeout: "0s"
  write_timeout: "0s"
  idle_timeout: "0s"
  shutdown_timeout: "0s"

database:
  url_env: ""

pokeapi:
  base_url: ""
  timeout: "0s"
  concurrency: 0

otel:
  enabled: true
  endpoint: ""
`)
		testastic.NoError(t, err)
		testastic.NoError(t, configFile.Close())

		_, err = config.Load(configFile.Name())

		testastic.NotNil(t, err)
		testastic.Contains(t, err.Error(), "log_config.level")
		testastic.Contains(t, err.Error(), "server.port")
		testastic.Contains(t, err.Error(), "server.read_timeout")
		testastic.Contains(t, err.Error(), "server.write_timeout")
		testastic.Contains(t, err.Error(), "server.idle_timeout")
		testastic.Contains(t, err.Error(), "server.shutdown_timeout")
		testastic.Contains(t, err.Error(), "database.url_env")
		testastic.Contains(t, err.Error(), "pokeapi.base_url")
		testastic.Contains(t, err.Error(), "pokeapi.timeout")
		testastic.Contains(t, err.Error(), "pokeapi.concurrency")
		testastic.Contains(t, err.Error(), "otel.endpoint")
	})

	t.Run("returns error when config file does not exist", func(t *testing.T) {
		t.Parallel()

		// when: loading a non-existent file
		_, err := config.Load("nonexistent.yaml")

		// then: it returns an error
		testastic.NotNil(t, err)
	})
}

func TestDatabaseURL(t *testing.T) {
	t.Run("reads the database url from the configured environment variable", func(t *testing.T) {
		t.Setenv("TEST_DATABASE_URL", "postgres://localhost:5432/app")

		cfg := config.Config{}
		cfg.Database.URLEnv = "TEST_DATABASE_URL"

		url, err := cfg.Database.URL()

		testastic.NoError(t, err)
		testastic.Equal(t, "postgres://localhost:5432/app", url)
	})

	t.Run("returns an error when the configured database url environment variable is empty", func(t *testing.T) {
		cfg := config.Config{}
		cfg.Database.URLEnv = "MISSING_DATABASE_URL"

		_, err := cfg.Database.URL()

		testastic.NotNil(t, err)
		testastic.Contains(t, err.Error(), "database env var is empty")
		testastic.Contains(t, err.Error(), "MISSING_DATABASE_URL")
	})
}
