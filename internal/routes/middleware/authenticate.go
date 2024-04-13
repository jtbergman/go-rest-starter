package middleware

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// Requires a user to be authenticated for the request to proceed
func (mw *Middleware) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)

		err := xerrors.ClientUnauthorized(user.IsAnonymous(), "middleware.Authenticated")
		if err != nil {
			mw.rest.Error(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
