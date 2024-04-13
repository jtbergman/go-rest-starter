package mocks

import (
	"log/slog"
	"os"
	"sync"
)

// ============================================================================
// Types
// ============================================================================

// Instance of the mock logger
type mockLogger struct {
	logger *slog.Logger
	logs   []logEntry
	mu     sync.Mutex
	record bool
}

// Represents a log entry
type logEntry struct {
	level   string
	message string
	args    []any
}

// ============================================================================
// Mock
// ============================================================================

// Mock Logger
func logger() *mockLogger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &mockLogger{logger: logger, logs: []logEntry{}, record: false}
}

// Begin capture
func (l *mockLogger) Begin() {
	l.mu.Lock()
	l.logs = []logEntry{}
	l.record = true
	l.mu.Unlock()
}

// End capture
func (l *mockLogger) End() {
	l.mu.Lock()
	l.record = false
	for _, log := range l.logs {
		switch log.level {
		case "DEBUG":
			l.logger.Debug(log.message, log.args...)

		case "INFO":
			l.logger.Info(log.message, log.args...)

		case "ERROR":
			l.logger.Error(log.message, log.args...)
		}
	}
	l.logs = nil
	l.mu.Unlock()
}

// ============================================================================
// Interface
// ============================================================================

func (l *mockLogger) Handler() slog.Handler {
	return l.logger.Handler()
}

func (l *mockLogger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.record {
		entry := logEntry{level: "DEBUG", message: msg, args: args}
		l.logs = append(l.logs, entry)
	}
}

func (l *mockLogger) Info(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.record {
		entry := logEntry{level: "INFO", message: msg, args: args}
		l.logs = append(l.logs, entry)
	}
}

func (l *mockLogger) Error(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.record {
		entry := logEntry{level: "ERROR", message: msg, args: args}
		l.logs = append(l.logs, entry)
	}
}
