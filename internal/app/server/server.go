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
	r.Post(
		"/",
		m.WithGzipReq(m.WithLog(h.HandlePost())),
	)
	r.Get(
		"/{shortUrl}",
		m.WithLog(h.HandleGet()),
	)
	r.Post(
		"/api/shorten",
		m.WithGzipReq(m.WithGzipResp(m.WithLog(h.HandleShorten()))),
	)
	return r
}