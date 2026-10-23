package urls_service

import (
	"context"
	"errors"
	"testing"
	"time"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
	memory "url-shortener/internal/features/urls/repository/memory"
)

type generatorFunc func() (string, error)

func (f generatorFunc) Generate() (string, error) { return f() }
func TestCreateRetriesCollisionWithoutOverwrite(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewURLsRepository()
	existing := domain.NewURL("aaaaaa", "https://existing.org", time.Now())
	repo.CreateURL(ctx, existing)
	calls := 0
	generator := generatorFunc(func() (string, error) {
		calls++
		if calls == 1 {
			return "aaaaaa", nil
		}
		return "bbbbbb", nil
	})
	service := NewURLsService(repo, generator)
	created, err := service.CreateURL(ctx, " https://new.org ")
	if err != nil || calls != 2 || created.ShortCode != "bbbbbb" || created.OriginalURL != "https://new.org" {
		t.Fatalf("create: %+v, %v, calls=%d", created, err, calls)
	}
	stored, _ := repo.GetURL(ctx, "aaaaaa")
	if stored != existing {
		t.Fatal("collision overwrote original URL")
	}
}
func TestCreateErrorsAndCancellation(t *testing.T) {
	repo := memory.NewURLsRepository()
	sentinel := errors.New("entropy unavailable")
	calls := 0
	generator := generatorFunc(func() (string, error) { calls++; return "", sentinel })
	service := NewURLsService(repo, generator)
	if _, err := service.CreateURL(context.Background(), "bad"); !errors.Is(err, core_errors.ErrInvalidArgument) || calls != 0 {
		t.Fatalf("invalid input: %v, calls=%d", err, calls)
	}
	if _, err := service.CreateURL(context.Background(), "https://example.org"); !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := service.CreateURL(ctx, "https://example.org"); !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("cancellation: %v, calls=%d", err, calls)
	}
	invalid := NewURLsService(repo, generatorFunc(func() (string, error) { return "!", nil }))
	if _, err := invalid.CreateURL(context.Background(), "https://example.org"); !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatal(err)
	}
	entries, _ := repo.GetURLs(context.Background())
	if len(entries) != 0 {
		t.Fatal("failed create persisted data")
	}
}

type failingRepository struct {
	URLsRepository
	err error
}

func (r failingRepository) CreateURL(context.Context, domain.URL) (domain.URL, error) {
	return domain.URL{}, r.err
}
func TestRepositoryErrorIsPreserved(t *testing.T) {
	sentinel := errors.New("storage unavailable")
	service := NewURLsService(failingRepository{err: sentinel}, generatorFunc(func() (string, error) { return "abc123", nil }))
	if _, err := service.CreateURL(context.Background(), "https://example.org"); !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
}
func TestCollisionRetryCanBeCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	service := NewURLsService(failingRepository{err: core_errors.ErrConflict}, generatorFunc(func() (string, error) {
		calls++
		if calls == 3 {
			cancel()
		}
		return "abc123", nil
	}))
	if _, err := service.CreateURL(ctx, "https://example.org"); !errors.Is(err, context.Canceled) || calls != 3 {
		t.Fatalf("retry: calls=%d, err=%v", calls, err)
	}
}
