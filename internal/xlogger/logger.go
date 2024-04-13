package xlogger

import "log/slog"

// ============================================================================
// Interace
// ============================================================================

type Logger interface {
	Handler() slog.Handler
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
