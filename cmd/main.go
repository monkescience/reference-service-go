package main

import (
	"log"
	"net/http"
	"os"
	statusapi "reference-service-go/internal/adapters/http/status"
	ui "reference-service-go/internal/adapters/http/ui"
	"reference-service-go/internal/middleware"

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

	// Create and register health API
	version := os.Getenv("VERSION")
	if version == "" {
		log.Fatal("VERSION environment variable is not set")
	}
	statusAPI := statusapi.NewAPI(version)
	// Mount the Status API under /api (e.g., /api/status/live)
	router.Route("/api", func(r chi.Router) {
		r.Mount("/", statusapi.Handler(statusAPI))
	})

	// Mount UI (tiles will call /api/status/live directly)
	dashboard, err := ui.NewDashboard()
	if err != nil {
		log.Fatalf("failed to init dashboard: %v", err)
	}
	dashboard.Routes(router)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
