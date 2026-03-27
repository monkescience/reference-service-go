package config_test

import (
	"reference-service-go/internal/config"
	"testing"

	"github.com/monkescience/testastic"
)

func TestLoad(t *testing.T) {
	t.Run("loads valid config", func(t *testing.T) {
		// when: loading a valid config file
		cfg, err := config.Load("../../config/config.yaml")

		// then: it returns the parsed config
		testastic.NoError(t, err)
		testastic.Equal(t, "info", cfg.LogConfig.Level)
		testastic.Equal(t, "text", cfg.LogConfig.Format)
		testastic.False(t, cfg.LogConfig.AddSource)
	})

	t.Run("returns error when config file does not exist", func(t *testing.T) {
		// when: loading a non-existent file
		_, err := config.Load("nonexistent.yaml")

		// then: it returns an error
		testastic.NotNil(t, err)
	})
}
