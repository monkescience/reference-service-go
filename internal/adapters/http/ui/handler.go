package ui

import (
	"html/template"
	"net/http"
)

// Dashboard is a tiny UI adapter that renders an HTMX dashboard
// with tiles that display the running version; tiles call /status/live directly.
type Dashboard struct {
	tmpl *template.Template
}

func NewDashboard() (*Dashboard, error) {
	// Parse all templates under templates/
	t, err := template.ParseFS(templatesFS, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	return &Dashboard{tmpl: t}, nil
}

// Routes mounts the UI routes under the given router.
func (d *Dashboard) Routes(r chiRouter) {
	// GET /ui -> full page
	r.Get("/ui", http.HandlerFunc(d.index))
}

// Minimal local interface to avoid importing chi directly here.
type chiRouter interface {
	Get(pattern string, h http.HandlerFunc)
}

func (d *Dashboard) index(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = d.tmpl.ExecuteTemplate(w, "index.gohtml", data)
}
