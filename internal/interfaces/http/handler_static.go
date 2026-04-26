package http

import (
	"net/http"
	"path/filepath"
)

// StaticHandler serves static files (CSS, JS, React build)
type StaticHandler struct {
	staticDir string
}

// NewStaticHandler creates a new static file handler
func NewStaticHandler(staticDir string) *StaticHandler {
	return &StaticHandler{
		staticDir: staticDir,
	}
}

// ServeHTTP serves static files
func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(h.staticDir, r.URL.Path)
	http.ServeFile(w, r, filePath)
}
