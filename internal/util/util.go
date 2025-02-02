package util

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func ParseBirthday(birthday interface{}) (time.Time, error) {
	switch v := birthday.(type) {
	case string:
		return time.Parse("2006-01-02", v)
	case []byte:
		return time.Parse("2006-01-02", string(v))
	default:
		return time.Time{}, fmt.Errorf("unexpected type for birthday: %T", v)
	}
}

func ParseDateTime(dateTime interface{}) (time.Time, error) {
	switch v := dateTime.(type) {
	case string:
		return time.Parse("2006-01-02 15:04:05", v)
	case []byte:
		return time.Parse("2006-01-02 15:04:05", string(v))
	default:
		return time.Time{}, fmt.Errorf("unexpected type for date time: %T", v)
	}
}

func GetHTTPRequestContext(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"user_agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
		"method":     r.Method,
		"path":       r.URL.Path,
	}
}

func MapToSlogAttrs(context map[string]interface{}) []slog.Attr {
	var attrs []slog.Attr
	for key, value := range context {
		switch v := value.(type) {
		case string:
			attrs = append(attrs, slog.String(key, v))
		case int:
			attrs = append(attrs, slog.Int(key, v))
		case bool:
			attrs = append(attrs, slog.Bool(key, v))
		case float64:
			attrs = append(attrs, slog.Float64(key, v))
		default:
			// If the value type is not supported, use Any
			attrs = append(attrs, slog.Any(key, v))
		}
	}
	return attrs
}

func LogWithContext(
	logger *slog.Logger,
	level slog.Level,
	msg string,
	additionalContext map[string]interface{},
	r *http.Request,
) {
	if r != nil {
		additionalContext = GetHTTPRequestContext(r)
	}

	// Convert the map to slog.Attr slice
	attrs := MapToSlogAttrs(additionalContext)

	// Create a slice to hold the arguments to pass to the logger
	var anyAttrs []interface{}
	for _, attr := range attrs {
		anyAttrs = append(anyAttrs, attr.Key, attr.Value)
	}

	// Log the message based on the specified level
	switch level {
	case slog.LevelError:
		logger.Error(msg, anyAttrs...)
	default:
		logger.Info(msg, anyAttrs...)
	}
}
