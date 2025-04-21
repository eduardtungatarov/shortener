package main

import (
	"github.com/eduardtungatarov/shortener/internal/app/server"
)

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
