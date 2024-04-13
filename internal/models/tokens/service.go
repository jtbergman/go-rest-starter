package tokens

import (
	"context"
	"time"

	"go-rest-starter.jtbergman.me/internal/models/core"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ===========================================================================
// Interface
// ===========================================================================

type TokensRepository interface {
	New(userID int64, expiryDuration time.Duration, scope string) (*Token, *xerrors.AppError)
	Insert(token *Token) (int64, *xerrors.AppError)
	Delete(plaintext string, scope string) (int64, *xerrors.AppError)
	DeleteAllForScope(userID int64, scope string) (int64, *xerrors.AppError)
}

func Repository(db core.Queryable) TokensRepository {
	return &Tokens{DB: db}
}

// ===========================================================================
// Implementation
// ===========================================================================

// Provides access to the Tokens database methods
type Tokens struct {
	DB core.Queryable
}

// Creates a Token with the given user ID, expiry, and scope
//
// Scopes:
//
//	ScopeActivation
//	ScopeAuthentication
//	ScopePasswordReset
func (Tokens) New(userID int64, expiryDuration time.Duration, scope string) (*Token, *xerrors.AppError) {
	token, err := new(userID, expiryDuration, scope)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// Insert token
func (m Tokens) Insert(token *Token) (int64, *xerrors.AppError) {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope, token.CreatedAt, token.UpdatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, xerrors.DatabaseError(err, "tokens.Insert")
	}

	return core.RowsAffected(result, "tokens.Insert")
}

// Delete specific token
//
// Scopes:
//
//	ScopeActivation
//	ScopeAuthentication
//	ScopePasswordReset
func (m Tokens) Delete(plaintext string, scope string) (int64, *xerrors.AppError) {
	hash := Hash(plaintext)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, "DELETE FROM tokens WHERE hash = $1 AND scope = $2", hash, scope)
	if err != nil {
		return 0, xerrors.DatabaseError(err, "tokens.Delete")
	}

	return core.RowsAffected(result, "tokens.Delete")
}

// Delete all tokens with a given scope
//
// Scopes:
//
//	ScopeActivation
//	ScopeAuthentication
//	ScopePasswordReset
func (m Tokens) DeleteAllForScope(userID int64, scope string) (int64, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, "DELETE FROM tokens WHERE user_id = $1 AND scope = $2", userID, scope)
	if err != nil {
		return 0, xerrors.DatabaseError(err, "tokens.DeleteAllForScope")
	}

	return core.RowsAffected(result, "tokens.DeleteAllForScope")
}
