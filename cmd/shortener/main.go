package main

import (
	"net/http"
	"github.com/eduardtungatarov/shortener/internal/app"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.MainHandler)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
