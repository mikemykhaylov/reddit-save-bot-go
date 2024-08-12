package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type (
	ContextKey string
)

var logKey ContextKey = "logger"

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := NewLogger()
		ctx := context.WithValue(r.Context(), logKey, logger)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(logKey).(*slog.Logger); ok {
		return logger
	}

	return NewLogger()
}
