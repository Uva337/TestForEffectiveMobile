// @title Subscription Service API
// @version 1.0
// @description This is a sample REST API for managing subscriptions.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"effective-mobile-task/internal/config"
	httpHandler "effective-mobile-task/internal/handler/http"
	"effective-mobile-task/internal/repository/postgres"
	"effective-mobile-task/internal/service"
	"effective-mobile-task/pkg/logger"
	_ "effective-mobile-task/docs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("не удалось загрузить конфигурацию", "error", err)
		os.Exit(1)
	}

	log := logger.New("local")
	log.Info("логгер инициализирован")
	log.Debug("отладочные сообщения включены")

	var dbPool *pgxpool.Pool
	const maxRetries = 5
	const retryDelay = 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Info("попытка подключения к БД", "attempt", i+1)
		dbPool, err = postgres.NewPostgresDB(cfg.Postgres)
		if err == nil {
			break
		}

		log.Warn("не удалось подключиться к БД, повтор...", "error", err, "next_attempt_in", retryDelay)
		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Error("не удалось подключиться к базе данных после нескольких попыток", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	log.Info("успешное подключение к базе данных")


	subRepo := postgres.NewSubscriptionRepository(dbPool)
	subService := service.NewSubscriptionService(subRepo)
	handler := httpHandler.NewHandler(subService, log)


	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler.RegisterRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}


	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	log.Info("сервер запускается", "port", cfg.HTTPPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("сервер аварийно остановлен", "error", err)
			os.Exit(1)
		}
	}()

	<-stop

	log.Info("сервер останавливается...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("ошибка при остановке сервера", "error", err)
		os.Exit(1)
	}

	log.Info("сервер успешно остановлен")
}

