package main

import (
	"log"
	"net/http"
	"os"
	"reference-service-go/internal/middleware"
	"reference-service-go/internal/status"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"reference-service-go/internal/incoming/order"
)

func main() {
	// Create a new router
	router := chi.NewRouter()

	responseTimeHistogramMetric := middleware.NewHttpResponseTimeHistogramMetric()

	// Add some middleware
	router.Use(responseTimeHistogramMetric.ResponseTimes)
	router.Use(chimiddleware.Recoverer)

	// Create a new order server
	orderServer := order.NewServer()

	// Register the order API handlers
	router.Mount("/v1", order.Handler(orderServer))

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
