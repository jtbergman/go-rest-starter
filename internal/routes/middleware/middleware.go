package middleware

import (
	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/models/permissions"
	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/rest"
)

type Middleware struct {
	permissions permissions.PermissionsRepository
	rest        *rest.Rest
	users       users.UsersRepository
}

func New(app *app.App) *Middleware {
	return &Middleware{
		permissions: app.Models.Permissions,
		rest:        app.Rest,
		users:       app.Models.Users,
	}
}
