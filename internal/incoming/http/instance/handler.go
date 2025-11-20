package instanceapi

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

// InstanceHandler handles instance information requests.
type InstanceHandler struct {
	version string
}

// NewInstanceHandler creates a new instance handler with the specified version.
func NewInstanceHandler(version string) *InstanceHandler {
	return &InstanceHandler{
		version: version,
	}
}

// GetInstanceInfo returns information about the running instance including version, hostname, uptime, and Go version.
func (h *InstanceHandler) GetInstanceInfo(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	uptime := time.Since(startTime)

	response := InstanceInfoResponse{
		Version:   h.version,
		Hostname:  hostname,
		Uptime:    uptime.String(),
		GoVersion: runtime.Version(),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
