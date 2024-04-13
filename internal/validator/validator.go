package validator

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// A Validator type which contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// Create validator
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Returns nil if valid otherwise a failed validation error
func (v *Validator) Valid(op string) *xerrors.AppError {
	if len(v.Errors) == 0 {
		return nil
	}

	return xerrors.ClientError(
		http.StatusUnprocessableEntity,
		v.Errors,
		op,
		xerrors.ErrFailedValidation,
	)
}

// Adds an error message to the map if the key does not exist
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Adds an error message to the map if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
