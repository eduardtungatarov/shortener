package main

import (
	"context"
	"fmt"

	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/logger"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/eduardtungatarov/shortener/internal/app/server"
	"github.com/eduardtungatarov/shortener/internal/app/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctx := context.Background()

	log, err := logger.MakeLogger()
	if err != nil {
		panic(err)
	}

	cfg := config.LoadFromFlag()

	s, err := storage.MakeStorage(cfg)
	if err != nil {
		log.Fatalf("failed to make storage: %v", err)
	}
	defer s.Close()
	err = s.Load(ctx)
	if err != nil {
		log.Fatalf("failed to load storage: %v", err)
	}

	m := middleware.MakeMiddleware(log)
	h := handlers.MakeHandler(s, cfg.BaseURL, log)

	go h.DeleteBatch(ctx)

	err = server.Run(cfg, h, m)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
