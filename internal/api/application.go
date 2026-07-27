package api

import (
	"flag"
	"log/slog"
	"os"

	"github.com/nongpal/Palge-Backend/internal/data"
)

const version = "0.1.0"

type Config struct {
	Port int
	Env  string
	db   struct {
		dsn string
	}
}

type Application struct {
	cfg    Config
	logger *slog.Logger
	models data.Models

	// INFO: in memory account state will be remove!
	accounts []*data.Account
	nextID   int64
}

func NewConfig(cfg *Config) {
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/palge?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()
}

func NewApplication(cfg Config) (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	app := &Application{
		cfg:      cfg,
		logger:   logger,
		models:   data.NewModels(db),
		accounts: make([]*data.Account, 0),
		nextID:   1,
	}
	return app, nil
}

func (app *Application) Close() error {
	return app.models.Accounts.DB.Close()
}
