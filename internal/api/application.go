package api

import (
	"log/slog"
	"os"

	"github.com/nongpal/Palge-Backend/internal/data"
)

const version = "0.1.0"

type Config struct {
	Port int
	Env  string
}

type Application struct {
	cfg      Config
	logger   *slog.Logger
	accounts []*data.Account
	nextID   int64
}

func NewApplication(cfg Config) *Application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Application{
		cfg:      cfg,
		logger:   logger,
		accounts: make([]*data.Account, 0),
		nextID:   1,
	}
}
