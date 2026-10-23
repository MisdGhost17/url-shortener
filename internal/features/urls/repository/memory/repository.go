package urls_memory_repository

import (
	"sync"
	"url-shortener/internal/core/domain"
)

// Values are copied at the boundary so callers cannot mutate stored entries.
type URLsRepository struct {
	mu   sync.RWMutex
	urls map[string]domain.URL
}

func NewURLsRepository() *URLsRepository {
	return &URLsRepository{urls: make(map[string]domain.URL)}
}
