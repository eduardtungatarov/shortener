// Package server роутинг и ф-я запуска http сервера приложения.
package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"

	"golang.org/x/crypto/acme/autocert"
)

// Run запуск http сервера приложения.
func Run(cfg config.Config, h *handlers.Handler, m *middleware.Middleware) error {
	r := getRouter(h, m)

	if cfg.EnableHTTPS {
		path := strings.Split(cfg.ServerHostPort, ":")
		host := path[0]

		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(host),
		}

		server := &http.Server{
			Addr:      host + ":443",
			Handler:   r,
			TLSConfig: manager.TLSConfig(),
		}
		return server.ListenAndServeTLS("", "")
	}

	return http.ListenAndServe(cfg.ServerHostPort, r)
}

func getRouter(h *handlers.Handler, m *middleware.Middleware) chi.Router {
	r := chi.NewRouter()
	r.Use(m.WithLog, m.WithAuth)

	r.Mount("/debug", chiMiddleware.Profiler())

	r.Get(
		"/{shortUrl}",
		h.HandleGet,
	)

	r.Get(
		"/ping",
		h.HandleGetPing,
	)

	r.Get(
		"/api/user/urls",
		h.HandleGetUserUrls,
	)

	gzipReqG := r.Group(func(r chi.Router) {
		r.Use(m.WithGzipReq)
	})
	gzipReqG.Post(
		"/",
		h.HandlePost,
	)
	gzipReqG.Group(func(r chi.Router) {
		r.Use(m.WithGzipResp, m.WithJSONReqCheck)
		r.Post("/api/shorten", h.HandleShorten)
		r.Post("/api/shorten/batch", h.HandleShortenBatch)
		r.Delete(
			"/api/user/urls",
			h.HandleDeleteUserUrls,
		)
	})

	return r
}
