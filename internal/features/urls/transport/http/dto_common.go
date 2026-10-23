package urls_transport_http

import (
	"time"
	"url-shortener/internal/core/domain"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
	Original  string `json:"original"`
}

type StatsResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int       `json:"clicks"`
}

func statsDTO(entry domain.URL) StatsResponse {
	return StatsResponse{ShortCode: entry.ShortCode, OriginalURL: entry.OriginalURL, CreatedAt: entry.CreatedAt, Clicks: entry.Clicks}
}
