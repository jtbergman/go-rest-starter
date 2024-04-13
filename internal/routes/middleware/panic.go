package middleware

import (
	"fmt"
	"net/http"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// Recovers a panic, logs an error, and sends a JSON response
func (mw *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				serverError := xerrors.ServerError(
					fmt.Sprintf("Panic.%s.%s", r.Method, r.URL.RequestURI()),
					fmt.Errorf("%w: %v", xerrors.ErrServerInternal, err),
				)

				mw.rest.Error(w, serverError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
