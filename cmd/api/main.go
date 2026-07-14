package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	cfg := config{
		port: 4000,
		env:  "development",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.run()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) run() error {
	//mux := app.routes()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),

		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.config.env,
	)
	return srv.ListenAndServe()
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	return mux
}
