package logger

import (
	"context"
	"github.com/rs/zerolog"
)

type ContextualLogger struct {
	logger zerolog.Logger
}

func NewContextualLogger(baseLogger zerolog.Logger) *ContextualLogger {
	return &ContextualLogger{logger: baseLogger}
}

func (cl *ContextualLogger) WithContext(ctx context.Context) *zerolog.Event {
	event := cl.logger.Info()
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		event.Str("trace_id", traceID)
	}
	if userID, ok := ctx.Value("user_id").(string); ok {
		event.Str("user_id", userID)
	}
	return event
} 