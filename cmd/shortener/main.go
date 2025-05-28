package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/logger"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/eduardtungatarov/shortener/internal/app/server"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
)

func main() {
	ctx := context.Background()

	log, err := logger.MakeLogger()
	if err != nil {
		panic(err)
	}

	cfg := config.LoadFromFlag()

	s, err := storage.MakeStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("failed to make storage: %v", err)
	}
	err = s.Load()
	if err != nil {
		log.Fatalf("failed to load storage: %v", err)
	}

	db, err := storage.MakeDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to make database: %v", err)
	}
	defer db.SQLDB.Close()

	m := middleware.MakeMiddleware(log)
	h := handlers.MakeHandler(ctx, s, cfg.BaseURL, db)

	err = server.Run(cfg, h, m)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
