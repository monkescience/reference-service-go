package main

import (
	"log"
	"net/http"
	"os"
	"reference-service-go/internal/adapters/http/order"
	repomem "reference-service-go/internal/adapters/repository/memory"
	"reference-service-go/internal/middleware"
	"reference-service-go/internal/status"
	uc "reference-service-go/internal/usecase/order"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Create a new router
	router := chi.NewRouter()

	responseTimeHistogramMetric := middleware.NewHttpResponseTimeHistogramMetric()

	// Add some middleware
	router.Use(responseTimeHistogramMetric.ResponseTimes)
	router.Use(chimiddleware.Recoverer)

	// Wire onion layers
	repo := repomem.NewRepository()
	service := uc.NewService(repo)
	h := order.NewAPI(service)

	// Register the order API handlers
	router.Mount("/v1", order.Handler(h))

	// Create and register health API
	version := os.Getenv("VERSION")
	if version == "" {
		log.Fatal("VERSION environment variable is not set")
	}
	healthAPI := status.New(version)
	healthAPI.RegisterRoutes(router)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
