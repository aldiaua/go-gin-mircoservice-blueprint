package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	"github.com/example/go-gin-blueprint/internal/config"
	"github.com/example/go-gin-blueprint/internal/database"
	"github.com/example/go-gin-blueprint/internal/logger"
	"github.com/example/go-gin-blueprint/internal/server"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.WithError(err).Error("failed to connect to database")
		return err
	}
	app := server.NewHTTPServer(cfg, db, log)
	defer app.Close()

	if err := app.Run(context.Background()); err != nil {
		log.WithError(err).Error("users service stopped unexpectedly")
		return err
	}
	log.Info("users service stopped")
	return nil
}
