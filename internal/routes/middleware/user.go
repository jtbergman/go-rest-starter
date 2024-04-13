package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ===========================================================================
// Authenticate Middleware
// ===========================================================================

// Adds a user to the request context. If there is no user (i.e. an
// Authorization header was not provided), then an AnonymousUser will be added
// to the request. If there is a token, but it does not map to an
// authenticated user, then an authorization error will be returned.
func (mw *Middleware) User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "Vary: Authorization" tells caches this response varies based on
		// the value of the Authorization head in the request
		w.Header().Add("Vary", "Authorization")

		// Read the token from the header
		token, validHeader := readAuthorizationHeader(r)
		err := xerrors.ClientUnauthorized(!validHeader, "middleware.Authenticate")
		if err != nil {
			mw.rest.Error(w, err)
			return
		}
		if token == "" {
			r = contextSetToken(r, token)
			r = contextSetUser(r, users.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the user's details and add them to the context
		user, err := mw.users.GetByToken(token)
		if err != nil {
			err.If(xerrors.ErrNotFound, func(err *xerrors.AppError) {
				err.StatusCode = http.StatusUnauthorized
				err.Data = "Auth token is invalid"
			})
			mw.rest.Error(w, err)
			return
		}

		// Add the user to the request context
		r = contextSetToken(r, token)
		r = contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// Context: User
// ===========================================================================

// A custom contextKey type to prevent key collisions
type contextKey string

// The contextKey for storing the request user
const userContextKey = contextKey("user")

// Retrieves the User struct from the request context. This value is set by
// Authentication middleware and can be trusted. However, you should always
// check for user.IsAnonymous().
func ContextGetUser(r *http.Request) *users.User {
	user, ok := r.Context().Value(userContextKey).(*users.User)

	if !ok {
		panic("missing user value in request context")
	}

	return user
}

// Returns a new copy of the request with the User struct added to the context
func contextSetUser(r *http.Request, user *users.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// ===========================================================================
// Context: Token
// ===========================================================================

// The contextKey for storing the extracted authorization token
const tokenContextKey = contextKey("token")

// Retrieves the token string from the request context. This value is set by
// Authentication middleware and can be trusted. However, you should always
// check for an an empty string (no user). If you use ContextGetUser, checking
// for user.IsAnonymous() is sufficient, and checking for an empty string is
// not ncessary.
func ContextGetToken(r *http.Request) string {
	token, ok := r.Context().Value(tokenContextKey).(string)

	if !ok {
		panic("missing token string in request context")
	}

	return token
}

// Returns a new copy of the request with the token added to the context
func contextSetToken(r *http.Request, token string) *http.Request {
	ctx := context.WithValue(r.Context(), tokenContextKey, token)
	return r.WithContext(ctx)
}

// ===========================================================================
// Helper
// ===========================================================================

// Reads the authorization header from the request. The token will be
// returned if the header is valid, an empty string will be returned
// if the  head is not present. If the header is malformed, a false
// ok value will be returned.
func readAuthorizationHeader(r *http.Request) (string, bool) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return authorizationHeader, true
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", false
	}

	return headerParts[1], true
}
