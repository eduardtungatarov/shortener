// Package server роутинг и ф-я запуска http сервера приложения.
package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
)

// Run запуск http сервера приложения.
func Run(ctx context.Context, cfg config.Config, h *handlers.Handler, m *middleware.Middleware) error {
	var server *http.Server
	serverErr := make(chan error, 1)
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
		go func() {
			serverErr <- server.ListenAndServeTLS("", "")
		}()
	}

	if !cfg.EnableHTTPS {
		server := &http.Server{Addr: cfg.ServerHostPort, Handler: r}
		go func() {
			serverErr <- server.ListenAndServe()
		}()
	}

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return server.Shutdown(context.Background())
	}
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
