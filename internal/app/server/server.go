package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
	"net/http"
)

const host = `localhost:8080`

func Run() error {
	handler := handlers.MakeHandler(
		storage.MakeStorage(),
		host,
	)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.MainHandler)

	return http.ListenAndServe(host, mux)
}