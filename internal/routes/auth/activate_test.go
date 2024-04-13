package auth

import (
	"fmt"
	"net/http"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
)

func TestActivate(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	// Seed - create user
	assert.Check(t, registerUser(handler, `{"email": "test@example.com", "password": "password"}`))

	// Invalid Token
	assert.RunHandlerTestCase[failure](t, handler, "PUT", ActivateRoute, assert.HandlerTestCase[failure]{
		Name:   "Activate/Invalid",
		Body:   `{"token": "token"}`,
		Status: http.StatusNotFound,
	})

	app.BG.Wait()
	token := mocks.Mailer(app).WelcomeActivationToken

	// Success
	assert.RunHandlerTestCase[user](t, handler, "PUT", ActivateRoute, assert.HandlerTestCase[user]{
		Name:   "Activate/Success",
		Body:   fmt.Sprintf(`{"token": "%s"}`, token),
		Status: http.StatusOK,
		FN: func(t *testing.T, result user) {
			assert.True(t, result.User.Activated)
		},
	})
}
