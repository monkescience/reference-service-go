package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"reference-service-go/internal/incoming/http/instance"
	"strconv"
	"time"
)

type FrontendHandler struct {
	templates      *template.Template
	instanceClient *http.Client
	instanceURL    string
}

type InstanceTileData struct {
	Index int
	Info  instanceapi.InstanceInfoResponse
}

type TilesData struct {
	Instances []InstanceTileData
}

type IndexData struct {
	Count int
}

func NewFrontendHandler(templatesPath, instanceURL string) (*FrontendHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templatesPath, "*.gohtml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &FrontendHandler{
		templates: tmpl,
		instanceClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		instanceURL: instanceURL,
	}, nil
}

func (h *FrontendHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Count: 3, // Default number of tiles
	}

	if err := h.templates.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *FrontendHandler) TilesHandler(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count := 3 // Default

	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 && parsedCount <= 20 {
			count = parsedCount
		}
	}

	// Fetch instance info for each tile
	instances := make([]InstanceTileData, count)
	for i := 0; i < count; i++ {
		info, err := h.fetchInstanceInfo()
		if err != nil {
			// Use error data if fetch fails
			info = instanceapi.InstanceInfoResponse{
				Version:   "error",
				Hostname:  "failed to fetch",
				Uptime:    "N/A",
				GoVersion: "N/A",
				Timestamp: time.Now(),
			}
		}
		instances[i] = InstanceTileData{
			Index: i + 1,
			Info:  info,
		}
	}

	data := TilesData{
		Instances: instances,
	}

	if err := h.templates.ExecuteTemplate(w, "tiles.gohtml", data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render tiles: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *FrontendHandler) fetchInstanceInfo() (instanceapi.InstanceInfoResponse, error) {
	resp, err := h.instanceClient.Get(h.instanceURL)
	if err != nil {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf("failed to fetch instance info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info instanceapi.InstanceInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return info, nil
}
