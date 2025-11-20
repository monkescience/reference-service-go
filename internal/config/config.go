package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	Version    string   `yaml:"-"`           // Version must be set via VERSION environment variable only
	TileColors []string `yaml:"tile-colors"` // Available colors for instance tiles based on version
}

// Load reads configuration from the specified YAML file and environment variables.
// The VERSION environment variable is required and must be set; it cannot be configured via the config file.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	cfg.Version = os.Getenv("VERSION")
	if cfg.Version == "" {
		return nil, fmt.Errorf("VERSION environment variable is required")
	}

	if len(cfg.TileColors) == 0 {
		return nil, fmt.Errorf("tile-colors must be configured in the config file")
	}

	return &cfg, nil
}
