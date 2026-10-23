package core_http_response

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	core_errors "url-shortener/internal/core/errors"
)

func TestServiceErrors(t *testing.T) {
	cases := []struct {
		err     error
		status  int
		message string
	}{
		{fmt.Errorf("wrapped: %w", core_errors.ValidationError{Message: "bad input"}), 400, "bad input"},
		{fmt.Errorf("wrapped: %w", core_errors.ErrNotFound), 404, "ссылка не найдена"},
		{core_errors.ErrConflict, 409, "конфликт ссылки"},
		{errors.New("private storage detail"), 500, "internal error"},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		ServiceError(w, tc.err, "internal error")
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.message) || strings.Contains(w.Body.String(), "private storage detail") {
			t.Fatalf("unexpected error response: %d %s", w.Code, w.Body.String())
		}
	}
}
