package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(cfg config.Config) error {
	h := handlers.MakeHandler(
		storage.MakeStorage(),
		cfg.BaseURL,
	)
	r := getRouter(h)

	return http.ListenAndServe(cfg.ServerHostPort, r)
}

func getRouter(h *handlers.Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.HandlePost)
	r.Get("/{shortUrl}", h.HandleGet)

	return r
}