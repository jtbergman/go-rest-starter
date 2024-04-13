package auth

import (
	"fmt"
	"net/http"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
)

func TestReset(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Bad Request
	assert.RunHandlerTestCase[failure](t, handler, "POST", ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/BadRequest",
		Body:   ``,
		Status: http.StatusBadRequest,
	})

	// User DNE
	assert.RunHandlerTestCase[failure](t, handler, "POST", ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/UserDNE",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusNotFound,
	})

	// Seed - create user
	assert.Check(t, registerUser(handler, credentials))

	// User Inactive
	assert.RunHandlerTestCase[failure](t, handler, "POST", ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/UserDNE",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusUnauthorized,
	})

	// Seed – activate user
	assert.Check(t, activateUser(handler, app))

	// Success
	assert.RunHandlerTestCase[message](t, handler, "POST", ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Success",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusAccepted,
	})

	// Wait for reset token
	app.BG.Wait()

	// Invalid Token
	assert.RunHandlerTestCase[failure](t, handler, "PUT", ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/BadToken",
		Body:   `{"password": "pa55word", "token": "token"}`,
		Status: http.StatusNotFound,
	})

	// Invalid password
	assert.RunHandlerTestCase[failures](t, handler, "PUT", ResetRoute, assert.HandlerTestCase[failures]{
		Name:   "Reset/InvalidPassword",
		Body:   fmt.Sprintf(`{"password": "please", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusUnprocessableEntity,
	})

	// Success
	assert.RunHandlerTestCase[failures](t, handler, "PUT", ResetRoute, assert.HandlerTestCase[failures]{
		Name:   "Reset/InvalidPassword",
		Body:   fmt.Sprintf(`{"password": "pa55word", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusOK,
	})
}
