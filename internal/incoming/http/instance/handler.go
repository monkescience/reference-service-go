package instanceapi

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type InstanceHandler struct {
	version string
}

func NewInstanceHandler(version string) *InstanceHandler {
	return &InstanceHandler{
		version: version,
	}
}

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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
