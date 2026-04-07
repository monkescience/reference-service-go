package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	LogConfig struct {
		Level     string `yaml:"level"`      // Log level (debug, info, warn, error)
		Format    string `yaml:"format"`     // Log format (json, text)
		AddSource bool   `yaml:"add_source"` // Include source file and line number
	} `yaml:"log_config"`

	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	PokeAPI  PokeAPIConfig  `yaml:"pokeapi"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	URL string `yaml:"url"`
}

// PokeAPIConfig holds settings for the PokeAPI client.
type PokeAPIConfig struct {
	BaseURL     string        `yaml:"base_url"`
	Timeout     time.Duration `yaml:"timeout"`
	Concurrency int           `yaml:"concurrency"`
}

// Load reads configuration from the specified YAML file.
func Load(path string) (*Config, error) {
	//nolint:gosec // Config file path is expected to be provided by trusted deployment configuration
	configFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	//nolint:noinlineerr,wsl // Defer close pattern is idiomatic for resource cleanup
	defer func() {
		if closeErr := configFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close config file: %w", closeErr))
		}
	}()

	var cfg Config

	decoder := yaml.NewDecoder(configFile)

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &cfg, nil
}
