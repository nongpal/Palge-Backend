package main

import (
<<<<<<< HEAD
=======
	"context"
	"errors"
>>>>>>> aa3203b (feat: Implement graceful shutdown)
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.cfg.port),
		Handler:           app.routes(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
		ErrorLog:          slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

<<<<<<< HEAD
=======
	shutdownErr := make(chan error)

>>>>>>> aa3203b (feat: Implement graceful shutdown)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
<<<<<<< HEAD
		app.logger.Info("caught signal", "signal", s.String())

		os.Exit(0)
=======
		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownErr <- srv.Shutdown(ctx)
>>>>>>> aa3203b (feat: Implement graceful shutdown)
	}()

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.cfg.env,
	)

<<<<<<< HEAD
	return srv.ListenAndServe()
=======
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
>>>>>>> aa3203b (feat: Implement graceful shutdown)
}
