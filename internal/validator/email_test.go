package validator

import (
	"strings"
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		Name    string
		Email   string
		IsEmail bool
	}{
		{
			Name:    "Valid/1",
			Email:   "test@example.com",
			IsEmail: true,
		},
		{
			Name:    "Valid/2",
			Email:   "test.user+123456789@example.com",
			IsEmail: true,
		},
		{
			Name:    "Valid/3",
			Email:   "a@b.c",
			IsEmail: true,
		},
		{
			Name:    "Invalid/1",
			Email:   "go-rest-starter",
			IsEmail: false,
		},
		{
			Name:    "Invalid/2",
			Email:   "go @ rest . starter",
			IsEmail: false,
		},
		{
			Name:    "Invalid/3",
			Email:   "test@example..com",
			IsEmail: false,
		},
		{
			Name:    "Invalid/4",
			Email:   "",
			IsEmail: false,
		},
		{
			Name:    "Invalid/5",
			Email:   strings.Repeat("A", 255),
			IsEmail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			v := New()
			v.IsEmail(tc.Email, "email", "is invalid")
			if tc.IsEmail {
				assert.True(t, v.Valid(tc.Name) == nil)
			} else {
				assert.False(t, v.Valid(tc.Name) == nil)
			}
		})
	}
}
