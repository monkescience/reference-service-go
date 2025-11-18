package statusapi

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// API implements the generated ServerInterface for the Status API
// and holds runtime process information like version and start time.
type API struct {
	version   string
	startTime time.Time
}

func NewAPI(version string) *API {
	return &API{version: version, startTime: time.Now()}
}

// GetStatusLive handles GET /status/live
func (api *API) GetStatusLive(w http.ResponseWriter, r *http.Request) {
	resp := LiveResponse{
		Status:    "UP",
		Version:   api.version,
		Timestamp: time.Now(),
		Uptime:    time.Since(api.startTime).String(),
		GoVersion: runtime.Version(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetStatusReady handles GET /status/ready
func (api *API) GetStatusReady(w http.ResponseWriter, r *http.Request) {
	resp := ReadyResponse{
		Status:    "UP",
		Timestamp: time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetStatusMetrics handles GET /status/metrics
func (api *API) GetStatusMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
