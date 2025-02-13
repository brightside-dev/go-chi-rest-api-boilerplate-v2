package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/rand"
)

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

// WithTransaction is a helper function that wraps database transaction logic.
// It accepts a context and a function (transactionCallback) containing the transaction logic.
// The transaction is committed if the function succeeds, or rolled back if the function fails or panics.
func WithTransaction(ctx context.Context, db *sql.DB, transactionCallback func(tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback/commit and handle panic
	defer func() {
		if p := recover(); p != nil {
			// Recover from panic and log the stack trace
			log.Printf("Panic occurred: %v", p)
			tx.Rollback() // Ensure rollback on panic
			panic(p)
		} else if err != nil {
			log.Printf("Error occurred: %v", err)
			// If error occurred, rollback transaction
			tx.Rollback()
		}
	}()

	var result interface{}

	// Run the transaction logic (which will contain the business logic)
	if result, err = transactionCallback(tx); err != nil {
		return nil, err
	}

	// Commit without shadowing the err variable
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

func GenerateVerificationCode() string {
	// Seed the random number generator to ensure random results each time
	rand.Seed(uint64(time.Now().UnixNano())) // Convert int64 to uint64

	// Define characters for letters (A-Z) and digits (0-9)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	// Create a slice for the result
	code := make([]rune, 5)

	// Randomly place 3 letters and 2 digits
	for i := 0; i < 3; i++ {
		code[i] = rune(letters[rand.Intn(len(letters))])
	}
	for i := 3; i < 5; i++ {
		code[i] = rune(digits[rand.Intn(len(digits))])
	}

	// Shuffle the slice to randomize the positions of letters and digits
	rand.Shuffle(5, func(i, j int) {
		code[i], code[j] = code[j], code[i]
	})

	return string(code)
}

func DD(data ...interface{}) {
	for _, d := range data {
		dumpData(d)
	}
	os.Exit(1)
}

func dumpData(d interface{}) {
	switch v := d.(type) {
	case validator.ValidationErrors:
		for _, err := range v {
			fmt.Printf("Field: %s, Tag: %s, ActualValue: %v\n", err.Field(), err.Tag(), err.Value())
		}
		return

	default:
		// Pretty-print JSON if it's a struct, map, or slice
		if reflect.TypeOf(d).Kind() == reflect.Struct || reflect.TypeOf(d).Kind() == reflect.Map || reflect.TypeOf(d).Kind() == reflect.Slice {
			prettyJSON, err := json.MarshalIndent(d, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
				return
			}
		}

		// Fallback to default formatting
		fmt.Printf("%#v\n", d)
	}
}

func FormatValidationError(errs validator.ValidationErrors) error {
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(fmt.Sprintf("Field: %s, Tag: %s, ActualValue: %v\n", err.Field(), err.Tag(), err.Value()))
	}
	return fmt.Errorf("%s", sb.String())
}
