package urls_service

import (
	"context"
	"fmt"
	"url-shortener/internal/core/domain"
)

func (s *URLsService) ResolveURL(ctx context.Context, code string) (domain.URL, error) {
	result, err := s.repository.ResolveURL(ctx, code)
	if err != nil {
		return domain.URL{}, fmt.Errorf("resolve_url: %w", err)
	}
	return result, nil
}
