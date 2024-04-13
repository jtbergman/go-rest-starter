// xerrors provides centralized error handling for the application
//
// All methods should return *AppError and wrap errors using DatabaseError.
// Use AppError.Matches to check if an error has a particular case.
package xerrors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

// ============================================================================
// Type
// ============================================================================

// Abstracts errors from the database into easy-to-check types
var (
	ErrCheckViolation      = errors.New("check_violation")
	ErrDeadlockDetected    = errors.New("deadlock_detected")
	ErrForeignKeyViolation = errors.New("foreign_key_violation")
	ErrNotFound            = errors.New("not_found")
	ErrNullViolation       = errors.New("null_violation")
	ErrUniqueViolation     = errors.New("unique_violation")
)

// Abstract errors from the api into easy-to-check types
var (
	ErrBadRequest       = errors.New("bad_request")
	ErrEntityTooLarge   = errors.New("entity_too_large")
	ErrFailedValidation = errors.New("failed_validation")
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrUnauthorized     = errors.New("unauthorized")
)

// Abstracts unknown errors
var (
	ErrServerInternal = errors.New("server_error")
	ErrMailerInternal = errors.New("mailer_error")
)

// ============================================================================
// Type
// ============================================================================

// Represents a class of errors that can be returned directly to clients
type AppError struct {
	StatusCode int
	Data       any
	Op         string
	Err        error
}

// If an error is a specific error type, replace the data
func (e *AppError) If(target error, fn func(err *AppError)) {
	if e.Matches(target) {
		fn(e)
	}
}

// Check if an AppError is a specific error type
func (e *AppError) Matches(target error) bool {
	return errors.Is(e, target)
}

// Supports unwrapping the errors for use with errors.Is(...)
func (e *AppError) Unwrap() error {
	return e.Err
}

// Makes this struct a true error type
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// ============================================================================
// Common Usage
// ============================================================================

// Creates a client error with the specified status code and data
func ClientError(statusCode int, data any, op string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Data:       data,
		Op:         op,
		Err:        err,
	}
}

// Returns an authorization error if the value is true
func ClientUnauthorized(value bool, op string) *AppError {
	if value {
		return ClientError(
			http.StatusUnauthorized,
			"Authorization required",
			op,
			ErrUnauthorized,
		)
	}
	return nil
}

// Creates a server error with the appropriate status code and message
func ServerError(op string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Data:       "A server error occurred",
		Op:         op,
		Err:        err,
	}
}

// ============================================================================
// Database Errors
// ============================================================================

// Wraps an error returned from a database operations
func DatabaseError(err error, op string) *AppError {
	if errors.Is(err, sql.ErrNoRows) {
		return ClientError(
			http.StatusNotFound,
			"The requested resource does not exist",
			op,
			fmt.Errorf("%w: %v", ErrNotFound, err),
		)
	}

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return ClientError(
				http.StatusConflict,
				"A resource with that identity already exists",
				op,
				fmt.Errorf("%w: %v", ErrUniqueViolation, err),
			)

		case "23503":
			return ClientError(
				http.StatusNotFound,
				"A referenced resource does not exist",
				op,
				fmt.Errorf("%w: %v", ErrForeignKeyViolation, err),
			)

		case "23502":
			return ClientError(
				http.StatusBadRequest,
				"The resource contained an unexpected null value",
				op,
				fmt.Errorf("%w: %v", ErrNullViolation, err),
			)

		case "23514":
			return ClientError(
				http.StatusBadRequest,
				"The resource contains invalid data",
				op,
				fmt.Errorf("%w: %v", ErrCheckViolation, err),
			)

		case "40P01":
			return ServerError(
				op,
				fmt.Errorf("%w: %v", ErrDeadlockDetected, err),
			)
		}
	}

	return ServerError(
		op,
		fmt.Errorf("%w: %v", ErrServerInternal, err),
	)
}
