package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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

	var wg sync.WaitGroup
	// Запускаем обработчик запросов на удаление ссылок.
	go func() {
		defer wg.Done()
		h.DeleteBatch(ctx)
	}()

	// Запускаем сервер.
	go func() {
		defer wg.Done()
		err = server.Run(ctx, cfg, h, m)
		if err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
