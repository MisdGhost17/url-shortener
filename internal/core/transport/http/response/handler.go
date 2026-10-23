package core_http_response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	core_errors "url-shortener/internal/core/errors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write JSON: %v", err)
	}
}
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}
func ServiceError(w http.ResponseWriter, err error, fallback string) {
	var validation core_errors.ValidationError
	switch {
	case errors.As(err, &validation):
		Error(w, http.StatusBadRequest, validation.Message)
	case errors.Is(err, core_errors.ErrNotFound):
		Error(w, http.StatusNotFound, "ссылка не найдена")
	case errors.Is(err, core_errors.ErrConflict):
		Error(w, http.StatusConflict, "конфликт ссылки")
	default:
		log.Printf("%s: %v", fallback, err)
		Error(w, http.StatusInternalServerError, fallback)
	}
}
