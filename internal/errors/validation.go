package errors

import (
	"fmt"
	"net/http"
)

// NewValidationError creates a validation-specific error
func NewValidationError(field string, reason string) *Error {
	return NewError(
		http.StatusBadRequest,
		fmt.Sprintf("Validation failed for field '%s': %s", field, reason),
		SeverityLow,
		CategoryValidation,
		false,
		nil,
	).WithSuggestion(fmt.Sprintf("Please check the value provided for '%s'", field))
}

// Common validation errors
var (
	ErrMissingRequired = func(field string) *Error {
		return NewValidationError(field, "field is required")
	}

	ErrInvalidFormat = func(field string, format string) *Error {
		return NewValidationError(field, fmt.Sprintf("must match format: %s", format))
	}
)
