package auth

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/mailer"
	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/routes/middleware"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

// ============================================================================
// Auth Type
// ============================================================================

// Encapsulates the Application dependencies required by routes
type Auth struct {
	bg     app.Backgrounder
	logger xlogger.Logger
	mailer mailer.Mailer
	rest   *rest.Rest
	tokens tokens.TokensRepository
	users  users.UsersRepository
}

func New(app *app.App) *Auth {
	return &Auth{
		bg:     app.BG,
		logger: app.Logger,
		mailer: app.Mailer,
		rest:   app.Rest,
		tokens: app.Models.Tokens,
		users:  app.Models.Users,
	}
}

// ============================================================================
// Route
// ============================================================================

func (auth *Auth) Route(mux *http.ServeMux, mw *middleware.Middleware) {
	mux.HandleFunc(ActivateRoute, auth.Activate)

	mux.HandleFunc(DeleteRoute, mw.Authenticated(auth.Delete))

	mux.HandleFunc(LoginRoute, auth.Login)

	mux.HandleFunc(LogoutRoute, mw.Authenticated(auth.Logout))

	mux.HandleFunc(RegisterRoute, auth.Register)

	mux.HandleFunc(ResetRoute, auth.Reset)
}

// ============================================================================
// Activate
// ============================================================================

const ActivateRoute = "/v1/auth/activate"

func (app *Auth) Activate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/activate.html")

	case "PUT":
		app.activatePut(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "GET, PUT")
	}
}

// ============================================================================
// Delete
// ============================================================================

const DeleteRoute = "/v1/auth/delete"

func (app *Auth) Delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		app.deletePost(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "POST")
	}
}

// ============================================================================
// Login
// ============================================================================

const LoginRoute = "/v1/auth/login"

func (app *Auth) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		app.loginPost(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "POST")
	}
}

// ============================================================================
// Logout
// ============================================================================

const LogoutRoute = "/v1/auth/logout"

func (app *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		app.logoutPost(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "POST")
	}
}

// ============================================================================
// Register
// ============================================================================

const RegisterRoute = "/v1/auth/register"

func (auth *Auth) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		auth.registerPost(w, r)

	default:
		auth.rest.MethodNotAllowed(w, r, "POST")
	}
}

// ============================================================================
// Reset
// ============================================================================

const ResetRoute = "/v1/auth/reset"

func (app *Auth) Reset(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/reset.html")

	case "POST":
		app.resetPost(w, r)

	case "PUT":
		app.resetPut(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "GET, POST, PUT")
	}
}
