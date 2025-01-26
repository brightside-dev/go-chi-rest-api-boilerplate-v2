package custom_error

import "fmt"

var (
	// General errors
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrUnAuthorized        = fmt.Errorf("unauthorized request")

	// Auth errors
	ErrInvalidEmailOrPassword = fmt.Errorf("invalid email or password")
	ErrInvalidRequestBody     = fmt.Errorf("invalid request body")

	// Request errors
	ErrMissingParam     = fmt.Errorf("missing required parameter")
	ErrInvalidParamType = fmt.Errorf("invalid parameter type")
)

// Use this function for system errors that will be logged
func NewSystemError(err error) error {
	return fmt.Errorf("internal server error: %w", err)
}
