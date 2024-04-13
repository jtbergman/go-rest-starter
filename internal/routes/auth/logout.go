package auth

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/routes/middleware"
)

// ============================================================================
// POST
// ============================================================================

// Logs the user out by deleting their access token from the tokens table
func (app *Auth) logoutPost(w http.ResponseWriter, r *http.Request) {
	token := middleware.ContextGetToken(r)

	if _, err := app.tokens.Delete(token, tokens.ScopeAuthentication); err != nil {
		app.rest.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
