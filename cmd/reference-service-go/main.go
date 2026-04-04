package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"reference-service-go/internal/build"
	"reference-service-go/internal/config"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/http/pokeapi"
	"reference-service-go/internal/outgoing/postgres"
	"reference-service-go/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"

	importsapi "reference-service-go/internal/incoming/http/imports"
	pokemonapi "reference-service-go/internal/incoming/http/pokemon"
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

func main() {
	err := run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func run() error {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger := setupLogger(cfg)

	ctx := context.Background()

	pool, err := postgres.Connect(ctx, cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	defer pool.Close()

	queries := postgres.New(pool)

	pokeapiClient, err := pokeapi.NewFetcher(
		&http.Client{Timeout: cfg.PokeAPI.Timeout},
		cfg.PokeAPI.BaseURL,
	)
	if err != nil {
		return fmt.Errorf("creating pokeapi client: %w", err)
	}

	importService := service.NewImportService(logger, pokeapiClient, queries, pool, cfg.PokeAPI.Concurrency)

	defer importService.Shutdown()

	gachaService := service.NewGachaService(logger, queries, domain.DefaultRand{})
	router := setupRouter(logger, importService, gachaService, queries)

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
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}

func setupRouter(
	logger *slog.Logger,
	importService *service.ImportService,
	gachaService *service.GachaService,
	queries *postgres.Queries,
) chi.Router {
	router := chi.NewRouter()
	router.Use(vital.Recovery(logger))
	router.Use(vital.RequestLogger(logger))

	importHandler := importsapi.NewImportHandler(logger, importService)
	importsapi.HandlerFromMux(importHandler, router)

	pokemonHandler := pokemonapi.NewPokemonHandler(logger, gachaService, queries)
	pokemonapi.HandlerFromMux(pokemonHandler, router)

	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(build.Version),
	)
	router.Mount("/health", healthHandler)

	return router
}
