package web_transport_http

import (
	"net/http"
	"path/filepath"
	core_http_server "url-shortener/internal/core/transport/http/server"
)

type WebHTTPHandler struct{ directory string }

func NewWebHTTPHandler(directory string) *WebHTTPHandler {
	return &WebHTTPHandler{directory: directory}
}
func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	files := http.StripPrefix("/static/", http.FileServer(http.Dir(h.directory)))
	return []core_http_server.Route{
		{Method: http.MethodGet, Path: "/{$}", Handler: h.Index},
		{Path: "/static/", Handler: files.ServeHTTP},
		{Path: "/favicon.svg", Handler: h.Favicon},
	}
}
func (h *WebHTTPHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(h.directory, "index.html"))
}
func (h *WebHTTPHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	// Legacy behavior: this HTML page is served with status 200.
	http.ServeFile(w, r, filepath.Join(h.directory, "404.html"))
}
func (h *WebHTTPHandler) Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, filepath.Join(h.directory, "favicon.svg"))
}
