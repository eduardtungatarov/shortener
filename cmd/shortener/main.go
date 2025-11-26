package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/logger"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/eduardtungatarov/shortener/internal/app/server"
	shortenerService "github.com/eduardtungatarov/shortener/internal/app/service/shortener"
	"github.com/eduardtungatarov/shortener/internal/app/storage"

	v1 "github.com/eduardtungatarov/shortener/internal/contracts/shortener/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	srv := shortenerService.New(s, cfg.BaseURL)
	h := handlers.MakeHandler(s, cfg.BaseURL, log, cfg.TrustedSubnet, srv)

	grp, ctx := errgroup.WithContext(ctx)

	// Запускаем обработчик запросов на удаление ссылок.
	grp.Go(func() error {
		return h.DeleteBatch(ctx)
	})

	// Запускаем http сервер.
	grp.Go(func() error {
		return server.Run(ctx, cfg, h, m)
	})

	// Запускаем grpc сервер.
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor))
	grpcHandler := handlers.NewGrpcHandler(srv)
	v1.RegisterShortenerServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)
	grp.Go(func() error {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			return err
		}
		serverErr := make(chan error, 1)
		go func() {
			serverErr <- grpcServer.Serve(lis)
		}()
		select {
		case err := <-serverErr:
			return err
		case <-ctx.Done():
			grpcServer.GracefulStop()
			return nil
		}
	})

	if err := grp.Wait(); err != nil {
		log.Info("error stopping the service: %v", err)
	}
	log.Info("service has been stopped")
}
