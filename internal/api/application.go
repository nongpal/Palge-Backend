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
		enabled          bool
		rps              float64
		burst            float64
		janitor_interval int64
		client_expiry    int64
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

	cfg.rl.enabled = getEnvAsBool("RATE_LIMIT_ENABLED", true)
	cfg.rl.rps = getEnvAsFloat("RATE_LIMIT_RPS", 2)
	cfg.rl.burst = getEnvAsFloat("RATE_LIMIT_BURST", 4)
	cfg.rl.janitor_interval = int64(getEnvAsInt("RATE_LIMIT_JANITOR_INTERVAL", 1))
	cfg.rl.client_expiry = int64(getEnvAsInt("RATE_LIMIT_CLIENT_EXPIRY", 10))
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

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseFloat(value, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
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
		cfg:         cfg,
		rateLimiter: NewRateLimiter(cfg.rl.rps, cfg.rl.burst),
		logger:      logger,
		models:      data.NewModels(db),
		mailer:      mailer,
	}
	return app, nil
}

func (app *Application) Close() error {
	return app.models.Accounts.DB.Close()
}
