package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const host = `localhost:8080`

func Run() error {
	h := handlers.MakeHandler(
		storage.MakeStorage(),
		host,
	)
	r := getRouter(h)

	return http.ListenAndServe(host, r)
}

func getRouter(h *handlers.Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.HandlePost)
	r.Get("/{shortUrl}", h.HandleGet)

	return r
}