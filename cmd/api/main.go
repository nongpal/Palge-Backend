package main

import (
	"log/slog"
	"os"
)

const version = "0.1.0"

type config struct {
	port int
	env  string
}

type application struct {
	cfg    config
	logger *slog.Logger
}

func main() {
	cfg := config{
		port: 4000,
		env:  "development",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
