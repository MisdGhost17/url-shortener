package urls_transport_http

import (
	"context"
	"net/http"
	"url-shortener/internal/core/domain"
	core_http_server "url-shortener/internal/core/transport/http/server"
)

type URLsService interface {
	CreateURL(context.Context, string) (domain.URL, error)
	GetURL(context.Context, string) (domain.URL, error)
	GetURLs(context.Context) ([]domain.URL, error)
	DeleteURL(context.Context, string) error
	ResolveURL(context.Context, string) (domain.URL, error)
}

type URLsHTTPHandler struct {
	service  URLsService
	notFound http.Handler
}

func NewURLsHTTPHandler(service URLsService, notFound http.Handler) *URLsHTTPHandler {
	return &URLsHTTPHandler{service: service, notFound: notFound}
}
func (h *URLsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/api/shorten", Handler: h.ShortenURL},
		{Path: "/api/stats/", Handler: h.Stats},
		{Path: "/api/list", Handler: h.ListAll},
		{Method: http.MethodDelete, Path: "/api/delete/", Handler: h.DeleteURL},
		{Method: http.MethodGet, Path: "/", Handler: h.Redirect},
	}
}
