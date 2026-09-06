package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"latihan-fiber/app/repository"
	"latihan-fiber/app/service"
	"latihan-fiber/config"
	"latihan-fiber/database"
)

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("Gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)

	app := config.NewApp(logger, pool, userService)

	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("Server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Server berjalan", slog.String("port", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Gagal menutup server dengan rapi",
			slog.String("error", err.Error()))
	}

	logger.Info("Server berhenti dengan rapi")
}
