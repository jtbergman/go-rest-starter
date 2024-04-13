package users

import (
	"context"
	"time"

	"go-rest-starter.jtbergman.me/internal/models/core"
	"go-rest-starter.jtbergman.me/internal/models/tokens"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ============================================================================
// Interface
// ============================================================================

// Defines a mockable interface for user operations
type UsersRepository interface {
	Delete(user *User) (int64, *xerrors.AppError)
	GetByEmail(email string) (*User, *xerrors.AppError)
	GetByToken(plaintext string) (*User, *xerrors.AppError)
	Insert(user *User) *xerrors.AppError
	New(email, plaintext string) (*User, *xerrors.AppError)
	Update(user *User) *xerrors.AppError
}

func Repository(db core.Queryable) UsersRepository {
	return &Users{DB: db}
}

// ============================================================================
// Implementation
// ============================================================================

// Provides access to User database methods
type Users struct {
	DB core.Queryable
}

// Create User
func (Users) New(email, plaintext string) (*User, *xerrors.AppError) {
	user, err := new(email, plaintext)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Insert a given user
//
// Check for xerrors.ErrUniqueViolation for email conflicts.
//
// Sets the following properties on the provided users:
//
// User.ID
// User.Activated
// User.CreatedAt
// User.Version
func (m Users) Insert(user *User) *xerrors.AppError {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, activated, created_at, version
	`
	args := []any{user.Email, user.Password}
	dest := []any{&user.ID, &user.Activated, &user.CreatedAt, &user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		return xerrors.DatabaseError(err, "users.Insert")
	}

	return nil
}

// Gets the user by their email
func (m Users) GetByEmail(email string) (*User, *xerrors.AppError) {
	query := `
		SELECT id, email, password, activated, created_at, version
		FROM users
		WHERE email = $1
	`
	var user User
	dest := []any{&user.ID, &user.Email, &user.Password, &user.Activated, &user.CreatedAt, &user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, email).Scan(dest...); err != nil {
		return nil, xerrors.DatabaseError(err, "users.GetByEmail")
	}

	return &user, nil
}

// Gets the user from one of their tokens
func (m Users) GetByToken(plaintext string) (*User, *xerrors.AppError) {
	query := `
		SELECT users.id, users.email, users.password, users.activated, users.created_at, users.version
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.expiry > $2
	`
	var user User
	args := []any{tokens.Hash(plaintext), time.Now()}
	dest := []any{&user.ID, &user.Email, &user.Password, &user.Activated, &user.CreatedAt, &user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		return nil, xerrors.DatabaseError(err, "users.GetByToken")
	}

	return &user, nil
}

// Updates a user using optimistic locking.
//
// Be careful not to provide a user with default values for the set fields.
//
// Sets:
//
// email
// password
// activated
func (m Users) Update(user *User) *xerrors.AppError {
	query := `
		UPDATE users
		SET email = $1, password = $2, activated = $3, version = version + 1
		WHERE id = $4 and version = $5
		RETURNING version
	`
	args := []any{user.Email, user.Password, user.Activated, user.ID, user.Version}
	dest := []any{&user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		return xerrors.DatabaseError(err, "users.Update")
	}

	return nil
}

// Deletes a user
func (m Users) Delete(user *User) (int64, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	if err != nil {
		return 0, xerrors.DatabaseError(err, "users.Delete")
	}

	return core.RowsAffected(result, "users.Delete")
}
