package urls_memory_repository

import (
	"context"
	"url-shortener/internal/core/domain"
)

func (r *URLsRepository) GetURLs(ctx context.Context) ([]domain.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries := make([]domain.URL, 0, len(r.urls))
	for _, entry := range r.urls {
		entries = append(entries, entry)
	}
	return entries, nil
}
