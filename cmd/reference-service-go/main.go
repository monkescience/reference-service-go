package main

import (
	"flag"
	"log"
	"log/slog"
	"reference-service-go/internal/build"
	"reference-service-go/internal/config"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"

	importsapi "reference-service-go/internal/incoming/http/imports"
)

const (
	serverPort         = 8080
	serverReadTimeout  = 10 * time.Second
	serverWriteTimeout = 10 * time.Second
	serverIdleTimeout  = 120 * time.Second
	shutdownTimeout    = 20 * time.Second
)

func setupLogger(cfg *config.Config) *slog.Logger {
	logConfig := vital.LogConfig{
		Level:     cfg.LogConfig.Level,
		Format:    cfg.LogConfig.Format,
		AddSource: cfg.LogConfig.AddSource,
	}

	handler, err := vital.NewHandlerFromConfig(logConfig, vital.WithBuiltinKeys())
	if err != nil {
		log.Fatalf("failed to create logger handler: %v", err)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func setupRouter(logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(vital.Recovery(logger))
	router.Use(vital.RequestLogger(logger))

	importHandler := importsapi.NewImportHandler(logger)
	importsapi.HandlerFromMux(importHandler, router)

	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(build.Version),
	)
	router.Mount("/health", healthHandler)

	return router
}

func main() {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := setupLogger(cfg)

	router := setupRouter(logger)

	server := vital.NewServer(
		router,
		vital.WithPort(serverPort),
		vital.WithReadTimeout(serverReadTimeout),
		vital.WithWriteTimeout(serverWriteTimeout),
		vital.WithIdleTimeout(serverIdleTimeout),
		vital.WithShutdownTimeout(shutdownTimeout),
		vital.WithLogger(logger),
	)

	err = server.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
