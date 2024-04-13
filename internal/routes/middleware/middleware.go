package middleware

import (
	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/models/permissions"
	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

type Middleware struct {
	logger      xlogger.Logger
	permissions permissions.PermissionsRepository
	rest        *rest.Rest
	users       users.UsersRepository
}

func New(app *app.App) *Middleware {
	return &Middleware{
		logger:      app.Logger,
		permissions: app.Models.Permissions,
		rest:        app.Rest,
		users:       app.Models.Users,
	}
}
