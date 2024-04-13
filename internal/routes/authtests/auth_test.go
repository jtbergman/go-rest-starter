package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-rest-starter.jtbergman.me/internal/app"
	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/routes/auth"
	"go-rest-starter.jtbergman.me/internal/routes/middleware"
)

func TestAuthE2E(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	// Shared credentials
	var bearer string
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Shared responses
	type token struct {
		Token string `json:"token"`
	}

	// Register
	assert.RunHandlerTestCase(t, handler, "POST", auth.RegisterRoute, assert.HandlerTestCase[user]{
		Name:   "Register",
		Body:   credentials,
		Status: http.StatusCreated,
		FN: func(t *testing.T, result user) {
			assert.False(t, result.User.Activated)
			assert.Equal(t, result.User.Email, "test@example.com")

			app.BG.Wait()
			assert.NotEqual(t, mocks.Mailer(app).WelcomeActivationToken, "")
		},
	})

	// Activate
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ActivateRoute, assert.HandlerTestCase[user]{
		Name:   "Activate",
		Body:   fmt.Sprintf(`{"token": "%s"}`, mocks.Mailer(app).WelcomeActivationToken),
		Status: http.StatusOK,
		FN: func(t *testing.T, result user) {
			assert.True(t, result.User.Activated)
		},
	})

	// Login
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[token]{
		Name:   "Login/1",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result token) {
			bearer = result.Token
			assert.NotEqual(t, bearer, "")
		},
	})

	// Logout
	assert.RunHandlerTestCase(t, handler, "POST", auth.LogoutRoute, assert.HandlerTestCase[struct{}]{
		Name:   "Logout",
		Auth:   bearer,
		Body:   ``,
		Status: http.StatusNoContent,
		FN:     nil,
	})

	// Request Reset
	assert.RunHandlerTestCase(t, handler, "POST", auth.ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Post",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusAccepted,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "An email will be sent with reset instructions")

			app.BG.Wait()
			assert.NotEqual(t, mocks.Mailer(app).PasswordResetToken, "")
		},
	})

	// Reset Password
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Put",
		Body:   fmt.Sprintf(`{"password": "pa55word", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your password was reset successfully")
		},
	})

	// Login
	credentials = `{"email": "test@example.com", "password": "pa55word"}`
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[token]{
		Name:   "Login/2",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result token) {
			bearer = result.Token
			assert.NotEqual(t, bearer, "")
		},
	})

	// Delete
	assert.RunHandlerTestCase(t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[message]{
		Name:   "Delete",
		Auth:   bearer,
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your account has been deleted")
		},
	})
}

// ============================================================================
// Helpers
// ============================================================================

// Creates a complete Auth handler including middleware
func authHandler(app *app.App) http.HandlerFunc {
	handler := func() http.Handler {
		mux := http.NewServeMux()

		middleware := middleware.New(app)
		auth := auth.New(app)
		auth.Route(mux, middleware)

		return middleware.User(mux)
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// Helper user type
type user struct {
	User users.User `json:"user"`
}

// Helper success type
type message struct {
	Message string `json:"message"`
}

// Helper failure type
type failure struct {
	Error string `json:"error"`
}

// Helper failures type
type failures struct {
	Error map[string]string `json:"error"`
}

// ============================================================================
// Seeds
// ============================================================================

// Helper to activate a user
func activateUser(handler http.HandlerFunc, app *app.App) bool {
	app.BG.Wait()
	body := fmt.Sprintf(`{"token": "%s"}`, mocks.Mailer(app).WelcomeActivationToken)
	return sendRequest(handler, "PUT", auth.ActivateRoute, body) == http.StatusOK
}

// Helper to login a user
func loginUser(handler http.HandlerFunc, credentials string) string {
	var result struct {
		Token string `json:"token"`
	}
	sendRequestGetResult(handler, "POST", auth.LoginRoute, credentials, &result)
	return result.Token
}

// Helper to create a user
func registerUser(handler http.HandlerFunc, credentials string) bool {
	statusCode := sendRequest(handler, "POST", auth.RegisterRoute, credentials)
	return statusCode == http.StatusCreated
}

// Sends a request and returns the HTTP status
func sendRequest(handler http.HandlerFunc, method, route, body string) int {
	req := httptest.NewRequest(method, route, bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	defer resp.Body.Close()
	return resp.StatusCode
}

func sendRequestGetResult[T any](handler http.HandlerFunc, method, route, body string, dst *T) *T {
	req := httptest.NewRequest(method, route, bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&dst)
	return dst
}
