package validator

import (
	"testing"

	"go-rest-starter.jtbergman.me/internal/assert"
)

func TestValidator(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		v := New()
		assert.True(t, v.Valid("") == nil)
	})

	t.Run("Check", func(t *testing.T) {
		v := New()
		v.Check(true, "error", "invalid")
		assert.True(t, v.Valid("") == nil)

		v.Check(1 == 2, "error", "invalid")
		assert.True(t, v.Valid("") != nil)
	})

	t.Run("AddError", func(t *testing.T) {
		v := New()
		v.AddError("error", "our request could not be processed")
		assert.True(t, v.Valid("") != nil)
	})
}
