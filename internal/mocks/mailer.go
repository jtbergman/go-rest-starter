package mocks

import (
	"sync"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// Mock mailer to test for emails sent
type Mail struct {
	mu                     sync.Mutex
	WelcomeCount           int
	WelcomeActivationToken string
	PasswordResetCount     int
	PasswordResetToken     string
}

// Create a mock mail
func mail() *Mail {
	return &Mail{}
}

// Sends a welcome email
func (m *Mail) SendWelcomeEmail(recipient string, data map[string]string) *xerrors.AppError {
	m.mu.Lock()
	m.WelcomeCount += 1
	m.WelcomeActivationToken = data["activateToken"]
	m.mu.Unlock()
	return nil
}

// Sends a password reset email
func (m *Mail) SendPasswordResetEmail(recipient string, data map[string]string) *xerrors.AppError {
	m.mu.Lock()
	m.PasswordResetCount += 1
	m.PasswordResetToken = data["passwordResetToken"]
	m.mu.Unlock()
	return nil
}
