package auth

import (
	"net/http"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
	"go-rest-starter.jtbergman.me/internal/routes/auth"
)

// Tests error cases for login
func TestLoginValidation(t *testing.T) {
	assert.Integration(t)

	app := mocks.App(t)
	handler := authHandler(app)

	type failure struct {
		Error map[string]string `json:"error"`
	}

	tests := []assert.HandlerTestCase[failure]{
		{
			Name:   "Email/Validation",
			Body:   `{"password": "password"}`,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error["email"], "must be provided")
			},
		},
		{
			Name:   "Password/Validation",
			Body:   `{"email": "test@example.com"}`,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error["password"], "must be provided")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, tc)
	}
}

// Tests error cases for login
func TestLoginUnauthorized(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	type success struct {
		Token string `json:"token"`
	}

	credentials := `{"email": "test@example.com", "password": "password"}`

	// User does not exist
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[failure]{
		Name:   "User/DoesNotExist",
		Body:   credentials,
		Status: http.StatusUnauthorized,
	})

	// Seed – create user
	assert.Check(t, registerUser(handler, credentials))

	// Password does not match
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[failure]{
		Name:   "User/WrongPassword",
		Body:   `{"email": "test@example.com", "password": "pa55word"}`,
		Status: http.StatusUnauthorized,
		FN: func(t *testing.T, result failure) {
			assert.Equal(t, result.Error, "The provided credentials are invalid")
		},
	})

	// User needs activation
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[failure]{
		Name:   "User/NeedsActivation",
		Body:   credentials,
		Status: http.StatusUnauthorized,
		FN: func(t *testing.T, result failure) {
			assert.Equal(t, result.Error, "Activate your account in order to sign in")
		},
	})

	// Seed - activate user
	assert.Check(t, activateUser(handler, app))

	// Success
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[success]{
		Name:   "Login/Success",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result success) {
			assert.True(t, len(result.Token) > 0)
		},
	})
}
