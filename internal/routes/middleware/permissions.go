package middleware

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// Requires a permission for a request to be performed
//
// Internally, this will require the user to be authenticated
func (mw *Middleware) RequirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)

		// Get permissions
		permissions, err := mw.permissions.GetByID(user.ID)
		if err != nil {
			mw.rest.Error(w, err)
			return
		}

		// Required permission exists
		err = xerrors.ClientUnauthorized(!permissions.Include(code), "middleware.RequirePermission.Include")
		if err != nil {
			mw.rest.Error(w, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return mw.Authenticated(fn)
}
