package main

import (
	"log/slog"
	"os"

	"github.com/nongpal/Palge-Backend/internal/data"
)

const version = "0.1.0"

type config struct {
	port int
	env  string
}

type application struct {
	cfg      config
	logger   *slog.Logger
	accounts []*data.Account
	nextID   int64
}

func main() {
	cfg := config{
		port: 4000,
		env:  "development",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		cfg:      cfg,
		logger:   logger,
		accounts: make([]*data.Account, 0),
		nextID:   1,
	}

	err := app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
