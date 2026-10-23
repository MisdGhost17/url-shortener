package urls_memory_repository

import (
	"context"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
)

func (r *URLsRepository) ResolveURL(ctx context.Context, code string) (domain.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return domain.URL{}, err
	}
	entry, found := r.urls[code]
	if !found {
		return domain.URL{}, core_errors.ErrNotFound
	}
	entry.Clicks++
	r.urls[code] = entry
	return entry, nil
}
