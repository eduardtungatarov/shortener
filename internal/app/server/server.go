package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(cfg config.Config, h *handlers.Handler, m *middleware.Middleware) error {
	r := getRouter(h, m)
	return http.ListenAndServe(cfg.ServerHostPort, r)
}

func getRouter(h *handlers.Handler, m *middleware.Middleware) chi.Router {
	r := chi.NewRouter()
	r.Use(m.WithLog, m.WithAuth)

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
	})

	return r
}