package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// Integration
// ============================================================================

// Helper to skip an integration test when running with the --short flag
func Integration(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.SkipNow()
	}
}

// Helper to decode a response body into a destination and stop tests on err
func Decode(t *testing.T, resp *http.Response, dst any) {
	t.Helper()

	err := json.NewDecoder(resp.Body).Decode(&dst)
	Check(t, err == nil)
}

// ============================================================================
// Handlers
// ============================================================================

// Defines a function on a generic type
type HandlerTestFunc[T any] func(t *testing.T, result T)

// Test case that allows verifying a specific body result
type HandlerTestCase[T any] struct {
	Name   string
	Auth   string
	Body   string
	Status int
	FN     HandlerTestFunc[T]
}

func RunHandlerTestCase[T any](
	t *testing.T,
	handler http.HandlerFunc,
	method string,
	url string,
	tc HandlerTestCase[T],
) {
	t.Helper()

	t.Run(tc.Name, func(t *testing.T) {
		req := httptest.NewRequest(method, url, bytes.NewBufferString(tc.Body))
		rr := httptest.NewRecorder()
		if tc.Auth != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.Auth))
		}

		handler.ServeHTTP(rr, req)
		resp := rr.Result()
		defer resp.Body.Close()

		Equal(t, resp.StatusCode, tc.Status)

		// Verify response is JSON unless no content
		if tc.FN != nil {
			Equal(t, resp.Header.Get("Content-Type"), "application/json")
			var body T
			Decode(t, resp, &body)

			// If a body validator was provided, validate
			if tc.FN != nil {
				tc.FN(t, body)
			}
		}
	})
}
