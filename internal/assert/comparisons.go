package assert

import (
	"errors"
	"testing"
)

// ============================================================================
// Equality
// ============================================================================

// Asserts that two values are equal
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// Asserts that two values are not equal
func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual == expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// Asserts that two errors are equal
func Is(t *testing.T, actual, expected error) {
	t.Helper()

	if !errors.Is(actual, expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// ============================================================================
// Booleans
// ============================================================================

// Asserts that the given value is true or stops execution
func Check(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Fatal("got: false, want: true")
	}
}

// Asserts that the given value is true
func True(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Error("got: false, want: true")
	}
}

// Asserts that the given value is false
func False(t *testing.T, value bool) {
	t.Helper()

	if value {
		t.Error("got: true; want: false")
	}
}
