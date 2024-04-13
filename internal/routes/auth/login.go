package auth

import (
	"net/http"
	"time"

	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/validator"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ============================================================================
// POST
// ============================================================================

// Validates the user-provided credentials and returns an access token if valid
func (app *Auth) loginPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "auth.loginPost", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.Email) > 0, "email", "must be provided")
	v.Check(len(input.Password) > 0, "password", "must be provided")
	if err := v.Valid("auth.loginPost"); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Get user
	user, err := app.users.GetByEmail(input.Email)
	if err != nil {
		err.If(xerrors.ErrNotFound, func(err *xerrors.AppError) {
			err.StatusCode = http.StatusUnauthorized
			err.Data = "The provided credentials are invalid"
		})
		app.rest.Error(w, err)
		return
	}

	// Verify password
	match, err := user.PasswordMatches(input.Password)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if !match {
		clientError := xerrors.ClientError(
			http.StatusUnauthorized,
			"The provided credentials are invalid",
			"auth.loginPost",
			xerrors.ErrUnauthenticated,
		)
		app.rest.Error(w, clientError)
		return
	}

	// Verify active
	if !user.Activated {
		clientError := xerrors.ClientError(
			http.StatusUnauthorized,
			"Activate your account in order to sign in",
			"auth.loginPost",
			xerrors.ErrUnauthorized,
		)
		app.rest.Error(w, clientError)
		return
	}

	// Create token
	token, err := app.tokens.New(user.ID, 30*24*time.Hour, tokens.ScopeAuthentication)
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Insert token
	if _, err := app.tokens.Insert(token); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Send response
	app.rest.WriteJSON(w, "auth.loginPost", http.StatusOK, rest.Envelope{"token": token.Plaintext})
}
