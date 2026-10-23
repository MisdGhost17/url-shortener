package core_http_server

import (
	"net/http"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// RegisterRoutes keeps the original JSON 405 responses. An empty Method
// intentionally preserves legacy method acceptance for list and statistics.
func RegisterRoutes(mux *http.ServeMux, routes ...Route) {
	for _, route := range routes {
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if route.Method != "" && r.Method != route.Method {
				core_http_response.Error(w, http.StatusMethodNotAllowed, "метод не поддерживается")
				return
			}
			route.Handler(w, r)
		})
	}
}
