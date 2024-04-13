package auth

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/rest"
	"go-rest-starter.jtbergman.me/internal/routes/middleware"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ============================================================================
// POST
// ============================================================================

// Deletes an authenticated user
//
// Users must also provide their credentials to confirm the deletion.
func (app *Auth) deletePost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Get user from context
	authUser := middleware.ContextGetUser(r)

	// Get user from request
	if err := app.rest.ReadJSON(w, r, "auth.deletePost", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Get user from DB
	requestUser, err := app.users.GetByEmail(input.Email)
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Compare passwords
	passwordIsCorrect, err := requestUser.PasswordMatches(input.Password)
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Password is valid
	err = xerrors.ClientUnauthorized(!passwordIsCorrect, "auth.deletePost.Password")
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Users match
	err = xerrors.ClientUnauthorized(authUser.ID != requestUser.ID, "auth.deletePost.ID")
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	// Delete user
	if _, err := app.users.Delete(authUser); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Send ID
	env := rest.Envelope{"message": "Your account has been deleted"}
	app.rest.WriteJSON(w, "auth.deletePost", http.StatusOK, env)
}
