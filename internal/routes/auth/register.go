package auth

import (
	"net/http"
	"time"

	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ============================================================================
// POST
// ============================================================================

// Registers a user with a given email and password and responds with http.StatusCreated
func (auth *Auth) registerPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request
	if err := auth.rest.ReadJSON(w, r, "auth.registerPost", &input); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Create user
	user, err := auth.users.New(input.Email, input.Password)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Insert user
	if err := auth.users.Insert(user); err != nil {
		err.If(xerrors.ErrUniqueViolation, func(err *xerrors.AppError) {
			err.Data = "That email is already taken"
		})
		auth.rest.Error(w, err)
		return
	}

	// Create activation token
	token, err := auth.tokens.New(user.ID, 7*24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Insert activation token
	if _, err := auth.tokens.Insert(token); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Send welcome email
	auth.bg.Run(func() {
		data := map[string]string{
			"activateToken": token.Plaintext,
		}

		err := auth.mailer.SendWelcomeEmail(user.Email, data)
		if err != nil {
			auth.logger.Error(err.Error())
		}
	})

	// Send the user response
	auth.rest.WriteJSON(w, "auth.registerPost", http.StatusCreated, rest.Envelope{"user": user})
}
