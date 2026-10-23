package urls_memory_repository

import (
	"context"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
)

func (r *URLsRepository) GetURL(ctx context.Context, code string) (domain.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return domain.URL{}, err
	}
	entry, found := r.urls[code]
	if !found {
		return domain.URL{}, core_errors.ErrNotFound
	}
	return entry, nil
}
