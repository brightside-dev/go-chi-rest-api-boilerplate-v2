package logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"runtime"
	"time"
)

// DBLogHandler logs messages to the database.
type DBLogHandler struct {
	db          *sql.DB
	minLogLevel slog.Level
}

// NewDBLogHandler creates a new DBLogHandler.
func NewDBLogHandler(db *sql.DB, minLogLevel slog.Level) *DBLogHandler {
	return &DBLogHandler{
		db:          db,
		minLogLevel: minLogLevel,
	}
}

// Enabled checks if the log level is enabled.
func (h *DBLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.minLogLevel
}

// Handle writes the log record to the database.
func (h *DBLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	// Extract message and attributes from the log record.
	var msg string
	attrs := make(map[string]interface{})
	var source string

	// Capture the message directly from the Record's Message field.
	msg = r.Message

	// Extract other attributes from the log record.
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "source" {
			// Capture source information (file and line)
			source = attr.Value.String()
		} else {
			attrs[attr.Key] = attr.Value.Any()
		}
		return true
	})

	// If the source is empty, manually capture it from the runtime.
	if source == "" {
		_, file, line, ok := runtime.Caller(5) // Capture file and line number, adjust the call depth as needed
		if ok {
			source = fmt.Sprintf("%s:%d", file, line)
		}
	}

	// Marshal attributes into a JSON string.
	attrJSON, err := json.Marshal(attrs)
	if err != nil {
		log.Printf("Failed to marshal attributes: %v", err)
		return err
	}

	// Insert the log into the database.
	_, err = h.db.ExecContext(ctx,
		"INSERT INTO logs (level, message, attributes, source, created_at) VALUES (?, ?, ?, ?, ?)",
		r.Level.String(), msg, string(attrJSON), source, time.Now(),
	)
	if err != nil {
		log.Printf("Failed to write log to database: %v", err)
		return err
	}

	return nil
}

// WithAttrs adds attributes to the handler (not implemented for DBLogHandler).
func (h *DBLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // Return the same handler as this is a simple implementation.
}

// WithGroup adds a group to the handler (not implemented for DBLogHandler).
func (h *DBLogHandler) WithGroup(name string) slog.Handler {
	return h // Return the same handler as this is a simple implementation.
}
