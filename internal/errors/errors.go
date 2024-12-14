package errors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

const (
	SeverityLow ErrorSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s ErrorSeverity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ErrorCategory represents the type of error
type ErrorCategory int

const (
	CategoryNetwork ErrorCategory = iota
	CategoryValidation
	CategorySecurity
	CategoryDatabase
	CategoryThirdParty
	CategoryInternal
)

func (c ErrorCategory) String() string {
	switch c {
	case CategoryNetwork:
		return "network"
	case CategoryValidation:
		return "validation"
	case CategorySecurity:
		return "security"
	case CategoryDatabase:
		return "database"
	case CategoryThirdParty:
		return "third_party"
	case CategoryInternal:
		return "internal"
	default:
		return "unknown"
	}
}

// Error represents a service error with enhanced metadata
type Error struct {
	Code               int           `json:"code"`
	Message            string        `json:"message"`
	Severity           ErrorSeverity `json:"severity"`
	Category           ErrorCategory `json:"category"`
	Timestamp          time.Time     `json:"timestamp"`
	Retryable          bool          `json:"retryable"`
	Context            interface{}   `json:"context,omitempty"`
	TraceID            string        `json:"trace_id"`
	Source             string        `json:"source"`
	RecoverySuggestion string        `json:"recovery_suggestion,omitempty"`
	Documentation      string        `json:"documentation_url,omitempty"`
	RelatedErrors      []string      `json:"related_errors,omitempty"`
	Err                error         `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v (Severity: %s, Category: %s)",
			e.Code, e.Message, e.Err, e.Severity.String(), e.Category.String())
	}
	return fmt.Sprintf("[%d] %s (Severity: %s, Category: %s)",
		e.Code, e.Message, e.Severity.String(), e.Category.String())
}

// ErrorStats represents error statistics
type ErrorStats struct {
	Count     int64
	LastError *Error
	FirstSeen time.Time
	LastSeen  time.Time
	Sources   map[string]int64
}

// NewError creates a new error with enhanced metadata
func NewError(code int, message string, severity ErrorSeverity, category ErrorCategory, retryable bool, err error) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Severity:  severity,
		Category:  category,
		Timestamp: time.Now(),
		Retryable: retryable,
		Err:       err,
		TraceID:   fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

// Builder methods
func (e *Error) WithSuggestion(suggestion string) *Error {
	e.RecoverySuggestion = suggestion
	return e
}

func (e *Error) WithCategory(category ErrorCategory) *Error {
	e.Category = category
	return e
}

func (e *Error) WithContext(ctx interface{}) *Error {
	e.Context = ctx
	return e
}

// Logging integration
func (e *Error) LogError(logger zerolog.Logger) {
	logEvent := logger.Error().
		Int("code", e.Code).
		Str("message", e.Message).
		Time("timestamp", e.Timestamp).
		Bool("retryable", e.Retryable).
		Str("severity", e.Severity.String()).
		Str("category", e.Category.String())

	if e.Context != nil {
		logEvent = logEvent.Interface("context", e.Context)
	}

	logEvent.Msg("Service error occurred")
}

// Common errors with enhanced metadata
var (
	ErrInvalidQuery = NewError(
		http.StatusBadRequest,
		"Invalid query parameters",
		SeverityLow,
		CategoryValidation,
		false,
		nil,
	)

	ErrNotFound = NewError(
		http.StatusNotFound,
		"Resource not found",
		SeverityLow,
		CategoryDatabase,
		false,
		nil,
	)

	ErrInternalServer = NewError(
		http.StatusInternalServerError,
		"Internal server error",
		SeverityHigh,
		CategoryInternal,
		true,
		nil,
	)

	ErrThirdPartyService = NewError(
		http.StatusServiceUnavailable,
		"Third-party service error",
		SeverityHigh,
		CategoryThirdParty,
		true,
		nil,
	)
)

// IsRetryable checks if an error should be retried
func IsRetryable(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Retryable
	}
	return false
}

// HandleError is a custom error handler for Fiber
func HandleError(ctx *fiber.Ctx, err error) error {
	var apiError *Error
	if e, ok := err.(*Error); ok {
		apiError = e
	} else {
		apiError = NewError(
			http.StatusInternalServerError,
			"An unexpected error occurred",
			SeverityHigh,
			CategoryInternal,
			false,
			err,
		)
	}

	return ctx.Status(apiError.Code).JSON(apiError)
}
