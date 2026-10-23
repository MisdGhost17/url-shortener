package urls_memory_repository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"url-shortener/internal/core/domain"
	core_errors "url-shortener/internal/core/errors"
)

func TestConcurrentCreateAndResolve(t *testing.T) {
	r := NewURLsRepository()
	ctx := context.Background()
	entry := domain.NewURL("abc123", "https://example.org", time.Now())
	var successes atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.CreateURL(ctx, entry)
			if err == nil {
				successes.Add(1)
			} else if !errors.Is(err, core_errors.ErrConflict) {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if successes.Load() != 1 {
		t.Fatalf("created %d entries with same code", successes.Load())
	}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := r.ResolveURL(ctx, entry.ShortCode); err != nil {
				t.Error(err)
			}
			if _, err := r.GetURLs(ctx); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	stored, err := r.GetURL(ctx, entry.ShortCode)
	if err != nil || stored.Clicks != 200 {
		t.Fatalf("lost clicks: %+v, %v", stored, err)
	}
	stored.Clicks = 999
	entries, _ := r.GetURLs(ctx)
	entries[0].OriginalURL = "modified"
	stored, _ = r.GetURL(ctx, entry.ShortCode)
	if stored.Clicks != 200 || stored.OriginalURL != entry.OriginalURL {
		t.Fatal("mutable storage leaked")
	}
	if err := r.DeleteURL(ctx, entry.ShortCode); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ResolveURL(ctx, entry.ShortCode); !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatal(err)
	}
	if _, err := r.GetURL(ctx, entry.ShortCode); !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatal(err)
	}
	if err := r.DeleteURL(ctx, entry.ShortCode); !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatal(err)
	}
}
func TestCancellationDoesNotMutateStorage(t *testing.T) {
	r := NewURLsRepository()
	live := context.Background()
	entry := domain.NewURL("abc123", "https://example.org", time.Now())
	r.CreateURL(live, entry)
	ctx, cancel := context.WithCancel(live)
	cancel()
	_, createErr := r.CreateURL(ctx, domain.NewURL("def456", "https://other.org", time.Now()))
	_, getErr := r.GetURL(ctx, entry.ShortCode)
	_, listErr := r.GetURLs(ctx)
	_, resolveErr := r.ResolveURL(ctx, entry.ShortCode)
	deleteErr := r.DeleteURL(ctx, entry.ShortCode)
	for _, err := range []error{createErr, getErr, listErr, resolveErr, deleteErr} {
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation: %v", err)
		}
	}
	entries, _ := r.GetURLs(live)
	if len(entries) != 1 || entries[0].Clicks != 0 {
		t.Fatalf("canceled mutation applied: %+v", entries)
	}
}
