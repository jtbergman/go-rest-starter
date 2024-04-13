package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/routes"
)

// Starts the server and handles graceful shutdown
func serve(app *app.App) error {
	// Define the server
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      routes.Mux(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
	}

	// Create a shutdown channel to receive errors from the Shutdown() function
	shutdownError := make(chan error)

	// Start a background routine
	go func() {
		// Listen for catchable stop signals with a buffer
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		s := <-quit

		// Log the caught signal
		app.Logger.Info("shutting down server", "signal", s.String())

		// Begin shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- srv.Shutdown(ctx)
		}

		// Log a message to say we're waiting for any background tasks
		app.Logger.Info("completing background tasks", "addr", srv.Addr)
		app.BG.Wait()
		shutdownError <- nil
	}()

	// Log the server start and address
	app.Logger.Info("starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait to receive the return value from Shutdown()
	if err := <-shutdownError; err != nil {
		return err
	}

	// Log the successful shutdown
	app.Logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
