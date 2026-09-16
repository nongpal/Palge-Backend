package api

import (
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
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

	rl struct {
		enabled bool
		rps     int
		burst   int
	}
}

type Application struct {
	cfg         Config
	rateLimiter *RateLimiter
	logger      *slog.Logger
	models      data.Models
	mailer      *mailer.Mailer
	wg          sync.WaitGroup
}

func NewConfig(cfg *Config) {
	_ = godotenv.Load()

	cfg.Port = getEnvAsInt("PORT", 4000)
	cfg.Env = getEnv("ENV", "development")
	cfg.db.dsn = getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5433/palge?sslmode=disable")

	cfg.smtp.host = getEnv("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	cfg.smtp.port = getEnvAsInt("SMTP_PORT", 2525)
	cfg.smtp.username = getEnv("SMTP_USERNAME", "")
	cfg.smtp.password = getEnv("SMTP_PASSWORD", "")
	cfg.smtp.sender = getEnv("SMTP_SENDER", "Palge <no-reply@github.com/nongpal/Palge-Backend>")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
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
