package middleware

import "net/http"

// Middleware to log requests
func (mw *Middleware) Requests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		mw.rest.Logger.Info("request", "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}
