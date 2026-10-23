package urls_transport_http

import (
	"net/http"
	"strings"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

func (h *URLsHTTPHandler) Stats(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/stats/")
	if code == "" {
		core_http_response.Error(w, http.StatusBadRequest, "укажите короткий код")
		return
	}
	entry, err := h.service.GetURL(r.Context(), code)
	if err != nil {
		core_http_response.ServiceError(w, err, "ошибка получения ссылки")
		return
	}
	core_http_response.JSON(w, http.StatusOK, statsDTO(entry))
}
