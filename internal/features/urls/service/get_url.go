package urls_service

import (
	"context"
	"fmt"
	"url-shortener/internal/core/domain"
)

func (s *URLsService) GetURL(ctx context.Context, code string) (domain.URL, error) {
	result, err := s.repository.GetURL(ctx, code)
	if err != nil {
		return domain.URL{}, fmt.Errorf("get_url: %w", err)
	}
	return result, nil
}
