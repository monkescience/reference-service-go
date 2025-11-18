package frontend

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tilesStr := r.URL.Query().Get("tiles")
	tiles, err := strconv.Atoi(tilesStr)
	if err != nil || tiles <= 0 {
		tiles = 3
	}

	data := struct {
		Tiles []struct{}
	}{
		Tiles: make([]struct{}, tiles),
	}

	templateName := "internal/adapters/http/frontend/templates/index.gohtml"
	if r.Header.Get("HX-Request") == "true" {
		templateName = "internal/adapters/http/frontend/templates/tiles.gohtml"
	}

	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		return
	}
}

func VersionTileHandler(w http.ResponseWriter, r *http.Request) {
	// In a real-world scenario, you might get this from a config or service discovery
	statusApiUrl := "http://localhost:8080/api/status/live"

	resp, err := http.Get(statusApiUrl)
	if err != nil {
		http.Error(w, "Failed to fetch version", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch version", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read version response", http.StatusInternalServerError)
		return
	}

	var statusResponse struct {
		Version string `json:"version"`
	}

	if err := json.Unmarshal(body, &statusResponse); err != nil {
		http.Error(w, "Failed to parse version response", http.StatusInternalServerError)
		return
	}

	hash := sha1.New()
	hash.Write([]byte(statusResponse.Version))
	hashBytes := hash.Sum(nil)
	hue, _ := strconv.ParseInt(hex.EncodeToString(hashBytes[:2]), 16, 64)
	hue = hue % 360

	data := struct {
		Version string
		Hue     int64
	}{
		Version: statusResponse.Version,
		Hue:     hue,
	}

	tmpl, err := template.ParseFiles("internal/adapters/http/frontend/templates/version-tile.gohtml")
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		return
	}
}
