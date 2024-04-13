package app

import (
	"fmt"
	"sync"

	"go-rest-starter.jtbergman.me/internal/xerrors"
	"go-rest-starter.jtbergman.me/internal/xlogger"
)

// ============================================================================
// Interface
// ============================================================================

// Defines a type that can run background tasks
type Backgrounder interface {
	Run(fn func())
	Wait()
}

// ============================================================================
// Type
// ============================================================================

// Concrete implementation that runs background tasks with a wait group
type Background struct {
	logger xlogger.Logger
	wg     sync.WaitGroup
}

// Creates a new Background instance
func NewBackground(logger xlogger.Logger) *Background {
	return &Background{logger: logger}
}

// ============================================================================
// Type
// ============================================================================

// Runs a background task
func (bg *Background) Run(fn func()) {
	bg.wg.Add(1)

	go func() {
		defer bg.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				serverError := xerrors.ServerError(
					"app.Background",
					fmt.Errorf("%w: %v", xerrors.ErrServerInternal, err),
				)
				bg.logger.Error(serverError.Error())
			}
		}()

		fn()
	}()
}

// Waits for background tasks
func (bg *Background) Wait() {
	bg.wg.Wait()
}
