package mocks

import (
	"testing"

	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/models"
	"go-rest-starter.jtbergman.me/internal/rest"
)

// Creates a mock App
func App(t *testing.T) *app.App {
	// Use mock config to get DSN
	cfg := cfg()

	// Create a new test db
	db := createTestDB(t, cfg.DB.DSN)

	// Create a shared logger
	logger := logger()

	mock := app.New(
		app.NewBackground(logger),
		cfg,
		logger,
		mail(),
		models.New(db),
		rest.New(logger),
	)

	return mock
}

// Provides access to the mock logger
func Logger(app *app.App) *mockLogger {
	return app.Logger.(*mockLogger)
}

// Provides access to the mock Mailer
func Mailer(app *app.App) *Mail {
	return app.Mailer.(*Mail)
}
