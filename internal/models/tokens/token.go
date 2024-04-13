package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ============================================================================
// Constants
// ============================================================================

const (
	ScopeActivation     = "activate"
	ScopeAuthentication = "authneticate"
	ScopePasswordReset  = "reset"
)

// ============================================================================
// Token
// ============================================================================

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	Scope     string    `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// New Token
func new(userID int64, expiryDuration time.Duration, scope string) (*Token, *xerrors.AppError) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, xerrors.ServerError(
			"tokens.new",
			xerrors.ErrServerInternal,
		)
	}

	// Convert the random bytes to base64 to get the plaintext
	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := Hash(plaintext)

	// Create the token with the duration added to the current time
	now := time.Now()
	return &Token{
		Plaintext: plaintext,
		Hash:      hash,
		UserID:    userID,
		Expiry:    now.Add(expiryDuration),
		Scope:     scope,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Performs a fast hash of a plaintext token
func Hash(plaintext string) []byte {
	hash := sha256.Sum256([]byte(plaintext))
	return hash[:]
}
