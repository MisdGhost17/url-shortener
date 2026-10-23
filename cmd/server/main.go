package main

import (
	"log"
	"net/http"
	"time"
	"url-shortener/internal/core/codegen"
	core_http_middleware "url-shortener/internal/core/transport/http/middleware"
	core_http_server "url-shortener/internal/core/transport/http/server"
	urls_memory_repository "url-shortener/internal/features/urls/repository/memory"
	urls_service "url-shortener/internal/features/urls/service"
	urls_transport_http "url-shortener/internal/features/urls/transport/http"
	web_transport_http "url-shortener/internal/features/web/transport/http"
)

func newHandler(staticDirectory string) http.Handler {
	repository := urls_memory_repository.NewURLsRepository()
	service := urls_service.NewURLsService(repository, codegen.RandomGenerator{})
	web := web_transport_http.NewWebHTTPHandler(staticDirectory)
	urls := urls_transport_http.NewURLsHTTPHandler(service, http.HandlerFunc(web.NotFound))
	mux := http.NewServeMux()
	core_http_server.RegisterRoutes(mux, urls.Routes()...)
	core_http_server.RegisterRoutes(mux, web.Routes()...)
	return core_http_middleware.CORS(core_http_middleware.Logging(mux))
}

func main() {
	srv := &http.Server{
		Addr: ":8080", Handler: newHandler("web/static"),
		ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second,
	}
	log.Println("URL Shortener running on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
