package auth

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/rest"
)

// ============================================================================
// PUT
// ============================================================================

// Activates the user's account given their activation token
func (app *Auth) activatePut(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
	}

	// Read activation token
	if err := app.rest.ReadJSON(w, r, "auth.activatePut", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Get user
	user, err := app.users.GetByToken(input.Token)
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Activate user
	user.Activated = true
	if err := app.users.Update(user); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Delete activation token
	if _, err := app.tokens.Delete(input.Token, tokens.ScopeActivation); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Send the updated user
	app.rest.WriteJSON(w, "auth.activatePut", http.StatusOK, rest.Envelope{"user": user})
}
