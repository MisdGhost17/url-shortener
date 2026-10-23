package urls_transport_http

import (
	"net/http"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

func (h *URLsHTTPHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.GetURLs(r.Context())
	if err != nil {
		core_http_response.ServiceError(w, err, "ошибка получения ссылок")
		return
	}
	result := make([]StatsResponse, 0, len(entries))
	for _, entry := range entries {
		result = append(result, statsDTO(entry))
	}
	core_http_response.JSON(w, http.StatusOK, result)
}
