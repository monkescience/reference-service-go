package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"reference-service-go/internal/config"
	"reference-service-go/internal/incoming/http/frontend"
	instanceapi "reference-service-go/internal/incoming/http/instance"
	"reference-service-go/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := chi.NewRouter()

	responseTimeHistogramMetric := middleware.NewHttpResponseTimeHistogramMetric()

	// Add some middleware
	router.Use(responseTimeHistogramMetric.ResponseTimes)
	router.Use(chimiddleware.Recoverer)

	// Instance API handler
	instanceHandler := instanceapi.NewInstanceHandler(cfg.Version)
	instanceapi.HandlerFromMux(instanceHandler, router)

	// Frontend handler
	templatesPath := filepath.Join("internal", "incoming", "http", "frontend", "templates")
	frontendHandler, err := frontend.NewFrontendHandler(templatesPath, "http://localhost:8080/instance/info")
	if err != nil {
		log.Fatalf("failed to create frontend handler: %v", err)
	}

	router.Get("/", frontendHandler.IndexHandler)
	router.Get("/tiles", frontendHandler.TilesHandler)

	log.Println("Starting server on :8080")
	log.Println("Open http://localhost:8080 in your browser")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
