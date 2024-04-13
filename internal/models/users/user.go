package users

import (
	"errors"
	"fmt"
	"time"

	"go-rest-starter.jtbergman.me/internal/validator"
	"go-rest-starter.jtbergman.me/internal/xerrors"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

// ============================================================================
// Type
// ============================================================================

// Encapsulates the database properties of a user. The
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

// Create a new User
func new(email, plaintext string) (*User, *xerrors.AppError) {
	v := validator.New()
	v.IsEmail(email, "email", "is invalid")
	v.Check(len(plaintext) >= 8, "password", "must be at least 8 characters")

	if err := v.Valid("users.new.valid"); err != nil {
		return nil, err
	}

	user := &User{Email: email}
	err := user.SetPassword(plaintext)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Set a user's password
func (u *User) SetPassword(plaintext string) *xerrors.AppError {
	hash, err := hash(plaintext, "models.SetPassword")

	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

// Checks a user's password
func (user *User) PasswordMatches(plaintext string) (bool, *xerrors.AppError) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, xerrors.ServerError(
				"models.PasswordMatches",
				xerrors.ErrServerInternal,
			)
		}
	}

	return true, nil
}

// ============================================================================
// Anonymous User
// ============================================================================

// An empty User
var AnonymousUser = &User{}

// Checks if a user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// ============================================================================
// Helper
// ============================================================================

// Hashes a password
func hash(plaintext string, op string) (string, *xerrors.AppError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)

	if err != nil {
		return "", xerrors.ServerError(
			op,
			fmt.Errorf("%w:%v", xerrors.ErrServerInternal, err),
		)
	}

	return string(hash), nil
}
