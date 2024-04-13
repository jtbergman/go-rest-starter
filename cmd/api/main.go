package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/config"
	"go-rest-starter.jtbergman.me/internal/mailer"
	"go-rest-starter.jtbergman.me/internal/models"
	"go-rest-starter.jtbergman.me/internal/rest"
)

func main() {
	// Create dependencies
	config := config.New()
	database := OpenDatabase(config.DB.DSN)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Log if successful connection
	logger.Info("database connection pool established")

	// Create App
	app := app.New(
		app.NewBackground(logger),
		config,
		logger,
		mailer.New(config, logger),
		models.New(database),
		rest.New(logger),
	)

	if err := serve(app); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
