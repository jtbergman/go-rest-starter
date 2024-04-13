package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
	"go-rest-starter.jtbergman.me/internal/config"
	"go-rest-starter.jtbergman.me/internal/xerrors"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

//go:embed "templates"
var templateFS embed.FS

// ============================================================================
// Interface
// ============================================================================

type Mailer interface {
	SendWelcomeEmail(recipient string, data map[string]string) *xerrors.AppError
	SendPasswordResetEmail(recipientemail string, data map[string]string) *xerrors.AppError
}

// ============================================================================
// Type
// ============================================================================

// The Mailer struct defines an SMTP server and the sender information
type Mail struct {
	dialer *mail.Dialer
	logger xlogger.Logger
	sender string
	skip   bool
}

// ============================================================================
// Implementation
// ============================================================================

const (
	welcomeTemplate       = "user_welcome.tmpl"
	passwordResetTemplate = "password_reset.tmpl"
)

// Creates a new Mailer
func New(cfg config.Config, logger xlogger.Logger) Mailer {
	isLocal := cfg.IsLocal()
	var dialer *mail.Dialer

	if isLocal {
		dialer = mail.NewDialer(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
		)
		dialer.Timeout = 5 * time.Second
	}

	return &Mail{
		dialer: dialer,
		logger: logger,
		sender: cfg.SMTP.Sender,
		skip:   isLocal,
	}
}

// Sends a welcome email
func (m Mail) SendWelcomeEmail(recipient string, data map[string]string) *xerrors.AppError {
	if m.skip {
		m.logger.Info("Welcome Email", data["activateToken"])
	}
	return m.send(recipient, welcomeTemplate, data)
}

// Sends a password reset email
func (m Mail) SendPasswordResetEmail(recipient string, data map[string]string) *xerrors.AppError {
	if m.skip {
		m.logger.Info("Password Reset", data["passwordResetToken"])
	}
	return m.send(recipient, passwordResetTemplate, data)
}

// ============================================================================
// Private
// ============================================================================

// Sends an email to a recipient using the specified template
func (m Mail) send(recipient, templateFile string, data any) *xerrors.AppError {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return xerrors.ServerError(
			"mailer.Send.ParseFS",
			fmt.Errorf("%w: %v", xerrors.ErrMailerInternal, err),
		)
	}

	// Execute the subject template with the given data
	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return xerrors.ServerError(
			"mailer.Send.Subject",
			fmt.Errorf("%w: %v", xerrors.ErrMailerInternal, err),
		)
	}

	// Execute the plaintext template with the given data
	plainBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
		return xerrors.ServerError(
			"mailer.Send.PlainBody",
			fmt.Errorf("%w: %v", xerrors.ErrMailerInternal, err),
		)
	}

	// Execute the htmlBody template with the given data
	htmlBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
		return xerrors.ServerError(
			"mailer.Send.HTMLBody",
			fmt.Errorf("%w: %v", xerrors.ErrMailerInternal, err),
		)
	}

	// Create the email
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Send the email with retry or return the last error
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)

		if nil == err {
			return nil
		}

		time.Sleep(time.Second)
	}

	return xerrors.ServerError(
		"mailer.Send.DialAndSend",
		fmt.Errorf("%w: %v", xerrors.ErrMailerInternal, err),
	)
}
