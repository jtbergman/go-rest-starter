package auth

import (
	"net/http"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
	"go-rest-starter.jtbergman.me/internal/mocks"
	"go-rest-starter.jtbergman.me/internal/models/users"
	"go-rest-starter.jtbergman.me/internal/routes/auth"
)

const registerSuccessBody = `{"email": "test@example.com", "password": "password"}`

// Test register success
func TestRegister(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	type success struct {
		User users.User `json:"user"`
	}

	tests := []assert.HandlerTestCase[success]{
		{
			Name:   "Success",
			Body:   registerSuccessBody,
			Status: http.StatusCreated,
			FN: func(t *testing.T, result success) {
				assert.True(t, result.User.ID > 0)
				assert.False(t, result.User.Activated)
				assert.Equal(t, result.User.Version, 0)
				assert.Equal(t, result.User.Email, "test@example.com")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, "POST", auth.RegisterRoute, tc)
	}

	t.Run("Success/WelcomeEmail", func(t *testing.T) {
		assert.Equal(t, mocks.Mailer(app).WelcomeCount, 1)
	})
}

// Test register validation
func TestRegisterValidation(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	tests := []assert.HandlerTestCase[failures]{
		{
			Name:   "Success",
			Body:   `{"email": "", "password": ""}`,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result failures) {
				assert.Equal(t, result.Error["email"], "is invalid")
				assert.Equal(t, result.Error["password"], "must be at least 8 characters")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			assert.RunHandlerTestCase(t, handler, "POST", auth.RegisterRoute, tc)
		})
	}
}

// Test register request
func TestRegisterFailure(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := authHandler(app)

	// Seed - create user
	assert.Check(t, registerUser(handler, registerSuccessBody))

	tests := []assert.HandlerTestCase[failure]{
		{
			Name:   "Conflict",
			Body:   registerSuccessBody,
			Status: http.StatusConflict,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error, "That email is already taken")
			},
		},
		{
			Name:   "BadRequest",
			Body:   ``,
			Status: http.StatusBadRequest,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error, "Request body cannot be empty")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, "POST", auth.RegisterRoute, tc)
	}
}
