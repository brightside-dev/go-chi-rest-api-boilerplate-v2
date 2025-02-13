package logger

import (
	"context"
	"log/slog"
)

type MultiHandler interface {
	Enabled(ctx context.Context, level slog.Level) bool
	Handle(ctx context.Context, r slog.Record) error
	WithAttrs(attrs []slog.Attr) slog.Handler
	WithGroup(name string) slog.Handler
}

// MultiHandler combines multiple handlers.
type multiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a new MultiHandler with the provided handlers.
func NewMultiHandler(handlers ...slog.Handler) MultiHandler {
	return &multiHandler{handlers: handlers}
}

// Handle writes the log to all handlers.
func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range m.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// Enabled checks if any handler is enabled for the log level.
func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithAttrs adds attributes to the MultiHandler (for all handlers).
func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m.handlers {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}
	return &multiHandler{handlers: newHandlers}
}

// WithGroup adds a group to the MultiHandler (for all handlers).
func (m *multiHandler) WithGroup(name string) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m.handlers {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}
	return &multiHandler{handlers: newHandlers}
}
