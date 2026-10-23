package urls_transport_http

import (
	"log"
	"net/http"
	"strings"
	core_http_response "url-shortener/internal/core/transport/http/response"
)

func (h *URLsHTTPHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/delete/")
	if err := h.service.DeleteURL(r.Context(), code); err != nil {
		core_http_response.ServiceError(w, err, "ошибка удаления ссылки")
		return
	}
	log.Printf("[DELETE] %s", code)
	core_http_response.JSON(w, http.StatusOK, map[string]string{"message": "удалено"})
}
