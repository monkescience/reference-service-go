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
	"reference-service-go/internal/core/catch"
	"reference-service-go/internal/core/pokemon"
	"reference-service-go/internal/incoming/referencehttp"
	"reference-service-go/internal/outgoing/pokeapi"
	"reference-service-go/internal/outgoing/referencepg"
	"reference-service-go/internal/outgoing/tracing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	defaultServerPort  = 8080
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
	migrateOnly := flag.Bool("migrate-only", false, "Run database migrations and exit")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if *migrateOnly {
		err = referencepg.Migrate(context.Background(), cfg.Database.URL)
		if err != nil {
			return fmt.Errorf("running migrations: %w", err)
		}

		return nil
	}

	return runServer(cfg)
}

func runServer(cfg *config.Config) error {
	logger := setupLogger(cfg)

	ctx := context.Background()

	err := tracing.Setup(ctx, cfg.OTel.Enabled, cfg.OTel.Endpoint)
	if err != nil {
		return fmt.Errorf("setting up tracing: %w", err)
	}

	store, err := referencepg.New(ctx, cfg.Database.URL)
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

	port := cfg.Server.Port
	if port == 0 {
		port = defaultServerPort
	}

	server := vital.NewServer(
		router,
		vital.WithPort(port),
		vital.WithReadTimeout(serverReadTimeout),
		vital.WithWriteTimeout(serverWriteTimeout),
		vital.WithIdleTimeout(serverIdleTimeout),
		vital.WithShutdownTimeout(shutdownTimeout),
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
		vital.WithVersion(build.Version),
	)
	router.Mount("/health", healthHandler)

	return router
}
