package urls_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
)

func (s *URLsService) CreateURL(ctx context.Context, original string) (domain.URL, error) {
	original = strings.TrimSpace(original)
	if err := domain.ValidateOriginalURL(original); err != nil {
		return domain.URL{}, fmt.Errorf("validate URL: %w", err)
	}
	for {
		if err := ctx.Err(); err != nil {
			return domain.URL{}, err
		}
		code, err := s.generator.Generate()
		if err != nil {
			return domain.URL{}, fmt.Errorf("generate code: %w", err)
		}
		entry := domain.NewURL(code, original, time.Now())
		if err := entry.Validate(); err != nil {
			return domain.URL{}, fmt.Errorf("validate generated URL: %w", err)
		}
		saved, err := s.repository.CreateURL(ctx, entry)
		if errors.Is(err, core_errors.ErrConflict) {
			continue
		}
		if err != nil {
			return domain.URL{}, fmt.Errorf("create URL: %w", err)
		}
		return saved, nil
	}
}
