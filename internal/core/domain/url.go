package domain

import (
	"strings"
	"time"
	core_errors "url-shortener/internal/core/errors"
)

// URL is independent of HTTP serialization and storage implementations.
type URL struct {
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	Clicks      int
}

const CodeLength = 6
const CodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewURL(code, original string, createdAt time.Time) URL {
	return URL{ShortCode: code, OriginalURL: strings.TrimSpace(original), CreatedAt: createdAt}
}

// ValidateOriginalURL preserves the existing URL acceptance rules.
func ValidateOriginalURL(original string) error {
	if original == "" {
		return core_errors.ValidationError{Message: "поле url обязательно"}
	}
	if !strings.HasPrefix(original, "http://") && !strings.HasPrefix(original, "https://") {
		return core_errors.ValidationError{Message: "URL должен начинаться с http:// или https://"}
	}
	return nil
}

func (u URL) Validate() error {
	if err := ValidateOriginalURL(u.OriginalURL); err != nil {
		return err
	}
	if len(u.ShortCode) != CodeLength {
		return core_errors.ValidationError{Message: "неверный короткий код"}
	}
	for _, c := range u.ShortCode {
		if !strings.ContainsRune(CodeAlphabet, c) {
			return core_errors.ValidationError{Message: "неверный короткий код"}
		}
	}
	if u.CreatedAt.IsZero() || u.Clicks < 0 {
		return core_errors.ValidationError{Message: "неверное состояние ссылки"}
	}
	return nil
}
