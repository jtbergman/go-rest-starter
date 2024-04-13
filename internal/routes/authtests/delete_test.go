package auth

import (
	"net/http"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
	"go-rest-starter.jtbergman.me/internal/routes/auth"
)

func TestDelete(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Seed – create user, activate user, login user
	assert.Check(t, registerUser(handler, credentials))
	assert.Check(t, activateUser(handler, app))
	token := loginUser(handler, credentials)
	assert.Check(t, len(token) > 0)

	// Auth Required
	assert.RunHandlerTestCase[failures](t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/AuthRequired",
		Body:   credentials,
		Status: http.StatusUnauthorized,
	})

	// User Not Found
	assert.RunHandlerTestCase[failures](t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/UserNotFound",
		Body:   `{"email": "test2@example.com", "password": "password"}`,
		Auth:   token,
		Status: http.StatusNotFound,
	})

	// Credentials Invalid
	assert.RunHandlerTestCase[failures](t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/CredentialsInvalid",
		Body:   `{"email": "test@example.com", "password": "pa55word"}`,
		Auth:   token,
		Status: http.StatusUnauthorized,
	})

	// Success
	assert.RunHandlerTestCase[message](t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[message]{
		Name:   "Delete/CredentialsInvalid",
		Body:   credentials,
		Auth:   token,
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your account has been deleted")
		},
	})
}
