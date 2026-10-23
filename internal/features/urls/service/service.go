package urls_service

import (
	"context"
	"url-shortener/internal/core/domain"
)

// URLsRepository belongs to the consumer, not to the storage adapter.
type URLsRepository interface {
	CreateURL(context.Context, domain.URL) (domain.URL, error)
	GetURL(context.Context, string) (domain.URL, error)
	GetURLs(context.Context) ([]domain.URL, error)
	DeleteURL(context.Context, string) error
	// ResolveURL must read and increment the click counter atomically.
	ResolveURL(context.Context, string) (domain.URL, error)
}

type CodeGenerator interface{ Generate() (string, error) }

type URLsService struct {
	repository URLsRepository
	generator  CodeGenerator
}

func NewURLsService(repository URLsRepository, generator CodeGenerator) *URLsService {
	return &URLsService{repository: repository, generator: generator}
}
