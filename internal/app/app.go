package app

import (
	"go-rest-starter.jtbergman.me/internal/config"
	"go-rest-starter.jtbergman.me/internal/mailer"
	"go-rest-starter.jtbergman.me/internal/models"
	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

// Container for app wide dependencies
type App struct {
	BG     Backgrounder
	Config config.Config
	Logger xlogger.Logger
	Mailer mailer.Mailer
	Models *models.Models
	Rest   *rest.Rest
}

// Create a new App struct
func New(
	backgrounder Backgrounder,
	config config.Config,
	logger xlogger.Logger,
	mailer mailer.Mailer,
	models *models.Models,
	rest *rest.Rest,
) *App {
	return &App{
		BG:     backgrounder,
		Config: config,
		Logger: logger,
		Mailer: mailer,
		Models: models,
		Rest:   rest,
	}
}
