package main

import (
	"log/slog"
	"os"
)

type config struct {
	port int
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	cfg := config{port: 4000}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application{
		config: cfg,
		logger: logger,
	}

	app.logger.Info("Bootstrap application")
}
