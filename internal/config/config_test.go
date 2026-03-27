package config_test

import (
	"reference-service-go/internal/config"
	"testing"

	"github.com/monkescience/testastic"
)

func TestLoad(t *testing.T) {
	t.Run("loads valid config", func(t *testing.T) {
		// GIVEN
		t.Setenv("VERSION", "1.0.0")

		// WHEN
		cfg, err := config.Load("../../config/config.yaml")

		// THEN
		testastic.NoError(t, err)
		testastic.Equal(t, "1.0.0", cfg.Version)
		testastic.Equal(t, "info", cfg.LogConfig.Level)
		testastic.Equal(t, "text", cfg.LogConfig.Format)
		testastic.False(t, cfg.LogConfig.AddSource)
	})

	t.Run("returns error when VERSION is missing", func(t *testing.T) {
		// GIVEN
		t.Setenv("VERSION", "")

		// WHEN
		_, err := config.Load("../../config/config.yaml")

		// THEN
		testastic.Equal(t, config.ErrVersionRequired, err)
	})

	t.Run("returns error when config file does not exist", func(t *testing.T) {
		// WHEN
		_, err := config.Load("nonexistent.yaml")

		// THEN
		testastic.NotNil(t, err)
	})
}
