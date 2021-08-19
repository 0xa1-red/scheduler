// Package logging contains the logic for logging both
// human-readable and structured log messages
package logging

import "go.uber.org/zap"

// Logger holds a pointer to a logger
type Logger struct {
	*zap.SugaredLogger
}

// New returns a pointer to a Logger or an error
func New() (*Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &Logger{l.Sugar()}, nil
}

// MustNew creates a new logger or panics
func MustNew() *Logger {
	l, err := New()
	if err != nil {
		panic(err)
	}
	return l
}
