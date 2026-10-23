package urls_transport_http

import (
	"encoding/json"
	"log"
	"net/http"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

func (h *URLsHTTPHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		core_http_response.Error(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	entry, err := h.service.CreateURL(r.Context(), request.URL)
	if err != nil {
		core_http_response.ServiceError(w, err, "ошибка генерации")
		return
	}
	log.Printf("[CREATE] %s -> %s", entry.ShortCode, entry.OriginalURL)
	// Preserve the deployed API's HTTPS short URLs, including behind a proxy.
	core_http_response.JSON(w, http.StatusCreated, ShortenResponse{
		ShortURL: "https://" + r.Host + "/" + entry.ShortCode, ShortCode: entry.ShortCode, Original: entry.OriginalURL,
	})
}
