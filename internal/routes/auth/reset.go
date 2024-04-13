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

// Creates a password reset request by generating one-time tokens
func (auth *Auth) resetPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	// Parse email
	if err := auth.rest.ReadJSON(w, r, "auth.resetPost", &input); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Get user
	user, err := auth.users.GetByEmail(input.Email)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Verify active
	err = xerrors.ClientUnauthorized(!user.Activated, "auth.resetPost")
	if err != nil {
		err.Data = "Please activate your account to reset password"
		auth.rest.Error(w, err)
		return
	}

	// Create reset token
	token, err := auth.tokens.New(user.ID, time.Hour, tokens.ScopePasswordReset)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Insert the token into the database
	if _, err := auth.tokens.Insert(token); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Send an email to the user
	auth.bg.Run(func() {
		data := map[string]string{
			"passwordResetToken": token.Plaintext,
		}

		err := auth.mailer.SendPasswordResetEmail(user.Email, data)
		if err != nil {
			auth.logger.Error(err.Error())
		}
	})

	// Notify the user their request is processing
	env := rest.Envelope{"message": "An email will be sent with reset instructions"}
	auth.rest.WriteJSON(w, "auth.resetPost", http.StatusAccepted, env)
}

// ============================================================================
// PUT
// ============================================================================

// Updates the user's password if they have a valid reset token
func (auth *Auth) resetPut(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	// Parse password and token
	if err := auth.rest.ReadJSON(w, r, "auth.resetPut", &input); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(len(input.Password) >= 8, "password", "must be at least 8 characters")

	// Error if invalid
	if err := v.Valid("auth.resetPut"); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Get user
	user, err := auth.users.GetByToken(input.Token)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Set password
	if err := user.SetPassword(input.Password); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Update user
	if err := auth.users.Update(user); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Delete token
	if _, err := auth.tokens.DeleteAllForScope(user.ID, tokens.ScopePasswordReset); err != nil {
		auth.rest.Error(w, err)
		return
	}

	env := rest.Envelope{"message": "Your password was reset successfully"}
	auth.rest.WriteJSON(w, "auth.resetPut", http.StatusOK, env)
}
