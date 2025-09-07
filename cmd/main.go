package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/shirr9/order-api/internal/config"
	"github.com/shirr9/order-api/internal/handlers"
	"github.com/shirr9/order-api/internal/kafka"
	"github.com/shirr9/order-api/internal/logger"
	"github.com/shirr9/order-api/internal/service"
	"github.com/shirr9/order-api/internal/storage/cache"
	"github.com/shirr9/order-api/internal/storage/postgresql"
)

func main() {
	CfgPath := "configs/config.yaml"
	cfg, err := config.Load(CfgPath)
	if err != nil {
		slog.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	logWriter, err := logger.NewWriter(cfg.LogPath)
	if err != nil {
		slog.Error("failed to init log writer", slog.Any("err", err))
		os.Exit(1)
	}
	log := logger.NewLogger(cfg.Env, logWriter)
	log.Info("logger initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := postgresql.New(ctx, cfg)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		os.Exit(1)
	}
	defer storage.Close()
	log.Info("database connected")
	repo := storage.NewPostgresRepository()

	cacheStorage, err := cache.NewRedis(ctx, cfg)
	if err != nil {
		log.Error("failed to init cache", slog.Any("err", err))
		os.Exit(1)
	}
	defer cacheStorage.Close()
	log.Info("cache connected")

	orderService := service.NewService(repo, log, cacheStorage)
	kafkaConsumer := kafka.NewConsumer(cfg.KafkaConfig, log, orderService)
	go func() {
		kafkaConsumer.Run(ctx)
	}()

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "file://", "null"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}).Handler)

	router.Get("/order/{id}", handlers.NewIdHandler(log, orderService))

	router.Post("/order", handlers.NewAddHandler(log, orderService))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start http server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	log.Info("application started", slog.String("env", cfg.Env))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	log.Info("stopping application", slog.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown failed", slog.Any("error", err))
	}
	cancel()

	if err := kafkaConsumer.Close(); err != nil {
		log.Error("failed to close kafka consumer", slog.Any("err", err))
	}

	log.Info("application stopped")
}
