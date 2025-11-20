package frontend

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	instanceapi "reference-service-go/internal/incoming/http/instance"
)

// FrontendHandler handles frontend HTTP requests for the web UI.
type FrontendHandler struct {
	templates      *template.Template
	instanceClient *http.Client
	instanceURL    string
	tileColors     []string
}

// InstanceTileData represents data for a single instance tile in the UI.
type InstanceTileData struct {
	Index int
	Info  instanceapi.InstanceInfoResponse
	Color string
}

// TilesData holds the collection of instance tiles to render.
type TilesData struct {
	Instances []InstanceTileData
}

// IndexData contains data for rendering the index page.
type IndexData struct {
	Count int
}

// NewFrontendHandler creates a new frontend handler with the specified templates path, instance API URL, and tile colors.
func NewFrontendHandler(
	templatesPath, instanceURL string,
	tileColors []string,
) (*FrontendHandler, error) {
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
		tileColors:  tileColors,
	}, nil
}

// IndexHandler serves the main index page with the default tile count.
func (h *FrontendHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Count: 3, // Default number of tiles
	}

	if err := h.templates.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to render template: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

// TilesHandler renders instance tiles based on the count query parameter.
func (h *FrontendHandler) TilesHandler(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count := 3 // Default

	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 &&
			parsedCount <= 20 {
			count = parsedCount
		}
	}

	instances := make([]InstanceTileData, count)
	for i := 0; i < count; i++ {
		info, err := h.fetchInstanceInfo()
		if err != nil {
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
			Color: h.getColorForVersion(info.Version),
		}
	}

	data := TilesData{
		Instances: instances,
	}

	if err := h.templates.ExecuteTemplate(w, "tiles.gohtml", data); err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to render tiles: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *FrontendHandler) fetchInstanceInfo() (instanceapi.InstanceInfoResponse, error) {
	resp, err := h.instanceClient.Get(h.instanceURL)
	if err != nil {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf(
			"failed to fetch instance info: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf(
			"unexpected status code: %d",
			resp.StatusCode,
		)
	}

	var info instanceapi.InstanceInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return instanceapi.InstanceInfoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return info, nil
}

// getColorForVersion returns a color from the configured tile colors based on the version string.
// Uses a hash function to deterministically select a color.
func (h *FrontendHandler) getColorForVersion(version string) string {
	if len(h.tileColors) == 0 {
		return "#667eea"
	}

	hasher := fnv.New32a()
	hasher.Write([]byte(version))
	index := hasher.Sum32() % uint32(len(h.tileColors))

	return h.tileColors[index]
}
