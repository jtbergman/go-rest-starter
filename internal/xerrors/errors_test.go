package xerrors

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/lib/pq"
	"go-rest-starter.jtbergman.me/internal/assert"
)

func TestDatabaseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		error      error
		wantStatus int
		wantError  error
	}{
		{
			name:       "ErrNoRows",
			error:      sql.ErrNoRows,
			wantStatus: http.StatusNotFound,
			wantError:  ErrNotFound,
		},
		{
			name:       "UniqueViolation",
			error:      &pq.Error{Code: "23505"},
			wantStatus: http.StatusConflict,
			wantError:  ErrUniqueViolation,
		},
		{
			name:       "ForeignKeyViolation",
			error:      &pq.Error{Code: "23503"},
			wantStatus: http.StatusNotFound,
			wantError:  ErrForeignKeyViolation,
		},
		{
			name:       "NullViolation",
			error:      &pq.Error{Code: "23502"},
			wantStatus: http.StatusBadRequest,
			wantError:  ErrNullViolation,
		},
		{
			name:       "CheckViolation",
			error:      &pq.Error{Code: "23514"},
			wantStatus: http.StatusBadRequest,
			wantError:  ErrCheckViolation,
		},
		{
			name:       "DeadlockDetected",
			error:      &pq.Error{Code: "40P01"},
			wantStatus: http.StatusInternalServerError,
			wantError:  ErrDeadlockDetected,
		},
		{
			name:       "Unexpected",
			error:      errors.New("some_random_error"),
			wantStatus: http.StatusInternalServerError,
			wantError:  ErrServerInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DatabaseError(tc.error, "TestOperation")
			assert.Equal(t, got.StatusCode, tc.wantStatus)
			assert.Is(t, got.Err, tc.wantError)
		})
	}
}

func TestClientUnauthorized(t *testing.T) {
	t.Parallel()

	t.Run("ClientUnauthorized/Nil", func(t *testing.T) {
		clientError := ClientUnauthorized(false, "xerrors.Nil")
		assert.True(t, clientError == nil)
	})

	t.Run("ClientUnauthorized/Error", func(t *testing.T) {
		clientError := ClientUnauthorized(true, "xerrors.Error")
		assert.True(t, clientError != nil)
		assert.Equal(t, clientError.StatusCode, http.StatusUnauthorized)
		assert.Equal(t, clientError.Op, "xerrors.Error")
		assert.Is(t, clientError, ErrUnauthorized)
	})
}
