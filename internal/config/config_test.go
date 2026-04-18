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

database:
  url: "postgres://localhost:5432/app"

pokeapi:
  base_url: "https://pokeapi.example/api/v2"
  timeout: "15s"
  concurrency: 7
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
		testastic.Equal(t, "postgres://localhost:5432/app", cfg.Database.URL)
		testastic.Equal(t, "https://pokeapi.example/api/v2", cfg.PokeAPI.BaseURL)
		testastic.Equal(t, 7, cfg.PokeAPI.Concurrency)
		testastic.Equal(t, "15s", cfg.PokeAPI.Timeout.String())
	})

	t.Run("returns error when config file does not exist", func(t *testing.T) {
		t.Parallel()

		// when: loading a non-existent file
		_, err := config.Load("nonexistent.yaml")

		// then: it returns an error
		testastic.NotNil(t, err)
	})
}
