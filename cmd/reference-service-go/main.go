package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reference-service-go/internal/build"
	"reference-service-go/internal/config"
	"reference-service-go/internal/core/catch"
	"reference-service-go/internal/core/pokemon"
	"reference-service-go/internal/incoming/referencehttp"
	"reference-service-go/internal/outgoing/pokeapi"
	"reference-service-go/internal/outgoing/referencepg"
	"reference-service-go/internal/outgoing/tracing"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func setupLogger(cfg *config.Config) (*slog.Logger, error) {
	logConfig := vital.LogConfig{
		Level:     cfg.LogConfig.Level,
		Format:    cfg.LogConfig.Format,
		AddSource: cfg.LogConfig.AddSource,
	}

	handler, err := vital.NewHandlerFromConfig(logConfig, vital.WithBuiltinKeys())
	if err != nil {
		return nil, fmt.Errorf("create logger handler: %w", err)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}

func main() {
	err := run()
	if err != nil {
		slog.ErrorContext(context.Background(), "server error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")
	migrateOnly := flag.Bool("migrate-only", false, "Run database migrations and exit")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger, err := setupLogger(cfg)
	if err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if *migrateOnly {
		databaseURL, databaseErr := cfg.Database.URL()
		if databaseErr != nil {
			return fmt.Errorf("loading database url: %w", databaseErr)
		}

		err = referencepg.Migrate(context.Background(), databaseURL)
		if err != nil {
			return fmt.Errorf("running migrations: %w", err)
		}

		logger.Info("migrations completed")

		return nil
	}

	return runServer(cfg, logger)
}

func runServer(cfg *config.Config, logger *slog.Logger) error {
	ctx := context.Background()

	err := tracing.Setup(ctx, cfg.OTel.Enabled, cfg.OTel.Endpoint)
	if err != nil {
		return fmt.Errorf("setting up tracing: %w", err)
	}

	databaseURL, err := cfg.Database.URL()
	if err != nil {
		return fmt.Errorf("loading database url: %w", err)
	}

	store, err := referencepg.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	defer store.Close()

	pokeapiClient, err := pokeapi.NewFetcher(
		&http.Client{
			Timeout:   cfg.PokeAPI.Timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		cfg.PokeAPI.BaseURL,
	)
	if err != nil {
		return fmt.Errorf("creating pokeapi client: %w", err)
	}

	pokemonService := pokemon.NewService(pokeapiClient, store, store, cfg.PokeAPI.Concurrency)

	defer pokemonService.Shutdown()

	catchService := catch.NewService(store, store, catch.DefaultRand{})
	router := setupRouter(logger, pokemonService, catchService)

	server := vital.NewServer(
		router,
		vital.WithPort(cfg.Server.Port),
		vital.WithReadTimeout(cfg.Server.ReadTimeout),
		vital.WithWriteTimeout(cfg.Server.WriteTimeout),
		vital.WithIdleTimeout(cfg.Server.IdleTimeout),
		vital.WithShutdownTimeout(cfg.Server.ShutdownTimeout),
		vital.WithLogger(logger),
		vital.WithShutdownFunc(tracing.Shutdown),
	)

	err = server.Run()
	if err != nil {
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}

func setupRouter(
	logger *slog.Logger,
	pokemonService *pokemon.Service,
	catchService *catch.Service,
) chi.Router {
	router := chi.NewRouter()
	router.Use(vital.Recovery(logger))
	router.Use(otelhttp.NewMiddleware(build.ServiceName))
	router.Use(vital.RequestLogger(logger))

	handler := referencehttp.NewHandler(pokemonService, catchService)
	referencehttp.HandlerFromMux(handler, router)

	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(build.Version()),
	)
	router.Mount("/health", healthHandler)

	return router
}
