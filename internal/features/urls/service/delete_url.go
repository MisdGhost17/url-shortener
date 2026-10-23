package urls_service

import (
	"context"
	"fmt"
)

func (s *URLsService) DeleteURL(ctx context.Context, code string) error {
	if err := s.repository.DeleteURL(ctx, code); err != nil {
		return fmt.Errorf("delete URL: %w", err)
	}
	return nil
}
