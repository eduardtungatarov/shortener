package main

import (
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/server"
)

func main() {
	cfg := config.LoadFromFlag()
	err := server.Run(cfg)
	if err != nil {
		panic(err)
	}
}
