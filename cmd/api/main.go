package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
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

	err := app.run()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) run() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.port),
		Handler: app.routes(),

		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.cfg.env,
	)
	return srv.ListenAndServe()
}
