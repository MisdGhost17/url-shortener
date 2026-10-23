package urls_memory_repository

import (
	"context"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
)

func (r *URLsRepository) CreateURL(ctx context.Context, entry domain.URL) (domain.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return domain.URL{}, err
	}
	if _, found := r.urls[entry.ShortCode]; found {
		return domain.URL{}, core_errors.ErrConflict
	}
	r.urls[entry.ShortCode] = entry
	return entry, nil
}
