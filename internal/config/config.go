package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
)

var (
	errDatabaseEnvEmpty       = errors.New("database env var is empty")
	errLogLevelEmpty          = errors.New("log_config.level must not be empty")
	errLogFormatEmpty         = errors.New("log_config.format must not be empty")
	errServerPortZero         = errors.New("server.port must not be zero")
	errServerReadTimeoutZero  = errors.New("server.read_timeout must not be zero")
	errServerWriteTimeoutZero = errors.New("server.write_timeout must not be zero")
	errServerIdleTimeoutZero  = errors.New("server.idle_timeout must not be zero")
	errServerShutdownZero     = errors.New("server.shutdown_timeout must not be zero")
	errDatabaseURLEnvEmpty    = errors.New("database.url_env must not be empty")
	errPokeAPIBaseURLEmpty    = errors.New("pokeapi.base_url must not be empty")
	errPokeAPITimeoutZero     = errors.New("pokeapi.timeout must not be zero")
	errPokeAPIConcurrencyZero = errors.New("pokeapi.concurrency must not be zero")
	errOTelEndpointEmpty      = errors.New("otel.endpoint must not be empty when otel.enabled is true")
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
	OTel     OTelConfig     `yaml:"otel"`
}

// OTelConfig holds OpenTelemetry tracing settings.
type OTelConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DatabaseConfig holds PostgreSQL connection settings.
//
// The database URL itself is loaded from the configured environment variable so
// secrets do not live in the YAML config file.
type DatabaseConfig struct {
	URLEnv string `yaml:"url_env"`
}

// PokeAPIConfig holds settings for the PokeAPI client.
type PokeAPIConfig struct {
	BaseURL     string        `yaml:"base_url"`
	Timeout     time.Duration `yaml:"timeout"`
	Concurrency int           `yaml:"concurrency"`
}

// Load reads configuration from the specified YAML file.
func Load(path string) (_ *Config, err error) {
	//nolint:gosec // Config file path is expected to be provided by trusted deployment configuration
	configFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	//nolint:noinlineerr,wsl // Defer close pattern is idiomatic for resource cleanup
	defer func() {
		if closeErr := configFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close config file: %w", closeErr))
		}
	}()

	var cfg Config

	decoder := yaml.NewDecoder(configFile)

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// URL loads the database URL from the configured environment variable.
func (d DatabaseConfig) URL() (string, error) {
	value := strings.TrimSpace(os.Getenv(d.URLEnv))
	if value == "" {
		return "", fmt.Errorf("%w: %s", errDatabaseEnvEmpty, d.URLEnv)
	}

	return value, nil
}

// Validate checks that required configuration values are present.
func (c Config) Validate() error {
	var err error

	if strings.TrimSpace(c.LogConfig.Level) == "" {
		err = errors.Join(err, errLogLevelEmpty)
	}

	if strings.TrimSpace(c.LogConfig.Format) == "" {
		err = errors.Join(err, errLogFormatEmpty)
	}

	if c.Server.Port == 0 {
		err = errors.Join(err, errServerPortZero)
	}

	if c.Server.ReadTimeout == 0 {
		err = errors.Join(err, errServerReadTimeoutZero)
	}

	if c.Server.WriteTimeout == 0 {
		err = errors.Join(err, errServerWriteTimeoutZero)
	}

	if c.Server.IdleTimeout == 0 {
		err = errors.Join(err, errServerIdleTimeoutZero)
	}

	if c.Server.ShutdownTimeout == 0 {
		err = errors.Join(err, errServerShutdownZero)
	}

	if strings.TrimSpace(c.Database.URLEnv) == "" {
		err = errors.Join(err, errDatabaseURLEnvEmpty)
	}

	if strings.TrimSpace(c.PokeAPI.BaseURL) == "" {
		err = errors.Join(err, errPokeAPIBaseURLEmpty)
	}

	if c.PokeAPI.Timeout == 0 {
		err = errors.Join(err, errPokeAPITimeoutZero)
	}

	if c.PokeAPI.Concurrency == 0 {
		err = errors.Join(err, errPokeAPIConcurrencyZero)
	}

	if c.OTel.Enabled && strings.TrimSpace(c.OTel.Endpoint) == "" {
		err = errors.Join(err, errOTelEndpointEmpty)
	}

	return err
}
