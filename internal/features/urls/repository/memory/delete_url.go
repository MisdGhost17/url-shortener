package urls_memory_repository

import (
	"context"
	core_errors "url-shortener/internal/core/errors"
)

func (r *URLsRepository) DeleteURL(ctx context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, found := r.urls[code]; !found {
		return core_errors.ErrNotFound
	}
	delete(r.urls, code)
	return nil
}
