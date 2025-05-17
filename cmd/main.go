package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"reference-service-go/internal/incoming/order"
)

func main() {
	// Create a new router
	r := chi.NewRouter()

	// Add some middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create a new order server
	orderServer := order.NewServer()

	// Register the order API handlers
	r.Mount("/v1", order.Handler(orderServer))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
