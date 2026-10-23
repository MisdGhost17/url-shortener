package domain

import (
	"errors"
	"testing"
	"time"
	core_errors "url-shortener/internal/core/errors"
)

func TestURLValidation(t *testing.T) {
	good := NewURL("aB123z", " https://example.org ", time.Now())
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	if good.OriginalURL != "https://example.org" {
		t.Fatal("URL not trimmed")
	}
	cases := []URL{good, good, good, good, good, good}
	cases[0].OriginalURL = ""
	cases[1].OriginalURL = "ftp://example.org"
	cases[2].ShortCode = "short"
	cases[3].ShortCode = "abcde!"
	cases[4].CreatedAt = time.Time{}
	cases[5].Clicks = -1
	for _, entry := range cases {
		if !errors.Is(entry.Validate(), core_errors.ErrInvalidArgument) {
			t.Fatalf("accepted invalid URL: %+v", entry)
		}
	}
}
