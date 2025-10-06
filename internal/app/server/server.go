package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
)

func Run(cfg config.Config, h *handlers.Handler, m *middleware.Middleware) error {
	r := getRouter(h, m)
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
