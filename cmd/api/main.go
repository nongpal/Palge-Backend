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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Bootstrap application")
}
