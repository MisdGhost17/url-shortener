package urls_transport_http

import (
	"errors"
	"log"
	"net/http"
	"strings"
	core_errors "url-shortener/internal/core/errors"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

func (h *URLsHTTPHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	entry, err := h.service.ResolveURL(r.Context(), code)
	if errors.Is(err, core_errors.ErrNotFound) {
		h.notFound.ServeHTTP(w, r)
		return
	}
	if err != nil {
		core_http_response.ServiceError(w, err, "ошибка перехода")
		return
	}
	log.Printf("[REDIRECT] %s -> %s", code, entry.OriginalURL)
	http.Redirect(w, r, entry.OriginalURL, http.StatusFound)
}
