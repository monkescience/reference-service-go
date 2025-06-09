package status

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// API represents the status API
type API struct {
	version   string
	startTime time.Time
}

// LiveResponse represents the liveness check response
type LiveResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	GoVersion string    `json:"go_version"`
}

// ReadyResponse represents the readiness check response
type ReadyResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// New creates a new status API
func New(version string) *API {
	return &API{
		version:   version,
		startTime: time.Now(),
	}
}

// RegisterRoutes registers the status API routes
func (api *API) RegisterRoutes(router chi.Router) {
	router.Group(
		func(router chi.Router) {
			router.Get("/status/live", api.LiveHandler)
			router.Get("/status/ready", api.ReadyHandler)
			router.Method("GET", "/status/metrics", api.MetricsHandler())
		},
	)
}

// LiveHandler handles liveness check requests
func (api *API) LiveHandler(w http.ResponseWriter, r *http.Request) {
	response := LiveResponse{
		Status:    "UP",
		Version:   api.version,
		Timestamp: time.Now(),
		Uptime:    time.Since(api.startTime).String(),
		GoVersion: runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadyHandler handles readiness check requests
func (api *API) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	response := ReadyResponse{
		Status:    "UP",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// MetricsHandler returns the Prometheus metrics handler
func (api *API) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
