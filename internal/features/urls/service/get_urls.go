package urls_service

import (
	"context"
	"fmt"
	"url-shortener/internal/core/domain"
)

func (s *URLsService) GetURLs(ctx context.Context) ([]domain.URL, error) {
	result, err := s.repository.GetURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get_urls: %w", err)
	}
	return result, nil
}
