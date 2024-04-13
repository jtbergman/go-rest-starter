package rest

import (
	"net/http"

	"go-rest-starter.jtbergman.me/internal/xerrors"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

// ============================================================================
// Interface
// ============================================================================

type Rest struct {
	Logger xlogger.Logger
}

func New(logger xlogger.Logger) *Rest {
	return &Rest{Logger: logger}
}

// ============================================================================
// Methods
// ============================================================================

func (rest *Rest) Error(w http.ResponseWriter, err *xerrors.AppError) {
	rest.Logger.Error(err.Error())
	rest.WriteJSON(w, err.Op, err.StatusCode, Envelope{"error": err.Data})
}

func (rest *Rest) MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowed string) {
	w.Header().Set("Allow", allowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	rest.Logger.Error("Method Not Allowed", "method", r.Method, "uri", r.URL.RequestURI())
}
