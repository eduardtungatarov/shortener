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
	r.Use(m.WithLog)

	r.Get(
		"/{shortUrl}",
		h.HandleGet(),
	)

	gzipReqG := r.Group(func(r chi.Router) {
		r.Use(m.WithGzipReq)
	})
	gzipReqG.Post(
		"/",
		h.HandlePost(),
	)
	gzipReqG.Group(func(r chi.Router) {
		r.Use(m.WithGzipResp)
		r.Post("/api/shorten", h.HandleShorten())
	})

	return r
}