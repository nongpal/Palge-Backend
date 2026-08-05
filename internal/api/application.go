package api

import (
	"flag"
	"log/slog"
	"os"
	"sync"

	"github.com/nongpal/Palge-Backend/internal/data"
	"github.com/nongpal/Palge-Backend/internal/mailer"
)

const version = "0.1.0"

type Config struct {
	Port int
	Env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type Application struct {
	cfg    Config
	logger *slog.Logger
	models data.Models
	mailer *mailer.Mailer
	wg     sync.WaitGroup
}

func NewConfig(cfg *Config) {
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/palge?sslmode=disable", "PostgreSQL DSN")

	// SMTP Config
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "8ee421160eb95b", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "f81a52d6f9a102", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Palge <no-reply@github.com/nongpal/Palge-Backend>", "SMTP sender")
	flag.Parse()
}

func NewApplication(cfg Config) (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	mailer, err := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		return nil, err
	}

	app := &Application{
		cfg:    cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer,
	}
	return app, nil
}

func (app *Application) Close() error {
	return app.models.Accounts.DB.Close()
}
